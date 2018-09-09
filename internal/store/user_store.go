package store

type UserStore struct {
	UserStoreEngine
}

func (s *UserStore) AddUser(userModel AddUserModel) (user UserSchema, err error) {
	if err := userModel.Validate(); err != nil {
		return user, err
	}

	if user, err = prepareUser(userModel); err != nil {
		return user, err
	}

	return s.UserStoreEngine.AddUser(user)
}

func (s *UserStore) GetUserByID(id string) (UserSchema, error) {
	return s.UserStoreEngine.GetUserByID(id)
}

func (s *UserStore) GetUserByAPIKey(apikey string) (UserSchema, error) {
	return s.UserStoreEngine.GetUserByAPIKey(apikey)
}
