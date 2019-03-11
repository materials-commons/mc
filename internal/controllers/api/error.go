package api

import (
	"net/http"

	"github.com/materials-commons/mc/pkg/mc"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
)

func ToHttpError(err error) error {
	if errors.Cause(err) == mc.ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
