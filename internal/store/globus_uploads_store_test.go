package store_test

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	r "gopkg.in/gorethink/gorethink.v4"
)

func TestGlobusUploadsStore_AddGlobusUpload(t *testing.T) {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	storeEngine := storengine.NewGlobusUploadsRethinkdb(session)
	globusUploads := store.NewGlobusUploadsStore(storeEngine)

	upload := model.AddGlobusUploadModel{
		ID:               "abc123",
		Owner:            "test@test.mc",
		Path:             "/tmp",
		ProjectID:        "38a14ddb-5369-4d0d-a2b8-d87dadbc5b1f",
		GlobusAclID:      "abc123",
		GlobusEndpointID: "def456",
		GlobusIdentityID: "an-id",
	}

	// Ensure that the id we passed in is used
	gu, err := globusUploads.AddGlobusUpload(upload)
	assert.Okf(t, err, "Unable to add upload: %s", err)
	assert.Truef(t, upload.ID == gu.ID, "IDs don't match %s/%s", gu.ID, upload.ID)
}
