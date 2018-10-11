package api

import (
	"net/http"

	"github.com/materials-commons/mc/internal/store"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/file"
)

type FileLoaderController struct{}

func NewFileLoaderController() *FileLoaderController {
	return &FileLoaderController{}
}

type LoadFilesReq struct {
	ProjectID string   `json:"project_id"`
	User      string   `json:"user"`
	Path      string   `json:"path"`
	Exclude   []string `json:"exclude"`
}

type fileLoaderStores struct {
	projectsStore  *store.ProjectsStore
	datafilesStore *store.DatafilesStore
	datadirsStore  *store.DatadirsStore
}

func (f *FileLoaderController) LoadFilesFromDirectory(c echo.Context) error {
	var req LoadFilesReq

	if err := c.Bind(&req); err != nil {
		return err
	}

	loadID, err := f.createLoadReq(req)
	if err != nil {
		return err
	}

	db := c.Get("DB").(store.DB)
	stores := &fileLoaderStores{
		projectsStore:  db.ProjectsStore(),
		datafilesStore: db.DatafilesStore(),
		datadirsStore:  db.DatadirsStore(),
	}

	proj, err := stores.projectsStore.GetProjectSimple(req.ProjectID)
	if err != nil {
		return err
	}

	go f.loadFiles(req, proj, stores)

	return c.JSON(http.StatusOK, map[string]interface{}{"load_id": loadID})
}

func (f *FileLoaderController) loadFiles(req LoadFilesReq, proj store.ProjectSimpleModel, stores *fileLoaderStores) {
	loader := file.NewMCFileLoader(req.Path, req.User, proj, stores.datafilesStore, stores.datadirsStore)
	skipper := file.NewExcludeListSkipper(req.Exclude)
	fl := file.NewFileLoader(skipper.Skipper, loader)
	fl.LoadFiles(req.Path)
}

func (f *FileLoaderController) createLoadReq(req LoadFilesReq) (id string, err error) {
	id = "abc123"
	return id, err
}
