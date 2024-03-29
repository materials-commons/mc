package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"

	"github.com/materials-commons/mc/internal/store/model"

	"golang.org/x/crypto/bcrypt"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
)

func TestUsersStore_AddUser(t *testing.T) {
	tests := []struct {
		user       model.AddUserModel
		shouldFail bool
		name       string
	}{
		{user: model.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: false, name: "Test add valid user"},
		{user: model.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test add duplicate user"},
		{user: model.AddUserModel{Email: "bi&$bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test add with bad email"},
		{user: model.AddUserModel{Email: "", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test empty email"},
		{user: model.AddUserModel{Email: "biff2@bob.com", Password: "", Fullname: "Bob"}, shouldFail: true, name: "Test empty password"},
		{user: model.AddUserModel{Email: "biff2@bob.com", Password: "abc123", Fullname: ""}, shouldFail: true, name: "Test empty fullname"},
	}

	s := newSEMemoryUsersStore()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.AddUser(test.user)
			if !test.shouldFail {
				assert.Okf(t, err, "Adding %#v should have failed", test.user)
				assert.Truef(t, user.APIKey != "", "Added user has invalid APIKey %#v", user)
				assert.Truef(t, user.Password != test.user.Password, "Password")
				perr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(test.user.Password))
				assert.Okf(t, perr, "Hashed and password don't compare %s for user %#v", perr, test.user)
			} else {
				assert.Errorf(t, err, "Expected adding %#v to have failed", test.user)
			}
		})
	}
}

func TestUsersStore_GetUserByID(t *testing.T) {
	tests := []struct {
		id         string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", shouldFail: false, name: "Test getting an existing user"},
		{id: "does-not-exist", shouldFail: true, name: "Test getting user that doesn't exist"},
	}

	s := newSEMemoryUsersStore()
	addDefaultUsersToStore(t, s)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.GetUserByID(test.id)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to find existing user %s, error: %s", test.id, err)
			} else {
				assert.Errorf(t, err, "Found user that doesn't exist %s, %#v", test.id, user)
			}
		})
	}
}

func TestUsersStore_GetUserByAPIKey(t *testing.T) {
	tests := []struct {
		apikey     string
		shouldFail bool
		name       string
	}{
		{apikey: "tuser@test.com apikey", shouldFail: false, name: "Test getting an existing user by apikey"},
		{apikey: "no such key", shouldFail: true, name: "Test getting an existing user by apikey that doesn't exist"},
	}

	s := newSEMemoryUsersStore()
	addDefaultUsersToStore(t, s)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.GetUserByAPIKey(test.apikey)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to find existing user by apikey (%s), error: %s", test.apikey, err)
			} else {
				assert.Errorf(t, err, "Found user by apikey that doesn't exist %s, %#v", test.apikey, user)
			}
		})
	}
}

func TestUsersStore_GetAndVerifyUser(t *testing.T) {
	users := []struct {
		user model.AddUserModel
	}{
		{user: model.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}},
	}

	s := newSEMemoryUsersStore()
	for _, u := range users {
		_, err := s.AddUser(u.user)
		assert.Okf(t, err, "Failed to add user %#v, error %s", u, err)
	}

	tests := []struct {
		id         string
		password   string
		shouldFail bool
		name       string
	}{
		{id: "biff@bob.com", password: "abc123", shouldFail: false, name: "Test existing user with correct password"},
		{id: "biff@bob.com", password: "wrong-abc123", shouldFail: true, name: "Test existing user with wrong password"},
		{id: "does-not-exist@bob.com", password: "abc123", shouldFail: true, name: "Test bad user id"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := s.GetAndVerifyUser(test.id, test.password)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed to verify existing user/password (%s)/(%s), error: (%s)", test.id, test.password, err)
			} else {
				assert.Errorf(t, err, "Expected verify to fail for user/password (%s)/(%s)", test.id, test.password)
			}
		})
	}
}

func TestUsersStore_ModifyUserFullname(t *testing.T) {
	tests := []struct {
		id         string
		fullname   string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", fullname: "tuser changed", shouldFail: false, name: "Test modify existing user with a valid fullname"},
		{id: "does-not-exist@test.com", fullname: "tuser changed", shouldFail: true, name: "Test modify user that doesn't exist"},
		{id: "tuser@test.com", fullname: "", shouldFail: true, name: "Test modify existing user with an invalid fullname"},
	}
	s := newSEMemoryUsersStore()
	addDefaultUsersToStore(t, s)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.ModifyUserFullname(test.id, test.fullname)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed modifying user (%#v) when it was expected to succeed, error: %s", test, err)
				assert.Truef(t, user.Fullname == test.fullname, "Expected user password (%s) to be modified to (%s)", user.Fullname, test.fullname)
			} else {
				assert.Errorf(t, err, "Expected test %#v to fail", test)
			}
		})
	}
}

func TestUsersStore_ModifyUserPassword(t *testing.T) {
	tests := []struct {
		id         string
		password   string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", password: "tuser changed", shouldFail: false, name: "Test modify existing user with a valid password"},
		{id: "does-not-exist@test.com", password: "tuser changed", shouldFail: true, name: "Test modify user that doesn't exist"},
		{id: "tuser@test.com", password: "", shouldFail: true, name: "Test modify existing user with an invalid password"},
	}
	s := newSEMemoryUsersStore()
	addDefaultUsersToStore(t, s)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.ModifyUserPassword(test.id, test.password)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed modifying user (%#v) when it was expected to succeed, error: %s", test, err)
				perr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(test.password))
				assert.Okf(t, perr, "Hashed and password don't compare %s for user %#v", perr, test)
			} else {
				assert.Errorf(t, err, "Expected test %#v to fail", test)
			}
		})
	}
}

func TestUsersStore_ModifyUserAPIKey(t *testing.T) {
	tests := []struct {
		id         string
		oldAPIKey  string
		shouldFail bool
		name       string
	}{
		{id: "tuser@test.com", oldAPIKey: "tuser@test.com apikey", shouldFail: false, name: "Test modify existing user with a valid oldAPIKey"},
		{id: "does-not-exist@test.com", oldAPIKey: "", shouldFail: true, name: "Test modify user that doesn't exist"},
	}
	s := newSEMemoryUsersStore()
	addDefaultUsersToStore(t, s)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			user, err := s.ModifyUserAPIKey(test.id)
			if !test.shouldFail {
				assert.Okf(t, err, "Failed modifying user (%#v) when it was expected to succeed, error: %s", test, err)
				assert.Truef(t, user.APIKey != test.oldAPIKey, "Expected user APIKey (%s) to be modified", user.APIKey)
			} else {
				assert.Errorf(t, err, "Expected test %#v to fail", test)
			}
		})
	}
}

func addDefaultUsersToStore(t *testing.T, s *store.UsersStore) {
	storengine.AddDefaultUsersToStoreEngine(t, s.UsersStoreEngine)
}

func newSEMemoryUsersStore() *store.UsersStore {
	return store.NewUsersStore(storengine.NewUsersMemory())
}
