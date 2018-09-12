package store

import "time"

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
	ModifyUserFullname(id, fullname string, updatedAt time.Time) (UserSchema, error)
	ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error)
	ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error)
	Name() string
}

type ProjectsStoreEngine interface {
	AddProject(project ProjectSchema) (ProjectSchema, error)
	GetProject(id string) (ProjectExtendedModel, error)
	GetAllProjectsForUser(user string) ([]ProjectExtendedModel, error)
	DeleteProject(id string) error
	UpdateProjectName(id string, name string, updatedAt time.Time) error
	UpdateProjectDescription(id string, description string, updatedAt time.Time) error
	Name() string
}
