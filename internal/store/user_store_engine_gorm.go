package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type UserStoreEngineGORM struct {
	db *gorm.DB
}

func (e *UserStoreEngineGORM) GetUserByID(id string) (UserModel, error) {
	var user UserModel
	err := e.db.First(&user, id).Error
	return user, getGORMError(err, fmt.Sprintf("No such user %s", id))
}

func (e *UserStoreEngineGORM) GetUserByAPIKey(apikey string) (UserModel, error) {
	var user UserModel
	err := e.db.First(&user, "apikey = ?", apikey).Error
	return user, getGORMError(err, fmt.Sprintf("No such user %s", apikey))
}

func getGORMError(err error, msg string) error {
	switch {
	case err == nil:
		return nil
	case gorm.IsRecordNotFoundError(err):
		return errors.Wrap(ErrNotFound, msg)
	default:
		return errors.Wrapf(err, "Internal error")
	}
}
