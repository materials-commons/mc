package store

import (
	"fmt"
	"time"
)

type UsersStoreEngineMemory struct {
	DB map[string]UserSchema
}

func NewUsersStoreEngineMemory() *UsersStoreEngineMemory {
	return &UsersStoreEngineMemory{
		DB: make(map[string]UserSchema),
	}
}

func (e *UsersStoreEngineMemory) AddUser(user UserSchema) (UserSchema, error) {
	_, ok := e.DB[user.ID]
	if ok {
		return UserSchema{}, fmt.Errorf("user already in shouldFail %s", user.ID)
	}

	e.DB[user.ID] = user

	return user, nil
}

func (e *UsersStoreEngineMemory) GetUserByID(id string) (UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return UserSchema{}, ErrNotFound
	}

	return user, nil
}

func (e *UsersStoreEngineMemory) GetUserByAPIKey(apikey string) (UserSchema, error) {
	for _, user := range e.DB {
		if user.APIKey == apikey {
			return user, nil
		}
	}

	return UserSchema{}, ErrNotFound
}

func (e *UsersStoreEngineMemory) ModifyUserFullname(id, fullname string, updatedAt time.Time) (UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return UserSchema{}, ErrNotFound
	}

	user.Fullname = fullname
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) ModifyUserPassword(id, password string, updatedAt time.Time) (UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return UserSchema{}, ErrNotFound
	}

	user.Password = password
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return UserSchema{}, ErrNotFound
	}

	user.APIKey = apikey
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) Name() string {
	return "UsersStoreEngineMemory"
}
