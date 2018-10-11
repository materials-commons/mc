package file_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/materials-commons/mc/internal/file"

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

	mcFileLoader := file.NewMCFileLoader("/tmp", "mctest", project, store.NewDatafilesStore(nil), ddStore)

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
	tmpFile, err := createTmpFile("/tmp", "loadFile")
	assert.Okf(t, err, "Unable to create temporary file: %s", err)
	defer os.Remove(tmpFile)

	secondFile := "/tmp/myfile.txt"
	err = ioutil.WriteFile(secondFile, []byte("second file"), 0644)
	assert.Okf(t, err, "Unable to create second file: %s", err)
	defer os.Remove(secondFile)

	finfo, err := os.Stat(tmpFile)
	assert.Okf(t, err, "Unable to stat current tmpFile (%s): %s", tmpFile, err)

	ddStoreEngine := store.NewDatadirsStoreEngineMemory()
	ddStore := store.NewDatadirsStore(ddStoreEngine)

	dfStoreEngine := store.NewDatafilesStoreEngineMemory()
	dfStore := store.NewDatafilesStore(dfStoreEngine)

	var (
		project     store.ProjectSimpleModel
		projDataDir store.DatadirSchema
		//dir         store.DatadirSchema
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

	mcFileLoader := file.NewMCFileLoader("/tmp", "mctest@test.com", project, dfStore, ddStore)

	t.Run("Check simple file load works", func(t *testing.T) {
		var df store.DatafileSchema
		err := mcFileLoader.LoadFileOrDir(tmpFile, finfo)
		assert.Okf(t, err, "Unable to load file %s: %s", tmpFile, err)

		fileBase := filepath.Base(tmpFile)
		df, err = dfStore.GetDatafileInDir(fileBase, projDataDir.ID)
		assert.Okf(t, err, "Unable to locate %s in datadir %s", fileBase, projDataDir.ID)

		assert.Truef(t, df.Size == finfo.Size(), "Datafile size in db not equal to actual file df = %d, actual = %d", df.Size, finfo.Size())
		assert.Truef(t, df.Parent == "", "Datafile parent is set and should not be")
		assert.Truef(t, df.Checksum != "", "Datafile checksum is blank")
		assert.Truef(t, df.Current, "Datafile Current flag is false")
	})

	t.Run("Check parent is correctly set", func(t *testing.T) {
		var (
			df  store.DatafileSchema
			df2 store.DatafileSchema
		)
		finfo, err := os.Stat(secondFile)
		assert.Okf(t, err, "Unable to stat %s: %s", secondFile, err)

		err = mcFileLoader.LoadFileOrDir(secondFile, finfo)
		assert.Okf(t, err, "Unable to load %s: %s", secondFile, err)

		fileBase := filepath.Base(secondFile)
		df, err = dfStore.GetDatafileInDir(fileBase, projDataDir.ID)
		assert.Okf(t, err, "Unable to locate %s in datadir %s", fileBase, projDataDir.ID)

		assert.Truef(t, df.Size == finfo.Size(), "Datafile size in db not equal to actual file df = %d, actual = %d", df.Size, finfo.Size())
		assert.Truef(t, df.Parent == "", "Datafile parent is set and should not be")
		assert.Truef(t, df.Checksum != "", "Datafile checksum is blank")
		assert.Truef(t, df.Current, "Datafile Current flag is false")

		// Change contents of file and load it a second time
		err = ioutil.WriteFile(secondFile, []byte("new contents here"), 0644)
		assert.Okf(t, err, "Unable to update %s: %s", secondFile, err)

		finfo2, err2 := os.Stat(secondFile)
		assert.Okf(t, err2, "Unable to stat %s: %s", secondFile, err)

		err = mcFileLoader.LoadFileOrDir(secondFile, finfo2)
		assert.Okf(t, err, "Unable to load %s: %s", secondFile, err)

		df2, err = dfStore.GetDatafileInDir(fileBase, projDataDir.ID)
		assert.Okf(t, err, "Unable to locate %s in datadir %s", fileBase, projDataDir.ID)

		assert.Truef(t, df2.Size == finfo2.Size(), "Datafile size in db not equal to actual file df = %d, actual = %d", df2.Size, finfo2.Size())
		assert.Truef(t, df2.Parent == df.ID, "Parent not equal to expected ID %s != %s", df2.Parent, df.ID)
		assert.Truef(t, df2.Current, "Datafile Current flag is false")

		id := df.ID
		df, err = dfStore.GetDatafileByID(df.ID)
		assert.Okf(t, err, "Unable to retrieve file %s: %s", id, err)
		assert.Truef(t, df.Size == finfo.Size(), "Datafile size in db not equal to actual file df = %d, actual = %d", df.Size, finfo.Size())
		assert.Truef(t, df.Parent == "", "Datafile parent is set and should not be")
		assert.Truef(t, df.Checksum != "", "Datafile checksum is blank")
		assert.Truef(t, df.Current == false, "Datafile Current flag is true")
	})
}
