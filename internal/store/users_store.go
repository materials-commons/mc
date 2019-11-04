package store

import (
	"time"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/hashicorp/go-uuid"
	"github.com/materials-commons/mc/internal/store/model"
	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	UsersStoreEngine storengine.UsersStoreEngine
}

func NewUsersStore(e storengine.UsersStoreEngine) *UsersStore {
	return &UsersStore{e}
}

func (s *UsersStore) AddUser(userModel model.AddUserModel) (user model.UserSchema, err error) {
	if err := userModel.Validate(); err != nil {
		return user, err
	}

	if user, err = prepareUser(userModel); err != nil {
		return user, err
	}

	return s.UsersStoreEngine.AddUser(user)
}

func (s *UsersStore) GetUserByID(id string) (model.UserSchema, error) {
	return s.UsersStoreEngine.GetUserByID(id)
}

func (s *UsersStore) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	return s.UsersStoreEngine.GetUserByAPIKey(apikey)
}

func (s *UsersStore) GetAndVerifyUser(id, password string) (model.UserSchema, error) {
	user, err := s.UsersStoreEngine.GetUserByID(id)
	if err != nil {
		return user, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return user, err
}

func prepareUser(userModel model.AddUserModel) (model.UserSchema, error) {
	var (
		err error
	)

	now := time.Now()

	u := model.UserSchema{
		ModelSimpleNoID: model.ModelSimpleNoID{
			Birthtime: now,
			MTime:     now,
			OType:     "user",
		},
		ID:       userModel.Email,
		Fullname: userModel.Fullname,
		Email:    userModel.Email,
	}

	if u.Password, err = generatePasswordHash(userModel.Password); err != nil {
		return u, err
	}

	if u.APIKey, err = uuid.GenerateUUID(); err != nil {
		return u, err
	}

	return u, nil
}

func generatePasswordHash(password string) (passwordHash string, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}
