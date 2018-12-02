package store_test

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	r "gopkg.in/gorethink/gorethink.v4"
)

func TestBackgroundProcessStore_AddBackgroundProcess(t *testing.T) {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	storeEngine := storengine.NewBackgroundProcessRethinkdb(session)
	globusUploads := store.NewBackgroundProcessStore(storeEngine)

	upload := model.AddBackgroundProcessModel{
		UserID:                "bogues.user@mc.org",
		ProjectID:             "ProjectId",
		BackgroundProcessID:   "BGProcessId",
		BackgroundProcessType: "bgp-type",
		Status:                "status",
		Message:               "message",
	}

	// Ensure that the id we passed in is used
	gu, err := globusUploads.AddBackgroundProcess(upload)
	assert.Okf(t, err, "Unable to add upload: %s", err)
	assert.Truef(t, upload.UserID == gu.UserID, "IDs don't match %s/%s", gu.UserID, upload.UserID)
}
