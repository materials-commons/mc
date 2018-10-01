package store_test

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
	r "gopkg.in/gorethink/gorethink.v4"
)

func testProjectsStoreEngineRethinkdbAddProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{
		{project: store.ProjectSchema{Model: store.Model{ID: "project1"}}, shouldFail: false, name: "Add project"},
		{project: store.ProjectSchema{Model: store.Model{ID: "project1"}}, shouldFail: true, name: "Add duplicate project"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proj, err := e.AddProject(test.project)
			if !test.shouldFail {
				assert.Okf(t, err, "Attempt to add project failed %#v, error: %s", test.project, err)
			} else {
				assert.Errorf(t, err, "Expected add project to fail for %#v, but it succeeded returning %#v", test.project, proj)
			}
		})
	}
}

func testProjectsStoreEngineRethinkdbDeleteProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		id         string
		shouldFail bool
		name       string
	}{
		{id: "project1", shouldFail: false, name: "Delete existing project"},
		{id: "project-does-not-exist", shouldFail: true, name: "Delete project that doesn't exist"},
	}

	addDefaultProjectsToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := e.DeleteProject(test.id)
			if !test.shouldFail {
				assert.Okf(t, err, "Attempt to delete project unexpectedly failed, id %s, err %s", test.id, err)
			} else {
				assert.Errorf(t, err, "Attempt to delete project should have failed id: %s", test.id)
			}
		})
	}
}

func testProjectsStoreEngineRethinkdbGetAllProjectsForUser(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	addDefaultProjectsToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdbGetProject(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	addDefaultProjectsToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdbUpdateProjectDescription(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	addDefaultProjectsToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func testProjectsStoreEngineRethinkdbUpdateProjectName(t *testing.T, e store.ProjectsStoreEngine) {
	tests := []struct {
		project    store.ProjectSchema
		shouldFail bool
		name       string
	}{{shouldFail: false}}

	addDefaultProjectsToStoreEngine(t, e)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Errorf("not implemented")
		})
	}
}

func addDefaultProjectsToStoreEngine(t *testing.T, e store.ProjectsStoreEngine) {
	projects := []store.ProjectSchema{
		{Model: store.Model{ID: "project1", Name: "project1", OType: "project", Owner: "tuser@test.com"}, Description: "project1 description"},
	}

	for _, project := range projects {
		_, err := e.AddProject(project)
		assert.Okf(t, err, "Failed to add project %s", project.ID)
	}

	//accessEntries := []store.AccessSchema{
	//	{ProjectID: "project1", UserID: "tuser@test.com"},
	//}
	//
}

func cleanupProjectsStoreEngine(e store.ProjectsStoreEngine) {
	if re, ok := e.(*store.ProjectsStoreEngineRethinkdb); ok {
		session := re.Session
		r.Table("projects").Delete().RunWrite(session)
		r.Table("datadirs").Delete().RunWrite(session)
		r.Table("project2datadir").Delete().RunWrite(session)
	}
}
