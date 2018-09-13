package store_test

import (
	"testing"
	"time"

	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/pkg/tutils/assert"
	r "gopkg.in/gorethink/gorethink.v4"
)

func testUsersStoreEngine_AddUser(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		user       store.UserSchema
		shouldFail bool
		name       string
	}{
		{user: store.UserSchema{ModelSimple: store.ModelSimple{ID: "tuser@test.com"}}, shouldFail: false, name: "New user"},
		{user: store.UserSchema{ModelSimple: store.ModelSimple{ID: "tuser@test.com"}}, shouldFail: true, name: "Existing user"},
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
	cleanupUsersStoreEngine(e)
}

func testUsersStoreEngine_GetUserByID(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		id         string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", shouldFail: false, name: "Find existing user"},
		{id: "nosuchuser@doesnot.exist", shouldFail: true, name: "Fail to find a non-existing user"},
	}

	addDefaultUsersToStoreEngine(t, e)
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
	cleanupUsersStoreEngine(e)
}

func testUsersStoreEngine_GetUserByAPIKey(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		id         string
		apikey     string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", apikey: "tuser@test.com apikey", shouldFail: false, name: "Look up existing apikey"},
		{id: "nosuchuser@doesnot.exist", apikey: "no such key", shouldFail: true, name: "Lookup apikey that doesn't exist"},
	}

	addDefaultUsersToStoreEngine(t, e)
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
	cleanupUsersStoreEngine(e)
}

func testUsersStoreEngine_ModifyUserFullname(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		id          string
		newFullname string
		shouldFail  bool
		name        string
	}{
		{id: "tuser@test.com", newFullname: "tuser-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newFullname: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	addDefaultUsersToStoreEngine(t, e)
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
	cleanupUsersStoreEngine(e)
}

func testUsersStoreEngine_ModifyUserPassword(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		id          string
		newPassword string
		shouldFail  bool
		name        string
	}{
		{id: "tuser@test.com", newPassword: "tuser-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newPassword: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	addDefaultUsersToStoreEngine(t, e)
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
	cleanupUsersStoreEngine(e)
}

func testUsersStoreEngine_ModifyUserAPIKey(t *testing.T, e store.UsersStoreEngine) {
	tests := []struct {
		id         string
		newAPIKey  string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", newAPIKey: "tuser-changed", shouldFail: false, name: "Set fullname for existing user"},
		{id: "doesnot@exist.com", newAPIKey: "nosuch-changed", shouldFail: true, name: "Set fullname for non-existing user"},
	}

	addDefaultUsersToStoreEngine(t, e)
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
	cleanupUsersStoreEngine(e)
}

func addDefaultUsersToStoreEngine(t *testing.T, e store.UsersStoreEngine) {
	users := []store.UserSchema{
		{ModelSimple: store.ModelSimple{ID: "tuser@test.com"}, APIKey: "tuser@test.com apikey", Fullname: "tuser", Password: "tuser-password"},
	}

	for _, user := range users {
		_, err := e.AddUser(user)
		assert.Okf(t, err, "Failed to add user %s", user.ID)
	}
}

func cleanupUsersStoreEngine(e store.UsersStoreEngine) {
	if re, ok := e.(*store.UsersStoreEngineRethinkdb); ok {
		session := re.Session
		r.Table("users").Delete().RunWrite(session)
	}
}
