package store

import "time"

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
	ModifyUserFullname(id, fullname string, updatedAt time.Time) (UserSchema, error)
	ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error)
	ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error)
}
