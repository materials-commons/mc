package store

import (
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
	return NewDatafilesStore(NewDatafilesStoreEngineRethinkdb(db.Session))
}

func (db *DBRethinkdb) DatadirsStore() *DatadirsStore {
	return NewDatadirsStore(NewDatadirsStoreEngineRethinkdb(db.Session))
}

func (db *DBRethinkdb) FileLoadsStore() *FileLoadsStore {
	return NewFileLoadsStore(NewFileLoadsStoreEngineRethinkdb(db.Session))
}

type DBMemory struct {
	DBProj      map[string]ProjectSchema
	DBUsers     map[string]UserSchema
	DBDatadirs  map[string]DatadirSchema
	DBDatafiles map[string]DatafileSchemaInMemory
	DBFileLoads map[string]FileLoadSchema
}

func NewDBMemory() *DBMemory {
	return &DBMemory{
		DBProj:      make(map[string]ProjectSchema),
		DBUsers:     make(map[string]UserSchema),
		DBDatadirs:  make(map[string]DatadirSchema),
		DBDatafiles: make(map[string]DatafileSchemaInMemory),
		DBFileLoads: make(map[string]FileLoadSchema),
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
		return NewDatafilesStore(NewDatafilesStoreEngineMemory())
	}

	return NewDatafilesStore(NewDatafilesStoreEngineMemoryWithDB(db.DBDatafiles))
}

func (db *DBMemory) DatadirsStore() *DatadirsStore {
	if db.DBDatadirs == nil {
		return NewDatadirsStore(NewDatadirsStoreEngineMemory())
	}

	return NewDatadirsStore(NewDatadirsStoreEngineMemoryWithDB(db.DBDatadirs))
}

func (db *DBMemory) FileLoadsStore() *FileLoadsStore {
	if db.DBFileLoads == nil {
		return NewFileLoadsStore(NewFileLoadsStoreEngineMemory())
	}

	return NewFileLoadsStore(NewFileLoadsStoreEngineMemoryWithDB(db.DBFileLoads))
}

var InMemory = NewDBMemory() // Global for testing purposes, allows a single db to be shared across test instances
