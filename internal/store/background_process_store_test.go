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
	bgps := store.NewBackgroundProcessStore(storeEngine)

	cleanupBackgroundProcessEngine(storeEngine)

	abgpModel := model.AddBackgroundProcessModel{
		UserID:                "bogues.user@mc.org",
		ProjectID:             "ProjectId",
		BackgroundProcessID:   "BGProcessId",
		BackgroundProcessType: "bgp-type",
		Status:                "status",
		Message:               "message",
	}

	bgp, err := bgps.AddBackgroundProcess(abgpModel)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)
}

func TestBackgroundProcessStore_GetListBackgroundProcess(t *testing.T) {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	storeEngine := storengine.NewBackgroundProcessRethinkdb(session)
	bgps := store.NewBackgroundProcessStore(storeEngine)

	cleanupBackgroundProcessEngine(storeEngine)

	abgpModel := model.AddBackgroundProcessModel{
		UserID:                "bogues.user@mc.org",
		ProjectID:             "ProjectId",
		BackgroundProcessID:   "BGProcessId",
		BackgroundProcessType: "bgp-type",
		Status:                "status",
		Message:               "message",
	}

	getListModel := model.GetListBackgroundProcessModel{
		UserID:              "bogues.user@mc.org",
		ProjectID:           "ProjectId",
		BackgroundProcessID: "BGProcessId",
	}

	bgp, err := bgps.AddBackgroundProcess(abgpModel)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

	bgp, err = bgps.AddBackgroundProcess(abgpModel)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

	bgpList, err := bgps.GetListBackgroundProcess(getListModel)
	assert.Okf(t, err, "Unable to get list of matching background_process records: %s", err)

	assert.Truef(t, len(bgpList) == 2,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	for _, record := range bgpList {
		assert.Truef(t, getListModel.UserID == record.UserID, "IDs don't match %s/%s", record.UserID, abgpModel.UserID)
	}
}

func TestBackgroundProcessStore_DeleteBackgroundProcess(t *testing.T) {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	storeEngine := storengine.NewBackgroundProcessRethinkdb(session)
	bgps := store.NewBackgroundProcessStore(storeEngine)

	cleanupBackgroundProcessEngine(storeEngine)
	abgpModel := model.AddBackgroundProcessModel{
		UserID:                "bogues.user@mc.org",
		ProjectID:             "ProjectId",
		BackgroundProcessID:   "BGProcessId",
		BackgroundProcessType: "bgp-type",
		Status:                "status",
		Message:               "message",
	}

	bgp, err := bgps.AddBackgroundProcess(abgpModel)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

//    id := bgp.ID

//    err = DeleteBackgroundProcess(id)
//    assert.Okf(t, err, "Unable to delete background_process record %s: %s", id, err)

//    assert.Truef(t, false, "Stop.")
}

func cleanupBackgroundProcessEngine(e storengine.BackgroundProcessStoreEngine) {
	if re, ok := e.(*storengine.BackgroundProcessRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("background_process").Delete().RunWrite(session)
	}
}
