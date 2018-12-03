package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"

	r "gopkg.in/gorethink/gorethink.v4"
)

func TestBackgroundProcessStoreEngineRethinkdb_BackgroundProcess(t *testing.T) {
	e := createRethinkdbBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_AddBackgroundProcess(t, e)
	e.Session.Close()
}

func TestBackgroundProcessStoreEngineRethinkdb_GetListBackgroundProcess(t *testing.T) {
	e := createRethinkdbBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_GetListBackgroundProcess(t, e)
	e.Session.Close()
}

func createRethinkdbBackgroundProcessStoreEngine() *storengine.BackgroundProcessRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := storengine.NewBackgroundProcessRethinkdb(session)
	cleanupBackgroundProcessEngine(e)
	return e
}
