package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
)

func testProjectsStoreEngineRethinkdb_AddProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdb_DeleteProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdb_GetAllProjectsForUser(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdb_GetProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdb_UpdateProjectDescription(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdb_UpdateProjectName(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}
