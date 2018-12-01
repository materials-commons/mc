package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"

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
	cleanupBackgroundProcessEngine(e)
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
}
