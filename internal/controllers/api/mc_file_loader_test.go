package api_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestMCFileLoaderLoadOnlyADirectory(t *testing.T) {
	// Create a finfo so we can pass it to LoadFileOrDir, doesn't matter what it points at as long as its a directory
	finfo, err := os.Stat(".")
	assert.Okf(t, err, "Unable to stat current dir %s", err)

	ddStoreEngine := store.NewDatadirsStoreEngineMemory()
	ddStore := store.NewDatadirsStore(ddStoreEngine)

	var (
		project     store.ProjectSimpleModel
		projDataDir store.DatadirSchema
		dir         store.DatadirSchema
	)
	project.Name = "My Project"
	project.ID = "My Project ID"

	dirToAdd := store.AddDatadirModel{
		Name:      project.Name,
		Owner:     "mctest",
		Parent:    "",
		ProjectID: project.ID,
	}

	projDataDir, err = ddStore.AddDatadir(dirToAdd)
	assert.Okf(t, err, "Unable to add project dir %s", err)

	mcFileLoader := api.NewMCFileLoader("/tmp", "mctest", project, store.NewDatafilesStore(nil), ddStore)

	// Load /tmp/dir for the first time
	t.Run("Initial load of directory /tmp/dir", func(t *testing.T) {
		err := mcFileLoader.LoadFileOrDir("/tmp/dir", finfo)
		assert.Okf(t, err, "Unable to load /tmp/dir %s", err)

		dirToFind := filepath.Join(project.Name, "dir")
		dir, err = ddStoreEngine.GetDatadirByPathInProject(dirToFind, project.ID)
		assert.Okf(t, err, "Unable to find path (%s) projectID (%s): %s", dirToFind, project.ID, err)
		assert.Truef(t, dir.Parent == projDataDir.ID, "Expected created dir to have parent of projDataDir %#v", dir)
		expectedName := filepath.Join(project.Name, "dir")
		assert.Truef(t, dir.Name == expectedName, "Expected dir to have name (%s) instead got (%s)", expectedName, dir.Name)
	})

	t.Run("Load /tmp/dir a second time and make sure no new entry is created", func(t *testing.T) {
		err := mcFileLoader.LoadFileOrDir("/tmp/dir", finfo)
		assert.Okf(t, err, "Unable to load /tmp/dir %s", err)

		dirToFind := filepath.Join(project.Name, "dir")
		dirAgain, err := ddStoreEngine.GetDatadirByPathInProject(dirToFind, project.ID)
		assert.Okf(t, err, "Unable to find path (%s) projectID (%s): %s", dirToFind, project.ID, err)
		assert.Truef(t, dirAgain.Parent == projDataDir.ID, "Expected created dir to have parent of projDataDir %#v", dirAgain)
		expectedName := filepath.Join(project.Name, "dir")
		assert.Truef(t, dirAgain.Name == expectedName, "Expected dir to have name (%s) instead got (%s)", expectedName, dirAgain.Name)
		assert.Truef(t, dirAgain.ID == dir.ID, "Expected dirAgain.ID (%s) == dir.ID (%s)", dirAgain.ID, dir.ID)
	})
}

func TestMCFileLoaderLoadOnlyAFile(t *testing.T) {
	//tmpFile, err := createTmpFile("/tmp", "loadFile")
	//assert.Okf(t, err, "Unable to create temporary file: %s", err)
	//defer os.Remove(tmpFile)
	//
	//finfo, err := os.Stat(tmpFile)
	//assert.Okf(t, err, "Unable to stat current tmpFile (%s): %s", tmpFile, err)
	//
	//ddStoreEngine := store.NewDatadirsStoreEngineMemory()
	//ddStore := store.NewDatadirsStore(ddStoreEngine)

}
