package store

type UsersStoreEngine interface {
	AddUser(user UserSchema) (UserSchema, error)
	GetUserByID(id string) (UserSchema, error)
	GetUserByAPIKey(apikey string) (UserSchema, error)
}
