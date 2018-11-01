package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-uuid"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/globus"
	"github.com/pkg/errors"
)

const globusBaseURL = "https://www.globus.org/app/transfer"

type GlobusController struct {
	client             *globus.Client
	globusUploadsStore *store.GlobusUploadsStore
	basePath           string
	globusEndpointID   string
}

func NewGlobusController(db store.DB, client *globus.Client, basePath, globusEndpointID string) *GlobusController {
	return &GlobusController{
		client:             client,
		globusUploadsStore: db.GlobusUploadsStore(),
		basePath:           basePath,
		globusEndpointID:   globusEndpointID,
	}
}

// CreateGlobusUploadRequests creates a new entry in the globus_uploads table that tracks to a directory on
// the materials commons globus endpoint that a user can upload to. The directory is created on the endpoint
// and that user has an ACL set on it to permit reading and writing to it. This directory is temporary and
// only available for the upload.
func (g *GlobusController) CreateGlobusUploadRequest(c echo.Context) error {
	var req struct {
		ProjectID string `json:"project_id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)

	globusResp, err := g.createAndSetupUploadReq(req.ProjectID, user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, globusResp)
}

type globusResp struct {
	GlobusURL string `json:"globus_url"`
	ID        string `json:"id"`
}

func (g *GlobusController) createAndSetupUploadReq(projectID string, user model.UserSchema) (globusResp, error) {
	var (
		err  error
		resp globusResp
	)

	gUploadModel := model.AddGlobusUploadModel{
		Owner:            user.ID,
		ProjectID:        projectID,
		GlobusEndpointID: g.globusEndpointID,
	}

	// Usually we let the database create the ID. However in this case we have multiple operations that depend
	// on the ID. So in this case we pre-create the ID and tell the database what the ID is.
	if gUploadModel.ID, err = uuid.GenerateUUID(); err != nil {
		return resp, err
	}

	gUploadModel.Path = filepath.Join(g.basePath, "__globus_uploads", gUploadModel.ID)
	if err := os.MkdirAll(gUploadModel.Path, 0700); err != nil {
		return resp, err
	}

	gUploadModel.GlobusIdentityID, gUploadModel.GlobusAclID, err = g.globusSetup(gUploadModel.ID, gUploadModel.Path, user.GlobusUser)
	if err != nil {
		return resp, err
	}

	if _, err := g.globusUploadsStore.AddGlobusUpload(gUploadModel); err != nil {
		return resp, err
	}

	resp.ID = gUploadModel.ID
	resp.GlobusURL = g.createEndpointURL(gUploadModel.ID)

	return resp, nil
}

// globusSetup performs a couple of opertions related to globus. It takes the users globus login and translates that into
// and identity id. The identity id is used to set the ACL on the directory in the end point for materials commons.
func (g *GlobusController) globusSetup(uploadID, path string, globusUser string) (globusIdentityID string, aclID int, err error) {
	identities, err := g.client.GetIdentities([]string{globusUser})
	if err != nil {
		return globusIdentityID, aclID, errors.WithMessage(err, fmt.Sprintf("Unable to retrieve globus user from globus api %s", globusUser))
	}

	globusIdentityID = identities.Included.Identities[0].ID

	rule := globus.EndpointACLRule{
		EndpointID:  g.globusEndpointID,
		Path:        path,
		IdentityID:  globusIdentityID,
		Permissions: "rw",
	}

	aclRes, err := g.client.AddEndpointACLRule(rule)
	if err != nil {
		msg := fmt.Sprintf("Unable to add endpoint rule for endpoint %s, path %s, user %s/%s", g.globusEndpointID, path, globusUser, globusIdentityID)
		return globusIdentityID, aclID, errors.WithMessage(err, msg)
	}

	return globusIdentityID, aclRes.AccessID, nil
}

// createEndpointURL creates the url that the front end can use to bring up the globus UI webapp and have the destination
// panel (right side) pointing to the materials commons endpoint and correct directory.
func (g *GlobusController) createEndpointURL(uploadID string) string {
	path := fmt.Sprintf("/__globus_uploads/%s", uploadID)
	return fmt.Sprintf("%s?destination_id=%s&destination_path=%s", globusBaseURL, g.globusEndpointID, path)
}
