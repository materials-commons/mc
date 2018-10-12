package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/materials-commons/mc/internal/controllers/api"
	"github.com/materials-commons/mc/internal/store"
)

func TestFileLoaderController_LoadFilesFromDirectory(t *testing.T) {
	memoryDB := &store.DBMemory{}
	uploadController := api.NewFileLoaderController(memoryDB)
	req := api.LoadFilesReq{
		ProjectID: "project-id",
		User:      "test@mc.org",
		Path:      "/tmp/dir",
		Exclude:   []string{"dir2/file.txt"},
	}

	j, err := json.Marshal(req)
	assert.Okf(t, err, "Unable to marshal req %#v into JSON: %s", req, err)

	c, rec := setupEcho(j)

	var resp struct {
		LoadID string `json:"load_id"`
	}

	err = uploadController.LoadFilesFromDirectory(c)
	assert.Okf(t, err, "Failed to process req %#v: %s", req, err)
	assert.Truef(t, rec.Code == http.StatusOK, "rec.Code (%d) != http.StatusOk", rec.Code)

	err = json.Unmarshal(rec.Body.Bytes(), &resp)
	assert.Okf(t, err, "Unable to marshal response: %s", err)

	assert.Truef(t, resp.LoadID == "abc123", "Expected LoadID = abc123, got %s", resp.LoadID)
}

func setupEcho(jsonBytes []byte) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", bytes.NewReader(jsonBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c, rec
}
