package store_test

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
	r "gopkg.in/gorethink/gorethink.v4"
)

func testDatafilesStoreEngineGetFile(t *testing.T, e store.DatafilesStoreEngine) {
	tests := []struct {
		fileid     string
		shouldFail bool
		name       string
	}{
		{fileid: "datafile1", shouldFail: false, name: "Get existing file"},
		{fileid: "does-not-exist", shouldFail: true, name: "Get file that doesn't exist"},
	}

	addDefaultDatafilesToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := e.GetFile(test.fileid)
			if !test.shouldFail {
				assert.Okf(t, err, "Could not find file id %s, error: %s", test.fileid, err)
			} else {
				assert.Errorf(t, err, "Found file that shouldn't exist %+v", f)
			}
		})
	}
}

func testDatafilesStoreEngineAddFile(t *testing.T, e store.DatafilesStoreEngine) {
	tests := []struct {
		file       store.DatafileSchema
		shouldFail bool
		name       string
	}{
		{file: store.DatafileSchema{Model: store.Model{ID: "datafile1"}}, shouldFail: true, name: "Add existing"},
		{file: store.DatafileSchema{Model: store.Model{Name: "newfile"}}, shouldFail: false, name: "Add new file"},
	}

	addDefaultDatafilesToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := e.AddFile(test.file, "", "")
			if !test.shouldFail {
				assert.Okf(t, err, "Could not add file id %+v, error: %s", test.file, err)
			} else {
				assert.Errorf(t, err, "Added file that already exists %+v", f)
			}
		})
	}
}

func testDatafilesStoreEngineGetFileWithChecksum(t *testing.T, e store.DatafilesStoreEngine) {
	tests := []struct {
		checksum   string
		shouldFail bool
		name       string
	}{
		{checksum: "csumdatafile1", shouldFail: false, name: "Get file with existing checksum"},
		{checksum: "does-not-exist", shouldFail: true, name: "Attempt to find a checksum that doesn't exist"},
	}

	addDefaultDatafilesToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := e.GetFileWithChecksum(test.checksum)
			if !test.shouldFail {
				assert.Okf(t, err, "Could not find file with checksum %s, error: %s", test.checksum, err)
			} else {
				assert.Errorf(t, err, "Found file with checksum %s that shouldn't exist %+v", test.checksum, f)
			}
		})
	}
}

func testDatafilesStoreEngineGetFileInDir(t *testing.T, e store.DatafilesStoreEngine) {
	// To be implemented
}

func addDefaultDatafilesToStoreEngine(t *testing.T, e store.DatafilesStoreEngine) {
	datafiles := []store.DatafileSchema{
		{Model: store.Model{ID: "datafile1"}, Checksum: "csumdatafile1"},
	}

	for _, datafile := range datafiles {
		_, err := e.AddFile(datafile, "", "")
		assert.Okf(t, err, "Failed to add file %s", datafile.ID)
	}
}

func cleanupDatafilesStoreEngine(e store.DatafilesStoreEngine) {
	if re, ok := e.(*store.DatafilesStoreEngineRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("datafiles").Delete().RunWrite(session)
		_, _ = r.Table("project2datafile").Delete().RunWrite(session)
		_, _ = r.Table("datadir2datafile").Delete().RunWrite(session)
	}
}
