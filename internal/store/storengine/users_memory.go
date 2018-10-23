package storengine

import (
	"fmt"
	"time"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type UsersStoreEngineMemory struct {
	DB map[string]model.UserSchema
}

func NewUsersStoreEngineMemory() *UsersStoreEngineMemory {
	return &UsersStoreEngineMemory{
		DB: make(map[string]model.UserSchema),
	}
}

func NewUsersStoreEngineMemoryWithDB(db map[string]model.UserSchema) *UsersStoreEngineMemory {
	return &UsersStoreEngineMemory{
		DB: db,
	}
}

func (e *UsersStoreEngineMemory) AddUser(user model.UserSchema) (model.UserSchema, error) {
	_, ok := e.DB[user.ID]
	if ok {
		return model.UserSchema{}, fmt.Errorf("user already in shouldFail %s", user.ID)
	}

	e.DB[user.ID] = user

	return user, nil
}

func (e *UsersStoreEngineMemory) GetUserByID(id string) (model.UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return model.UserSchema{}, mc.ErrNotFound
	}

	return user, nil
}

func (e *UsersStoreEngineMemory) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	for _, user := range e.DB {
		if user.APIKey == apikey {
			return user, nil
		}
	}

	return model.UserSchema{}, mc.ErrNotFound
}

func (e *UsersStoreEngineMemory) ModifyUserFullname(id, fullname string, updatedAt time.Time) (model.UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return model.UserSchema{}, mc.ErrNotFound
	}

	user.Fullname = fullname
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) ModifyUserPassword(id, password string, updatedAt time.Time) (model.UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return model.UserSchema{}, mc.ErrNotFound
	}

	user.Password = password
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) ModifyUserAPIKey(id, apikey string, updatedAt time.Time) (model.UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return model.UserSchema{}, mc.ErrNotFound
	}

	user.APIKey = apikey
	user.MTime = updatedAt
	e.DB[id] = user
	return user, nil
}

func (e *UsersStoreEngineMemory) Name() string {
	return "UsersStoreEngineMemory"
}
