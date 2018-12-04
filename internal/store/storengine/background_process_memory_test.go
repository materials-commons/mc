package storengine_test

import (
	"testing"

	"github.com/materials-commons/mc/internal/store/storengine"
)

func TestBackgroundProcessStoreEngineMemory_BackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_AddBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_GetBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_GetBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_DeleteBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_DeleteBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_GetListBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_GetListBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_UpdateStatusBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_UpdateStatusBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_SetFinishedBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_SetFinishedBackgroundProcess(t, e)
}

func TestBackgroundProcessStoreEngineMemory_SetOKBackgroundProcess(t *testing.T) {
	e := createMemoryBackgroundProcessStoreEngine()
	testBackgroundProcessStoreEngine_SetOkBackgroundProcess(t, e)
}

func createMemoryBackgroundProcessStoreEngine() *storengine.BackgroundProcessMemory {
	e := storengine.NewBackgroundProcessMemory()
	return e
}
