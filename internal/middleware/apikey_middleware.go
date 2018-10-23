package middleware

import (
	"fmt"
	"net/http"

	"github.com/materials-commons/mc/internal/store/model"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type APIKeyRetriever func(string, echo.Context) (*model.UserSchema, error)

type APIKeyConfig struct {
	Skipper   middleware.Skipper
	Keyname   string
	Retriever APIKeyRetriever
}

func APIKeyAuth(config APIKeyConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			value, err := getAPIKey(config.Keyname, c)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, err.Error())
			}

			user, err := config.Retriever(value, c)
			switch {
			case err != nil:
				return err
			case user == nil:
				return echo.ErrUnauthorized
			default:
				c.Set("User", *user)
				return next(c)
			}
		}
	}
}

func getAPIKey(key string, c echo.Context) (string, error) {
	if value, err := keyFromHeader(key, c); err == nil {
		return value, nil
	}

	if value, err := keyFromQuery(key, c); err == nil {
		return value, nil
	}

	return "", fmt.Errorf("no apikey '%s' as query param or header", key)
}

func keyFromHeader(key string, c echo.Context) (string, error) {
	value := c.Request().Header.Get(key)
	if value == "" {
		return "", fmt.Errorf("no apikey '%s' as header", key)
	}
	return value, nil
}

func keyFromQuery(key string, c echo.Context) (string, error) {
	value := c.QueryParam(key)
	if value == "" {
		return "", fmt.Errorf("no apikey '%s' as query param", key)
	}
	return value, nil
}
