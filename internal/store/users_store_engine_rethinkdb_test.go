package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
	r "gopkg.in/gorethink/gorethink.v4"
)

func TestUsersStoreEngineRethinkdb_AddUser(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_AddUser(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_GetUserByID(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_GetUserByID(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_GetUserByAPIKey(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_GetUserByAPIKey(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_ModifyUserFullname(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_ModifyUserFullname(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_ModifyUserPassword(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_ModifyUserPassword(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_ModifyUserAPIKey(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngine_ModifyUserAPIKey(t, e)
	e.Session.Close()
}

func createRethinkdbUsersStoreEngine() *store.UsersStoreEngineRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := store.NewUsersStoreEngineRethinkdb(session)
	cleanupUsersStoreEngine(e)
	return e
}
