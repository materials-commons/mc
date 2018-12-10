package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	"github.com/materials-commons/mc/pkg/tutils/assert"

	r "gopkg.in/gorethink/gorethink.v4"
)

func testBackgroundProcessStoreEngine_AddBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {
	tests := []struct {
		bgp        model.BackgroundProcessSchema
		shouldFail bool
		name       string
	}{
		{bgp: model.BackgroundProcessSchema{ModelSimple: model.ModelSimple{ID: "bgp1"}},
			shouldFail: true,
			name:       "Add existing"},
		{bgp: model.BackgroundProcessSchema{ModelSimple: model.ModelSimple{ID: "bgp_new"}},
			shouldFail: false,
			name:       "Add new background_process"},
	}

	addDefaultBackgroundProcessToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bgp, err := e.AddBackgroundProcess(test.bgp)
			if !test.shouldFail {
				assert.Okf(t, err, "Could not add background_process id %+v, error: %s", test.bgp, err)
			} else {
				assert.Errorf(t, err, "Added background_process that already exists %+v", bgp)
			}
		})
	}
}

func testBackgroundProcessStoreEngine_GetBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {

	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}

	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add BackgroundProcess: %s", err)
	assert.Truef(t, bgpSchema.UserID == bgp.UserID, "Add for Get: User IDs don't match '%s'/'%s'", bgp.UserID, bgpSchema.UserID)

	id := bgp.ID
	bgp, err = e.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get BackgroundProcess, %s: %s", id, err)
	assert.Truef(t, bgpSchema.UserID == bgp.UserID, "Get: User IDs don't match '%s'/'%s'", bgp.UserID, bgpSchema.UserID)
}

func testBackgroundProcessStoreEngine_SetFinishedBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {

	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}

	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, !bgp.IsFinished, "Initial background_process record incorrectly marked finished")

	id := bgp.ID

	err = e.SetFinishedBackgroundProcess(id, true)
	assert.Okf(t, err, "Unable to set finished flag on background_process record, %s: %s", id, err)

	bgp, err = e.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get background process, %s: %s", id, err)

	assert.Truef(t, bgp.IsFinished, "Updated background_process record incorrectly marked not finished")
}

func testBackgroundProcessStoreEngine_SetOkBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {

	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}

	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add abgpModel: %s", err)
	assert.Truef(t, !bgp.IsOk, "Initial background_process record incorrectly marked ok")

	id := bgp.ID

	err = e.SetOkBackgroundProcess(id, true)
	assert.Okf(t, err, "Unable to set ok flag on background_process record, %s: %s", id, err)

	bgp, err = e.GetBackgroundProcess(id)
	assert.Okf(t, err, "Unable to get background process, %s: %s", id, err)

	assert.Truef(t, bgp.IsOk, "Updated background_process record incorrectly marked not ok")
}

func testBackgroundProcessStoreEngine_GetListBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {

	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}

	getListModel := model.GetListBackgroundProcessModel{
		UserID:           "bogues.user@mc.org",
		ProjectID:        "ProjectId",
		BackgroundTaskID: "BGProcessId",
	}

	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add bgpSchema: %s", err)
	assert.Truef(t, bgpSchema.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, bgpSchema.UserID)

	bgp, err = e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add bgpSchema: %s", err)
	assert.Truef(t, bgpSchema.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, bgpSchema.UserID)

	bgpList, err := e.GetListBackgroundProcess(getListModel)
	assert.Okf(t, err, "Unable to get list of matching background_process records: %s", err)

	assert.Truef(t, len(bgpList) == 2,
		"Unexpected length in returned list of background_process records, %v", len(bgpList))

	for _, record := range bgpList {
		assert.Truef(t, getListModel.UserID == record.UserID, "IDs don't match %s/%s", record.UserID, bgpSchema.UserID)
	}
}

func testBackgroundProcessStoreEngine_DeleteBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {
	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}
	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add bgpSchema: %s", err)
	assert.Truef(t, bgpSchema.UserID == bgp.UserID, "IDs don't match %s/%s", bgp.UserID, bgpSchema.UserID)

	id := bgp.ID

	err = e.DeleteBackgroundProcess(id)
	assert.Okf(t, err, "Unable to delete bgpSchema: %s", err)
}

func testBackgroundProcessStoreEngine_UpdateStatusBackgroundProcess(t *testing.T, e storengine.BackgroundProcessStoreEngine) {
	bgpSchema := model.BackgroundProcessSchema{
		UserID:             "bogues.user@mc.org",
		ProjectID:          "ProjectId",
		BackgroundTaskID:   "BGProcessId",
		BackgroundTaskType: "bgp-type",
		Status:             "status",
		Message:            "message",
	}
	bgp, err := e.AddBackgroundProcess(bgpSchema)
	assert.Okf(t, err, "Unable to add bgpSchema: %s", err)
	assert.Truef(t, bgpSchema.Status == bgp.Status, "Status Fields don't match '%s'/'%s' ", bgp.Status, bgpSchema.Status)

	id := bgp.ID
	newStatus := "new status"
	newMessage := "new message"
	err = e.UpdateStatusBackgroundProcess(id, newStatus, newMessage)
	assert.Okf(t, err, "Unable to update background_process record, %s: %s", id, err)

	bgp, err = e.GetBackgroundProcess(id)
	assert.Truef(t, newStatus == bgp.Status, "Status Fields don't match '%s'/'%s' ", bgp.Status, newStatus)
}

func addDefaultBackgroundProcessToStoreEngine(t *testing.T, e storengine.BackgroundProcessStoreEngine) {
	background_process_records := []model.BackgroundProcessSchema{
		{ModelSimple: model.ModelSimple{ID: "bgp1"}},
	}

	for _, bgp := range background_process_records {
		_, err := e.AddBackgroundProcess(bgp)
		assert.Okf(t, err, "Failed to add background_process %s", bgp.ID)
	}
}

func cleanupBackgroundProcessEngine(e storengine.BackgroundProcessStoreEngine) {
	if re, ok := e.(*storengine.BackgroundProcessRethinkdb); ok {
		session := re.Session
		_, _ = r.Table("background_process").Delete().RunWrite(session)
	}
	if me, ok := e.(*storengine.BackgroundProcessMemory); ok {
		me.DB = make(map[string]model.BackgroundProcessSchema)
	}
}
