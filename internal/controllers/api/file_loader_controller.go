package api

import (
	"net/http"

	"github.com/materials-commons/mc/internal/store/model"

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

// LoadFilesFromDirectory accepts LoadFilesReq and creates a file load request
// in the file_loads table. The background file loader (see mc/cmd/mcserv/root.go)
// will pull items out of the database and process the load request.
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
	flAdd := model.AddFileLoadModel{
		ProjectID: req.ProjectID,
		Owner:     req.User,
		Path:      req.Path,
		Exclude:   req.Exclude,
	}

	fl, err := f.fileloadsStore.AddFileLoad(flAdd)
	return fl.ID, err
}

func (f *FileLoaderController) GetFilesLoadRequest(c echo.Context) error {
	var req struct {
		ID string `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	fileLoad, err := f.fileloadsStore.GetFileLoad(req.ID)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, fileLoad)
}
