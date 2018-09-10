package store

import (
	"testing"
	"time"

	"github.com/materials-commons/mc/pkg/tutils/assert"
)

func TestAddUser(t *testing.T) {
	e := NewUsersStoreEngineMemory()
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
				assert.Okf(t, err, "Failed to add user %s", test.user.ID)
			} else {
				assert.Errorf(t, err, "Added existing user %s", test.user.ID)
			}
		})
	}
}

func TestGetUserByID(t *testing.T) {
	tests := []struct {
		id         string
		shouldFail bool
		name       string
	}{
		{id: "gtarcea@umich.edu", shouldFail: false, name: "Find existing user"},
		{id: "nosuchuser@doesnot.exist", shouldFail: true, name: "Fail to find a non-existing user"},
	}

	e := NewUsersStoreEngineMemory()
	addDefaultUsers(t, e)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := e.GetUserByID(test.id)
			if !test.shouldFail {
				assert.Okf(t, err, "Should have found user %s, error %s", test.id, err)
				assert.Truef(t, user.ID == test.id, "User and retrieved user not equal %s != %s", user.ID, test.id)
			} else {
				assert.Errorf(t, err, "Should have failed retrieving user %s", test.id)
			}
		})
	}
}

func TestGetUserByAPIKey(t *testing.T) {
	tests := []struct {
		id         string
		apikey     string
		shouldFail bool
		name       string
	}{
		{id: "gtarcea@umich.edu", apikey: "gtarcea@umich.edu apikey", shouldFail: false, name: "Look up existing apikey"},
		{id: "nosuchuser@doesnot.exist", apikey: "no such key", shouldFail: true, name: "Lookup apikey that doesn't exist"},
	}

	e := NewUsersStoreEngineMemory()
	addDefaultUsers(t, e)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := e.GetUserByAPIKey(test.apikey)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to look up existing apikey %s", test.apikey)
				assert.Truef(t, user.APIKey == test.apikey, "APIKeys don't match for user/expected %s/%s", user.APIKey, test.apikey)
				assert.Truef(t, user.ID == test.id, "User IDs don't match for user/expected %s/%s", user.ID, test.id)
			} else {
				assert.Errorf(t, err, "Found apikey that should not exist %s, user %s", test.apikey, user.ID)
			}
		})
	}

}

func TestModifyUserFullname(t *testing.T) {
	tests := []struct {
		id          string
		newFullname string
		shouldFail  bool
		name        string
	}{
		{id: "gtarcea@umich.edu", newFullname: "gtarcea-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newFullname: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	e := NewUsersStoreEngineMemory()
	addDefaultUsers(t, e)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := e.ModifyUserFullname(test.id, test.newFullname, time.Now())
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to modify existing user fullname id %s", test.id)
				assert.Truef(t, user.Fullname == test.newFullname, "Expected fullname to equal %s, instead got %s", test.newFullname, user.Fullname)
			} else {
				assert.Errorf(t, err, "Attempt to modify user (%s) who does not exist succeeded", test.id)
			}
		})
	}
}

func TestModifyUserPassword(t *testing.T) {
	tests := []struct {
		id          string
		newPassword string
		shouldFail  bool
		name        string
	}{
		{id: "gtarcea@umich.edu", newPassword: "gtarcea-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newPassword: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	e := NewUsersStoreEngineMemory()
	addDefaultUsers(t, e)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := e.ModifyUserPassword(test.id, test.newPassword, time.Now())
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to modify existing user password id %s", test.id)
				assert.Truef(t, user.Password == test.newPassword, "Expected password to equal %s, instead got %s", test.newPassword, user.Password)
			} else {
				assert.Errorf(t, err, "Attempt to modify user (%s) who does not exist succeeded", test.id)
			}
		})
	}
}

func TestModifyUserAPIKey(t *testing.T) {
	tests := []struct {
		id         string
		newAPIKey  string
		shouldFail bool
		name       string
	}{
		{id: "gtarcea@umich.edu", newAPIKey: "gtarcea-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newAPIKey: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	e := NewUsersStoreEngineMemory()
	addDefaultUsers(t, e)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := e.ModifyUserAPIKey(test.id, test.newAPIKey, time.Now())
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to modify existing user APIKey id %s", test.id)
				assert.Truef(t, user.APIKey == test.newAPIKey, "Expected APIKey to equal %s, instead got %s", test.newAPIKey, user.APIKey)
			} else {
				assert.Errorf(t, err, "Attempt to modify user (%s) who does not exist succeeded", test.id)
			}
		})
	}
}

func addDefaultUsers(t *testing.T, e UsersStoreEngine) {
	users := []UserSchema{
		{ID: "gtarcea@umich.edu", APIKey: "gtarcea@umich.edu apikey", Fullname: "gtarcea", Password: "gtarcea-password"},
	}

	for _, user := range users {
		_, err := e.AddUser(user)
		assert.Okf(t, err, "Failed to add user %s", user.ID)
	}
}
