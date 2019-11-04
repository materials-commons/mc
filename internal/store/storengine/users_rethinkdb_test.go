package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"

	r "gopkg.in/gorethink/gorethink.v4"
)

func TestUsersStoreEngineRethinkdb_AddUser(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineAddUser(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_GetUserByID(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineGetUserByID(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_GetUserByAPIKey(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineGetUserByAPIKey(t, e)
	e.Session.Close()
}

func createRethinkdbUsersStoreEngine() *storengine.UsersRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := storengine.NewUsersRethinkdb(session)
	storengine.CleanupUsersStoreEngine(e)
	return e
}
