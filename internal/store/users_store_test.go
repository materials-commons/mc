package store_test

import (
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/store"
)

func TestUsersStore_AddUser(t *testing.T) {
	tests := []struct {
		user       store.AddUserModel
		shouldFail bool
		name       string
	}{
		{user: store.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: false, name: "Test add valid user"},
		{user: store.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test add duplicate user"},
		{user: store.AddUserModel{Email: "bi&$bob.com", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test add with bad email"},
		{user: store.AddUserModel{Email: "", Password: "abc123", Fullname: "Bob"}, shouldFail: true, name: "Test empty email"},
		{user: store.AddUserModel{Email: "biff2@bob.com", Password: "", Fullname: "Bob"}, shouldFail: true, name: "Test empty password"},
		{user: store.AddUserModel{Email: "biff2@bob.com", Password: "abc123", Fullname: ""}, shouldFail: true, name: "Test empty fullname"},
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
		user store.AddUserModel
	}{
		{user: store.AddUserModel{Email: "biff@bob.com", Password: "abc123", Fullname: "Bob"}},
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

}

func TestUsersStore_ModifyUserPassword(t *testing.T) {

}

func TestUsersStore_ModifyUserAPIKey(t *testing.T) {

}

func addDefaultUsersToStore(t *testing.T, s *store.UsersStore) {
	addDefaultUsersToStoreEngine(t, s.UsersStoreEngine)
}

func newSEMemoryUsersStore() *store.UsersStore {
	return store.NewUsersStore(store.NewUsersStoreEngineMemory())
}
