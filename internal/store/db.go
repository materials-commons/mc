package store

import (
	r "gopkg.in/gorethink/gorethink.v4"
)

type DB interface {
	ProjectsStore() *ProjectsStore
	UsersStore() *UsersStore
	DatafilesStore() *DatafilesStore
	DatadirsStore() *DatadirsStore
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

type DBMemory struct{}

func (db *DBMemory) ProjectsStore() *ProjectsStore {
	return NewProjectsStore(NewProjectsStoreEngineMemory())
}

func (db *DBMemory) UsersStore() *UsersStore {
	return NewUsersStore(NewUsersStoreEngineMemory())
}

func (db *DBMemory) DatafilesStore() *DatafilesStore {
	return NewDatafilesStore(NewDatafilesStoreEngineMemory())
}

func (db *DBMemory) DatadirsStore() *DatadirsStore {
	return NewDatadirsStore(NewDatadirsStoreEngineMemory())
}
