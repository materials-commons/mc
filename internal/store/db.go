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
	dbProj      map[string]ProjectSchema
	dbUsers     map[string]UserSchema
	dbDatadirs  map[string]DatadirSchema
	dbDatafiles map[string]DatafileSchemaInMemory
	dbFileLoads map[string]FileLoadSchema
}

func NewDBMemory() *DBMemory {
	return &DBMemory{
		dbProj:      make(map[string]ProjectSchema),
		dbUsers:     make(map[string]UserSchema),
		dbDatadirs:  make(map[string]DatadirSchema),
		dbDatafiles: make(map[string]DatafileSchemaInMemory),
		dbFileLoads: make(map[string]FileLoadSchema),
	}
}

func (db *DBMemory) ProjectsStore() *ProjectsStore {
	if db.dbProj == nil {
		return NewProjectsStore(NewProjectsStoreEngineMemory())
	}

	return NewProjectsStore(NewProjectsStoreEngineMemoryWithDB(db.dbProj))
}

func (db *DBMemory) UsersStore() *UsersStore {
	if db.dbUsers == nil {
		return NewUsersStore(NewUsersStoreEngineMemory())
	}

	return NewUsersStore(NewUsersStoreEngineMemoryWithDB(db.dbUsers))
}

func (db *DBMemory) DatafilesStore() *DatafilesStore {
	if db.dbDatafiles == nil {
		return NewDatafilesStore(NewDatafilesStoreEngineMemory())
	}

	return NewDatafilesStore(NewDatafilesStoreEngineMemoryWithDB(db.dbDatafiles))
}

func (db *DBMemory) DatadirsStore() *DatadirsStore {
	if db.dbDatadirs == nil {
		return NewDatadirsStore(NewDatadirsStoreEngineMemory())
	}

	return NewDatadirsStore(NewDatadirsStoreEngineMemoryWithDB(db.dbDatadirs))
}

func (db *DBMemory) FileLoadsStore() *FileLoadsStore {
	if db.dbFileLoads == nil {
		return NewFileLoadsStore(NewFileLoadsStoreEngineMemory())
	}

	return NewFileLoadsStore(NewFileLoadsStoreEngineMemoryWithDB(db.dbFileLoads))
}

var InMemory = NewDBMemory() // Global for testing purposes, allows a single db to be shared across test instances
