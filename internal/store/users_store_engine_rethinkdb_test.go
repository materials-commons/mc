package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
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

func TestUsersStoreEngineRethinkdb_ModifyUserFullname(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineModifyUserFullname(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_ModifyUserPassword(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineModifyUserPassword(t, e)
	e.Session.Close()
}

func TestUsersStoreEngineRethinkdb_ModifyUserAPIKey(t *testing.T) {
	e := createRethinkdbUsersStoreEngine()
	testUsersStoreEngineModifyUserAPIKey(t, e)
	e.Session.Close()
}

func createRethinkdbUsersStoreEngine() *store.UsersStoreEngineRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := store.NewUsersStoreEngineRethinkdb(session)
	cleanupUsersStoreEngine(e)
	return e
}
