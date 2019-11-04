package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"
	"github.com/materials-commons/mc/pkg/mc"
)

type UsersMemory struct {
	DB map[string]model.UserSchema
}

func NewUsersMemory() *UsersMemory {
	return &UsersMemory{
		DB: make(map[string]model.UserSchema),
	}
}

func NewUsersMemoryWithDB(db map[string]model.UserSchema) *UsersMemory {
	return &UsersMemory{
		DB: db,
	}
}

func (e *UsersMemory) AddUser(user model.UserSchema) (model.UserSchema, error) {
	_, ok := e.DB[user.ID]
	if ok {
		return model.UserSchema{}, fmt.Errorf("user already in shouldFail %s", user.ID)
	}

	e.DB[user.ID] = user

	return user, nil
}

func (e *UsersMemory) GetUserByID(id string) (model.UserSchema, error) {
	user, ok := e.DB[id]
	if !ok {
		return model.UserSchema{}, mc.ErrNotFound
	}

	return user, nil
}

func (e *UsersMemory) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	for _, user := range e.DB {
		if user.APIKey == apikey {
			return user, nil
		}
	}

	return model.UserSchema{}, mc.ErrNotFound
}
