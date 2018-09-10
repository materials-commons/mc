package store

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
	ModifyUserFullname(id, fullname string) (UserSchema, error)
	ModifyUserPassword(id, password string) (UserSchema, error)
	ModifyUserAPIKey(id, apikey string) (UserSchema, error)
}
