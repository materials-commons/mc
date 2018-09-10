package store_test

import (
	"fmt"
	"testing"

	"github.com/materials-commons/mc/internal/store"
)

func TestUsersStore_AddUser(t *testing.T) {
	s := newSEMemoryUsersStore()
	fmt.Println(s)
}

func TestUsersStore_GetUserByID(t *testing.T) {

}

func TestUsersStore_GetUserByAPIKey(t *testing.T) {

}

func TestUsersStore_GetAndVerifyUser(t *testing.T) {

}

func TestUsersStore_ModifyUserFullname(t *testing.T) {

}

func TestUsersStore_ModifyUserPassword(t *testing.T) {

}

func TestUsersStore_ModifyUserAPIKey(t *testing.T) {

}

func initUserStore(t *testing.T, s *store.UsersStore) {
	addDefaultUsersToStoreEngine(t, s.UsersStoreEngine)
}

func newSEMemoryUsersStore() *store.UsersStore {
	return &store.UsersStore{UsersStoreEngine: store.NewUsersStoreEngineMemory()}
}
