package store

import (
	"time"

	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	UsersStoreEngine
}

func (s *UsersStore) AddUser(userModel AddUserModel) (user UserSchema, err error) {
	if err := userModel.Validate(); err != nil {
		return user, err
	}

	if user, err = prepareUser(userModel); err != nil {
		return user, err
	}

	return s.UsersStoreEngine.AddUser(user)
}

func (s *UsersStore) GetUserByID(id string) (UserSchema, error) {
	return s.UsersStoreEngine.GetUserByID(id)
}

func (s *UsersStore) GetUserByAPIKey(apikey string) (UserSchema, error) {
	return s.UsersStoreEngine.GetUserByAPIKey(apikey)
}

func (s *UsersStore) GetAndVerifyUser(id, password string) (UserSchema, error) {
	user, err := s.UsersStoreEngine.GetUserByID(id)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, err
}

func (s *UsersStore) ModifyUserFullname(id, fullname string) (UserSchema, error) {
	return s.UsersStoreEngine.ModifyUserFullname(id, fullname, time.Now())
}

func (s *UsersStore) ModifyUserPassword(id, password string) (UserSchema, error) {
	passwordHash, err := generatePasswordHash(password)
	if err != nil {
		return UserSchema{}, err
	}

	return s.UsersStoreEngine.ModifyUserPassword(id, passwordHash, time.Now())
}

func (s *UsersStore) ModifyUserAPIKey(id string) (UserSchema, error) {
	apikey, err := uuid.GenerateUUID()
	if err != nil {
		return UserSchema{}, err
	}

	return s.UsersStoreEngine.ModifyUserAPIKey(id, apikey, time.Now())
}
