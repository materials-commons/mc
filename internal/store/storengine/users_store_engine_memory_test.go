package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"
)

func TestUsersStoreEngineMemory_AddUser(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineAddUser(t, e)
}

func TestUsersStoreEngineMemory_GetUserByID(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineGetUserByID(t, e)
}

func TestUsersStoreEngineMemory_GetUserByAPIKey(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineGetUserByAPIKey(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserFullname(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserFullname(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserPassword(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserPassword(t, e)
}

func TestUsersStoreEngineMemory_ModifyUserAPIKey(t *testing.T) {
	e := storengine.NewUsersStoreEngineMemory()
	testUsersStoreEngineModifyUserAPIKey(t, e)
}
