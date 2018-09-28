package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/materials-commons/mc/internal/store"
)

type MCFileLoader struct {
	root             string
	datafilesStore   *store.DatafilesStore
	project          store.ProjectSimpleModel
	currentDatadirID string
}

func NewMCFileLoader(root string, project store.ProjectSimpleModel, dfStore *store.DatafilesStore) *MCFileLoader {
	return &MCFileLoader{
		root:           root,
		datafilesStore: dfStore,
		project:        project,
	}
}

func (l *MCFileLoader) LoadFileOrDir(path string, finfo os.FileInfo) error {
	switch {

	case !finfo.Mode().IsRegular() && !finfo.IsDir():
		// Not a regular file, don't do anything
		return nil

	case finfo.IsDir():
		return l.loadDirectory(path, finfo)

	default:
		// Is a file
		return l.loadFile(path, finfo)
	}
}

func (l *MCFileLoader) loadDirectory(path string, finfo os.FileInfo) error {
	fmt.Printf("loadDirectory '%s' '%s'\n", l.root, filepath.Join(l.project.Name, strings.TrimPrefix(path, l.root+"/")))
	return nil
}

func (l *MCFileLoader) loadFile(path string, finfo os.FileInfo) error {
	return nil
}
