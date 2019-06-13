package api

import (
	"net/http"
	"path/filepath"

	"github.com/apex/log"
	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/ds"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/globusapi"
	r "gopkg.in/gorethink/gorethink.v4"
)

type GlobusDatasetController struct {
	client   *globusapi.Client
	dsStore  *store.DatasetsStore
	session  *r.Session
	basePath string
}

func NewGlobusDatasetController(db store.DB, client *globusapi.Client, basePath string, session *r.Session) *GlobusDatasetController {
	return &GlobusDatasetController{
		client:   client,
		dsStore:  db.DatasetsStore(),
		session:  session,
		basePath: basePath,
	}
}

func (g *GlobusDatasetController) CreateGlobusDatasetDownload(c echo.Context) error {
	var (
		req struct {
			DatasetID string `json:"dataset_id"`
			ProjectID string `json:"project_id"`
		}

		datasetGlobusPath string
	)

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user := c.Get("User").(model.UserSchema)

	dataset, err := g.dsStore.GetDataset(req.DatasetID)
	if err != nil {
		log.Infof("Failed getting dataset %s", req.DatasetID)
		return ToHttpError(err)
	}

	if dataset.Owner != user.ID {
		log.Infof("User doesn't have permission to create dataset directory")
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	if dataset.Published {
		datasetGlobusPath = filepath.Join(g.basePath, "__published_datasets/%s", req.DatasetID)
	} else {
		datasetGlobusPath = filepath.Join(g.basePath, "__datasets/%s", req.DatasetID)
	}

	dsDirLoader := ds.NewDirLoader(datasetGlobusPath, g.session)
	if err := dsDirLoader.LoadDirFromDataset(dataset, req.ProjectID); err != nil {
		return ToHttpError(err)
	}

	return c.JSON(http.StatusOK, req)
}
