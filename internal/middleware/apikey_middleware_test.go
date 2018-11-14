package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/materials-commons/mc/pkg/tutils/assert"

	"github.com/labstack/echo/middleware"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"

	"github.com/labstack/echo"
)

func TestAPIKeyAuth(t *testing.T) {
	e := echo.New()
	usersStore := store.InMemory.UsersStore()

	config := APIKeyConfig{
		Skipper: middleware.DefaultSkipper,
		Keyname: "apikey",
		Retriever: func(apikey string, c echo.Context) (*model.UserSchema, error) {
			user, err := usersStore.GetUserByAPIKey(apikey)
			return &user, err
		},
	}

	// Create a handler that is wrapped with the APIKeyAuth middleware.
	// Now we can perform various tests using the handler.
	handler := APIKeyAuth(config)(func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	})

	testUserToAdd := model.AddUserModel{
		Email:    "t1@mc.org",
		Fullname: "t1",
		Password: "abc123",
	}

	user, err := usersStore.AddUser(testUserToAdd)
	assert.Okf(t, err, "Unable to add user %#v: %s", testUserToAdd, err)

	tests := []struct {
		apikey     string
		shouldFail bool
		name       string
	}{
		{apikey: user.APIKey, shouldFail: false, name: "Test valid apikey"},
		{apikey: "no-such-key", shouldFail: true, name: "Test invalid apikey"},
	}

	// Test apikey from query param
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s from query param", test.name), func(t *testing.T) {
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			q := req.URL.Query()
			q.Add("apikey", test.apikey)
			req.URL.RawQuery = q.Encode()
			err := handler(c)
			if !test.shouldFail {
				assert.Okf(t, err, "Unable to find valid apikey %s: %s", test.apikey, err)
			} else {
				assert.Errorf(t, err, "Found invalid apikey %s", test.apikey)
			}
		})
	}

	// Test apikey from header
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s from header", test.name), func(t *testing.T) {
			req := httptest.NewRequest(echo.GET, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			req.Header.Set("apikey", test.apikey)
			err := handler(c)
			if !test.shouldFail {
				assert.Okf(t, err, "Unable to find valid apikey %s: %s", test.apikey, err)
			} else {
				assert.Errorf(t, err, "Found invalid apikey %s", test.apikey)
			}
		})
	}
}
