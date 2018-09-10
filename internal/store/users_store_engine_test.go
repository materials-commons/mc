package store

import (
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestAddUser(t *testing.T) {
	e := &UsersStoreEngineMemory{}
	tests := []struct {
		user       UserSchema
		shouldFail bool
		name       string
	}{
		{user: UserSchema{ID: "gtarcea@umich.edu"}, shouldFail: false, name: "New user"},
		{user: UserSchema{ID: "gtarcea@umich.edu"}, shouldFail: true, name: "Existing user"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := e.AddUser(test.user)
			if !test.shouldFail {
				assert.Ok(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {

}

func TestGetUserByAPIKey(t *testing.T) {

}

func TestModifyUserFullname(t *testing.T) {

}

func TestModifyUserPassword(t *testing.T) {

}

func TestModifyUserAPIKey(t *testing.T) {

}
