package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"
)

func TestUsersStoreEngineMemory_AddUser(t *testing.T) {
	e := storengine.NewUsersMemory()
	testUsersStoreEngineAddUser(t, e)
}

func TestUsersStoreEngineMemory_GetUserByID(t *testing.T) {
	e := storengine.NewUsersMemory()
	testUsersStoreEngineGetUserByID(t, e)
}

func TestUsersStoreEngineMemory_GetUserByAPIKey(t *testing.T) {
	e := storengine.NewUsersMemory()
	testUsersStoreEngineGetUserByAPIKey(t, e)
}
