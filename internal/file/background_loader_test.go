package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/materials-commons/mc/internal/store/model"

	"github.com/materials-commons/mc/pkg/tutils"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
)

func TestBackgroundLoader_loadFiles(t *testing.T) {
	var (
		loadTree string
		tmpFile  string
		proj     model.ProjectSchema
		fl       model.FileLoadSchema
	)
	tdir, err := tutils.PrepareTestDirTree("mcdir")
	assert.Okf(t, err, "Unable to create temporary dir: %s", err)
	defer os.RemoveAll(tdir)
	mcdir := filepath.Join(tdir, "mcdir")

	loadTree, err = tutils.PrepareTestDirTree("dir1")
	assert.Okf(t, err, "Unable to create loadTree dir: %s", err)

	tmpFile, err = tutils.CreateTmpFile(filepath.Join(loadTree, "dir1"), "hello world")
	assert.Okf(t, err, "Unable to create temporary file: %s", err)
	fmt.Println("tmpFile", tmpFile)

	loader := NewBackgroundLoader(mcdir, 1, store.InMemory)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	loader.c = ctx
	fileloadsStore := store.InMemory.FileLoadsStore()
	projectsStore := store.InMemory.ProjectsStore()
	ddirStore := store.InMemory.DatadirsStore()

	pAdd := model.AddProjectModel{
		Name:  "proj1",
		Owner: "test@mc.org",
	}
	proj, err = projectsStore.AddProject(pAdd)
	assert.Okf(t, err, "Unable to add project: %s", err)

	dAdd := model.AddDatadirModel{
		Name:      proj.Name,
		Owner:     proj.Owner,
		ProjectID: proj.ID,
	}
	_, err = ddirStore.AddDatadir(dAdd)
	assert.Okf(t, err, "Unable to add directory: %s", err)

	flAddModel := model.AddFileLoadModel{
		ProjectID: proj.ID,
		Path:      loadTree,
		Owner:     proj.Owner,
	}
	fl, err = fileloadsStore.AddFileLoad(flAddModel)
	assert.Okf(t, err, "Unable to add file load request: %s", err)

	rv := loader.worker(fl)
	if rv == nil {
		err = nil
	} else {
		err = rv.(error)
	}

	assert.Okf(t, err, "Unable to load request: %s", err)

	assert.Truef(t, len(store.InMemory.DBDatadirs) == 2, "Didn't get expected number (2) of dirs: %d", len(store.InMemory.DBDatadirs))

	assert.Truef(t, len(store.InMemory.DBDatafiles) == 1, "Didn't get expected number (1) of files: %d", len(store.InMemory.DBDatafiles))

	for _, dfile := range store.InMemory.DBDatafiles {
		fmt.Printf("%#v\n", dfile)
		filePath := filepath.Join(MCFileDir(mcdir, dfile.DataFile.ID), dfile.DataFile.ID)
		assert.Truef(t, Exists(filePath), "File doesn't exist: %s", filePath)
	}
}
