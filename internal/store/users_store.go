package store

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/hashicorp/go-uuid"
	"golang.org/x/crypto/bcrypt"
)

type UsersStore struct {
	UsersStoreEngine UsersStoreEngine
}

func NewUsersStore(e UsersStoreEngine) *UsersStore {
	return &UsersStore{e}
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
	if err := validation.Validate(fullname, validation.Required, validation.Length(1, 40)); err != nil {
		return UserSchema{}, err
	}
	return s.UsersStoreEngine.ModifyUserFullname(id, fullname, time.Now())
}

func (s *UsersStore) ModifyUserPassword(id, password string) (UserSchema, error) {
	if err := validation.Validate(password, validation.Required, validation.Length(1, 100)); err != nil {
		return UserSchema{}, err
	}

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

func prepareUser(userModel AddUserModel) (UserSchema, error) {
	var (
		err error
	)

	now := time.Now()

	u := UserSchema{
		ModelSimple: ModelSimple{
			Birthtime: now,
			MTime:     now,
			ID:        userModel.Email,
			OType:     "user",
		},
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
