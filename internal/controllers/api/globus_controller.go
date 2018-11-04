package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/materials-commons/mc/pkg/mc"

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

// GetGlobusUploadRequest will retrieve the specified request so long as the
// user is the owner of the request, or the user has the Admin flag set to true
func (g *GlobusController) GetGlobusUploadRequest(c echo.Context) error {
	var req struct {
		ID string `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)

	globusUploadReq, err := g.globusUploadsStore.GetGlobusUpload(req.ID)

	switch {
	case errors.Cause(err) == mc.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return err
	case user.ID != globusUploadReq.Owner && !user.Admin:
		return echo.ErrUnauthorized
	default:
		return c.JSON(http.StatusOK, globusUploadReq)
	}
}

// ListGlobusUploadRequests will retrieve all the upload requests for a specific user. Admin
// users can set user to "all" to retrieve all the known upload requests. Users can only retrieve
// their own upload requests unless they are an admin. Admins can retrieve other users requests.
func (g *GlobusController) ListGlobusUploadRequests(c echo.Context) error {
	var (
		req struct {
			User string `json:"user"`
		}
		uploads []model.GlobusUploadSchema
		err     error
	)

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)
	switch {
	case req.User == "all" && user.Admin:
		// Admin user is allowed to get all requests
		uploads, err = g.globusUploadsStore.GetAllGlobusUploads()
	case req.User == "all" && !user.Admin:
		return echo.ErrUnauthorized
	case req.User != user.ID && !user.Admin:
		return echo.ErrUnauthorized
	default:
		// Either req.User == user.ID, or req.User != user.ID but user.Admin == True
		uploads, err = g.globusUploadsStore.GetAllGlobusUploadsForUser(req.User)
	}

	if err != nil && errors.Cause(err) == mc.ErrNotFound {
		return echo.ErrNotFound
	} else if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, uploads)
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
func (g *GlobusController) globusSetup(uploadID, path string, globusUser string) (globusIdentityID string, aclID string, err error) {
	identities, err := g.client.GetIdentities([]string{globusUser})
	if err != nil {
		return globusIdentityID, aclID, errors.WithMessage(err, fmt.Sprintf("Unable to retrieve globus user from globus api %s", globusUser))
	}

	globusIdentityID = identities.Identities[0].ID

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
