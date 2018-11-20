package uiapi

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
)

type DatadirsController struct {
	datadirsStore *store.DatadirsStore
}

func NewDatadirsController(db store.DB) *DatadirsController {
	return &DatadirsController{datadirsStore: db.DatadirsStore()}
}

func (d *DatadirsController) GetDirectoryForProject(c echo.Context) error {
	var req struct {
		ProjectID   string `json:"project_id"`
		DirectoryID string `json:"directory_id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)

	if dirs, err := d.datadirsStore.GetDatadirForProject(req.ProjectID, user.ID, req.DirectoryID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, dirs)
	}
}

func (d *DatadirsController) GetFilesForDirectory(c echo.Context) error {
	var req struct {
		ProjectID   string `json:"project_id"`
		DirectoryID string `json:"directory_id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user := c.Get("User").(model.UserSchema)
	if files, err := d.datadirsStore.GetFilesForDatadir(req.ProjectID, user.ID, req.DirectoryID); err != nil {
		return err
	} else {
		return c.JSON(http.StatusOK, files)
	}
}
