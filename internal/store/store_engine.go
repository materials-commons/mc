package store

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
	GetAndVerifyUser(id, password string) (UserSchema, error)
	ModifyUser(id, fullname string) (UserSchema, error)
	ModifyUserPassword(id, password string) (UserSchema, error)
	ResetAPIKey(id string) (UserSchema, error)
}
