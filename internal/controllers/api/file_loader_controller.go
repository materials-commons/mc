package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo"
)

type FileLoaderController struct {
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

	loadID, err := f.createLoadReq(req)
	if err != nil {
		return err
	}

	go f.loadFiles(req)

	return c.JSON(http.StatusOK, map[string]interface{}{"load_id": loadID})
}

func (f *FileLoaderController) loadFiles(req LoadFilesReq) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() && info.Name() == "" {
			fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
			return filepath.SkipDir
		}
		fmt.Printf("visited file or dir: %q\n", path)
		return nil
	})
	fmt.Println("err", err)
}

func (f *FileLoaderController) createLoadReq(req LoadFilesReq) (id string, err error) {
	return id, err
}
