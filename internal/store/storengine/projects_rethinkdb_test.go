package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"

	r "gopkg.in/gorethink/gorethink.v4"
)

func TestProjectsStoreEngineRethinkdb_AddProject(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbAddProject(t, e)
	e.Session.Close()
}

func TestProjectsStoreEngineRethinkdb_DeleteProject(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbDeleteProject(t, e)
	e.Session.Close()
}

func TestProjectsStoreEngineRethinkdb_GetAllProjectsForUser(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbGetAllProjectsForUser(t, e)
	e.Session.Close()
}

func TestProjectsStoreEngineRethinkdb_GetProject(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbGetProject(t, e)
	e.Session.Close()
}

func TestProjectsStoreEngineRethinkdb_UpdateProjectDescription(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbUpdateProjectDescription(t, e)
	e.Session.Close()
}

func TestProjectsStoreEngineRethinkdb_UpdateProjectName(t *testing.T) {
	e := createRethinkdbProjectsStoreEngine()
	testProjectsStoreEngineRethinkdbUpdateProjectName(t, e)
	e.Session.Close()
}

func createRethinkdbProjectsStoreEngine() *storengine.ProjectsRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := storengine.NewProjectsRethinkdb(session)
	cleanupProjectsStoreEngine(e)
	return e
}
