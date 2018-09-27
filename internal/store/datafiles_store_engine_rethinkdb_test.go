package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
	r "gopkg.in/gorethink/gorethink.v4"
)

func TestDatafilesStoreEngineRethinkdb_AddFile(t *testing.T) {
	e := createRethinkdbDatafilesStoreEngine()
	testDatafilesStoreEngineAddFile(t, e)
	e.Session.Close()
}

func TestDatafilesStoreEngineRethinkdb_GetFile(t *testing.T) {
	e := createRethinkdbDatafilesStoreEngine()
	testDatafilesStoreEngineGetFile(t, e)
	e.Session.Close()
}

func TestDatafilesStoreEngineRethinkdb_GetFileWithChecksum(t *testing.T) {
	e := createRethinkdbDatafilesStoreEngine()
	testDatafilesStoreEngineGetFileWithChecksum(t, e)
	e.Session.Close()
}

func TestDatafilesStoreEngineRethinkdb_GetFileInDir(t *testing.T) {
	e := createRethinkdbDatafilesStoreEngine()
	testDatafilesStoreEngineGetFileInDir(t, e)
	e.Session.Close()
}

func createRethinkdbDatafilesStoreEngine() *store.DatafilesStoreEngineRethinkdb {
	session, _ := r.Connect(r.ConnectOpts{Database: "mctest", Address: "localhost:30815"})
	r.SetTags("r")
	e := store.NewDatafilesStoreEngineRethinkdb(session)
	cleanupDatafilesStoreEngine(e)
	return e
}
