package api

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/materials-commons/mc/internal/store"
)

type MCFileLoader struct {
	root               string
	dfStore            *store.DatafilesStore
	ddStore            *store.DatadirsStore
	project            store.ProjectSimpleModel
	currentRootDatadir store.DatadirSchema
	owner              string
}

func NewMCFileLoader(root, owner string, project store.ProjectSimpleModel, dfStore *store.DatafilesStore, ddStore *store.DatadirsStore) *MCFileLoader {
	return &MCFileLoader{
		root:    root,
		owner:   owner,
		dfStore: dfStore,
		ddStore: ddStore,
		project: project,
	}
}

func (l *MCFileLoader) LoadFileOrDir(path string, finfo os.FileInfo) error {
	var err error
	dirName := filepath.Dir(filepath.Join(l.project.Name, strings.TrimPrefix(path, l.root+"/")))
	if l.currentRootDatadir.Name != dirName {
		l.currentRootDatadir, err = l.ddStore.GetDatadirByPathInProject(dirName, l.project.ID)
		if err != nil {
			errors.Wrapf(err, "Unable to get dir for path (%s) and project ID (%s) as new current dir", dirName, l.project.ID)
		}
	}

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
	dirPath := filepath.Join(l.project.Name, strings.TrimPrefix(path, l.root+"/"))
	if _, err := l.ddStore.GetDatadirByPathInProject(dirPath, l.project.ID); err == nil {
		return nil
	}

	dir := store.AddDatadirModel{
		Name:      dirPath,
		Owner:     l.owner,
		Parent:    l.currentRootDatadir.ID,
		ProjectID: l.project.ID,
	}

	if _, err := l.ddStore.AddDatadir(dir); err != nil {
		return errors.Wrapf(err, "Unable to create Datadir for path (%s) in project (%s)", dirPath, l.project.ID)
	}

	return nil
}

func (l *MCFileLoader) loadFile(path string, finfo os.FileInfo) error {
	return nil
}
