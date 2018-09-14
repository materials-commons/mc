package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
)

func TestUsersStoreEngineMemory_AddUser(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineAddUser(t, e)
}

func TestUsersStoreEngineMemory_GetUserByID(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineGetUserByID(t, e)
}

func TestUsersStoreEngineMemory_GetUserByAPIKey(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineGetUserByAPIKey(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserFullname(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserFullname(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserPassword(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserPassword(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserAPIKey(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserAPIKey(t, e)
}
