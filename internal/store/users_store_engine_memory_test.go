package store_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store"
)

func TestUsersStoreEngineMemory_AddUser(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_AddUser(t, e)
}

func TestUsersStoreEngineMemory_GetUserByID(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_GetUserByID(t, e)
}

func TestUsersStoreEngineMemory_GetUserByAPIKey(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_GetUserByAPIKey(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserFullname(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_ModifyUserFullname(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserPassword(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_ModifyUserPassword(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserAPIKey(t *testing.T) {
	e := store.NewUsersStoreEngineMemory()
	testUsersStoreEngine_ModifyUserAPIKey(t, e)
}
