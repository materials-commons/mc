package storengine

import (
	"fmt"

	"github.com/materials-commons/mc/internal/store/model"

	r "gopkg.in/gorethink/gorethink.v4"
)

type UsersRethinkdb struct {
	Session *r.Session
}

func NewUsersRethinkdb(session *r.Session) *UsersRethinkdb {
	return &UsersRethinkdb{Session: session}
}

func (e *UsersRethinkdb) AddUser(user model.UserSchema) (model.UserSchema, error) {
	errMsg := fmt.Sprintf("Unable to add user %+v", user)
	resp, err := r.Table("users").Insert(user, r.InsertOpts{ReturnChanges: true}).RunWrite(e.Session)
	return user, checkRethinkdbInsertError(resp, err, errMsg)
}

func (e *UsersRethinkdb) GetUserByID(id string) (model.UserSchema, error) {
	var user model.UserSchema
	errMsg := fmt.Sprintf("No such user %s", id)
	res, err := r.Table("users").Get(id).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}
	defer res.Close()

	err = res.One(&user)
	return user, err
}

func (e *UsersRethinkdb) GetUserByAPIKey(apikey string) (model.UserSchema, error) {
	var user model.UserSchema
	errMsg := fmt.Sprintf("No such apikey %s", apikey)
	res, err := r.Table("users").GetAllByIndex("apikey", apikey).Run(e.Session)
	if err := checkRethinkdbQueryError(res, err, errMsg); err != nil {
		return user, err
	}
	defer res.Close()

	err = res.One(&user)
	return user, err
}
