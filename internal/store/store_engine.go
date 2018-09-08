package store

type UserStoreEngine interface {
	GetUserByID(id string) (UserModel, error)
	GetUserByAPIKey(apikey string) (UserModel, error)
}
