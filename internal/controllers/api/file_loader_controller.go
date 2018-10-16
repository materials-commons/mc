package api

import (
	"net/http"

	"github.com/materials-commons/mc/internal/store"

	"github.com/labstack/echo"
)

type FileLoaderController struct {
	fileloadsStore *store.FileLoadsStore
}

func NewFileLoaderController(db store.DB) *FileLoaderController {
	return &FileLoaderController{fileloadsStore: db.FileLoadsStore()}
}

type LoadFilesReq struct {
	ProjectID string   `json:"project_id"`
	User      string   `json:"user"`
	Path      string   `json:"path"`
	Exclude   []string `json:"exclude"`
}

func (f *FileLoaderController) LoadFilesFromDirectory(c echo.Context) error {
	var req LoadFilesReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	fileLoadID, err := f.createLoadReq(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"file_load_id": fileLoadID})
}

func (f *FileLoaderController) createLoadReq(req LoadFilesReq) (id string, err error) {
	flAdd := store.AddFileLoadModel{
		ProjectID: req.ProjectID,
		Owner:     req.User,
		Path:      req.Path,
		Exclude:   req.Exclude,
	}

	fl, err := f.fileloadsStore.AddFileLoad(flAdd)
	return fl.ID, err
}
