package api

import (
	"net/http"

	"github.com/materials-commons/mc/internal/store"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/file"
)

type fileLoaderStores struct {
	projectsStore  *store.ProjectsStore
	datafilesStore *store.DatafilesStore
	datadirsStore  *store.DatadirsStore
	fileloadsStore *store.FileLoadsStore
}

type FileLoaderController struct {
	stores *fileLoaderStores
}

func NewFileLoaderController(db store.DB) *FileLoaderController {
	stores := &fileLoaderStores{
		projectsStore:  db.ProjectsStore(),
		datafilesStore: db.DatafilesStore(),
		datadirsStore:  db.DatadirsStore(),
		fileloadsStore: db.FileLoadsStore(),
	}
	return &FileLoaderController{stores: stores}
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

	proj, err := f.stores.projectsStore.GetProjectSimple(req.ProjectID)
	if err != nil {
		return err
	}

	go f.loadFiles(req, proj)

	return c.JSON(http.StatusOK, map[string]interface{}{"file_load_id": fileLoadID})
}

func (f *FileLoaderController) loadFiles(req LoadFilesReq, proj store.ProjectSimpleModel) {
	loader := file.NewMCFileLoader(req.Path, req.User, proj, f.stores.datafilesStore, f.stores.datadirsStore)
	skipper := file.NewExcludeListSkipper(req.Exclude)
	fl := file.NewFileLoader(skipper.Skipper, loader)
	_ = fl.LoadFiles(req.Path)
}

func (f *FileLoaderController) createLoadReq(req LoadFilesReq) (id string, err error) {
	flAdd := store.AddFileLoadModel{
		ProjectID: req.ProjectID,
		Owner:     req.User,
		Path:      req.Path,
		Exclude:   req.Exclude,
	}

	fl, err := f.stores.fileloadsStore.AddFileLoad(flAdd)
	return fl.ID, err
}
