package store

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
