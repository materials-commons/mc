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
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "User IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

	cleanupBackgroundProcessEngine(storeEngine)
}

func TestBackgroundProcessStore_GetBackgroundProcess(t *testing.T) {
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

	id := bgp.ID

	bgp, err = bgps.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get background process, %s: %s", id, err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "User IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

	cleanupBackgroundProcessEngine(storeEngine)
}

func TestBackgroundProcessStore_SetFinishedBackgroundProcess(t *testing.T) {
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
	assert.Truef(t, !bgp.IsFinished, "Initial background_process record incorrectly marked finished")

	id := bgp.ID

	err = bgps.SetFinishedBackgroundProcess(id, true)
	assert.Okf(t, err, "Unable to set finished flag on background_process record, %s: %s", id, err)

	bgp, err = bgps.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get background process, %s: %s", id, err)

	assert.Truef(t, bgp.IsFinished, "Updated background_process record incorrectly marked not finished")
	cleanupBackgroundProcessEngine(storeEngine)
}

func TestBackgroundProcessStore_SetOKBackgroundProcess(t *testing.T) {
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
	assert.Truef(t, !bgp.IsOk, "Initial background_process record incorrectly marked ok")

	id := bgp.ID

	err = bgps.SetOkBackgroundProcess(id, true)
	assert.Okf(t, err, "Unable to set ok flag on background_process record, %s: %s", id, err)

	bgp, err = bgps.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get background process, %s: %s", id, err)

	assert.Truef(t, bgp.IsOk, "Updated background_process record incorrectly marked not ok")
	cleanupBackgroundProcessEngine(storeEngine)
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

	cleanupBackgroundProcessEngine(storeEngine)
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

	getListModel := model.GetListBackgroundProcessModel{
		UserID:              "bogues.user@mc.org",
		ProjectID:           "ProjectId",
		BackgroundProcessID: "BGProcessId",
	}

	bgp, err := bgps.AddBackgroundProcess(abgpModel)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, abgpModel.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, abgpModel.UserID)

	bgpList, err := bgps.GetListBackgroundProcess(getListModel)
	assert.Okf(t, err, "Unable to get list of matching background_process records: %s", err)

	assert.Truef(t, len(bgpList) == 1,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	for _, record := range bgpList {
		assert.Truef(t, getListModel.UserID == record.UserID, "IDs don't match %s/%s", record.UserID, abgpModel.UserID)
	}

	id := bgp.ID

	err = bgps.DeleteBackgroundProcess(id)
	assert.Okf(t, err, "Unable to delete background_process record %s: %s", id, err)

	bgpList, _ = bgps.GetListBackgroundProcess(getListModel)
	assert.Truef(t, len(bgpList) == 0,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	cleanupBackgroundProcessEngine(storeEngine)
}

func TestBackgroundProcessStore_UpdateStatusBackgroundProcess(t *testing.T) {
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
	assert.Truef(t, abgpModel.Status == bgp.Status, "Status Fields don't match %s/%s", bgp.Status, abgpModel.Status)

	bgpList, err := bgps.GetListBackgroundProcess(getListModel)
	assert.Okf(t, err, "Unable to get list of matching background_process records: %s", err)

	assert.Truef(t, len(bgpList) == 1,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	for _, record := range bgpList {
		assert.Truef(t, getListModel.UserID == record.UserID, "IDs don't match %s/%s", record.UserID, abgpModel.UserID)
	}

	id := bgp.ID

	newStatus := "new status"
	newMessage := "new message"
	err = bgps.UpdateStatusBackgroundProcess(id, newStatus, newMessage)
	assert.Okf(t, err, "Unable to update background_process record, %s: %s", id, err)

	bgpList, err = bgps.GetListBackgroundProcess(getListModel)
	assert.Okf(t, err, "Unable to get list of matching background_process records: %s", err)

	assert.Truef(t, len(bgpList) == 1,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	bgp = bgpList[0]
	assert.Truef(t, newStatus == bgp.Status, "Status Fields don't match %s/%s", bgp.Status, newStatus)
	cleanupBackgroundProcessEngine(storeEngine)
}

func cleanupBackgroundProcessEngine(e storengine.BackgroundProcessStoreEngine) {
	if re, ok := e.(*storengine.BackgroundProcessRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("background_process").Delete().RunWrite(session)
	}
}
