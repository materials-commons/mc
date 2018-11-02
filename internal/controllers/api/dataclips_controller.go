package api

import (
	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
)

type DataclipsController struct {
	db store.DB
}

func NewDataclipsController(db store.DB) *DataclipsController {
	return &DataclipsController{db: db}
}

func CreateDataclip(c echo.Context) error {
	var req struct {
		ProjectID    string
		ExperimentID string
		Public       bool
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	return nil
}

func DeleteDataclip(c echo.Context) error {
	return nil
}

func GetDataForDataclip(c echo.Context) error {
	return nil
}
