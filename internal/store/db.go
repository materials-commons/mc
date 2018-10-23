package store

import (
	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/internal/store/storengine"
	r "gopkg.in/gorethink/gorethink.v4"
)

type DB interface {
	ProjectsStore() *ProjectsStore
	UsersStore() *UsersStore
	DatafilesStore() *DatafilesStore
	DatadirsStore() *DatadirsStore
	FileLoadsStore() *FileLoadsStore
}

type DBRethinkdb struct {
	Session *r.Session
}

func NewDBRethinkdb(session *r.Session) *DBRethinkdb {
	return &DBRethinkdb{Session: session}
}

func (db *DBRethinkdb) ProjectsStore() *ProjectsStore {
	return NewProjectsStore(NewProjectsStoreEngineRethinkdb(db.Session))
}

func (db *DBRethinkdb) UsersStore() *UsersStore {
	return NewUsersStore(NewUsersStoreEngineRethinkdb(db.Session))
}

func (db *DBRethinkdb) DatafilesStore() *DatafilesStore {
	return NewDatafilesStore(storengine.NewDatafilesRethinkdb(db.Session))
}

func (db *DBRethinkdb) DatadirsStore() *DatadirsStore {
	return NewDatadirsStore(storengine.NewDatadirsRethinkdb(db.Session))
}

func (db *DBRethinkdb) FileLoadsStore() *FileLoadsStore {
	return NewFileLoadsStore(NewFileLoadsStoreEngineRethinkdb(db.Session))
}

type DBMemory struct {
	DBProj      map[string]model.ProjectSchema
	DBUsers     map[string]model.UserSchema
	DBDatadirs  map[string]model.DatadirSchema
	DBDatafiles map[string]storengine.DatafileSchemaInMemory
	DBFileLoads map[string]model.FileLoadSchema
}

func NewDBMemory() *DBMemory {
	return &DBMemory{
		DBProj:      make(map[string]model.ProjectSchema),
		DBUsers:     make(map[string]model.UserSchema),
		DBDatadirs:  make(map[string]model.DatadirSchema),
		DBDatafiles: make(map[string]storengine.DatafileSchemaInMemory),
		DBFileLoads: make(map[string]model.FileLoadSchema),
	}
}

func (db *DBMemory) ProjectsStore() *ProjectsStore {
	if db.DBProj == nil {
		return NewProjectsStore(NewProjectsStoreEngineMemory())
	}

	return NewProjectsStore(NewProjectsStoreEngineMemoryWithDB(db.DBProj))
}

func (db *DBMemory) UsersStore() *UsersStore {
	if db.DBUsers == nil {
		return NewUsersStore(NewUsersStoreEngineMemory())
	}

	return NewUsersStore(NewUsersStoreEngineMemoryWithDB(db.DBUsers))
}

func (db *DBMemory) DatafilesStore() *DatafilesStore {
	if db.DBDatafiles == nil {
		return NewDatafilesStore(storengine.NewDatafilesMemory())
	}

	return NewDatafilesStore(storengine.NewDatafilesMemoryWithDB(db.DBDatafiles))
}

func (db *DBMemory) DatadirsStore() *DatadirsStore {
	if db.DBDatadirs == nil {
		return NewDatadirsStore(storengine.NewDatadirsMemory())
	}

	return NewDatadirsStore(storengine.NewDatadirsMemoryWithDB(db.DBDatadirs))
}

func (db *DBMemory) FileLoadsStore() *FileLoadsStore {
	if db.DBFileLoads == nil {
		return NewFileLoadsStore(NewFileLoadsStoreEngineMemory())
	}

	return NewFileLoadsStore(NewFileLoadsStoreEngineMemoryWithDB(db.DBFileLoads))
}

var InMemory = NewDBMemory() // Global for testing purposes, allows a single db to be shared across test instances
