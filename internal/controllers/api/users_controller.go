package api

import (
	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
	"net/http"
)

type UsersController struct {
	UsersStore *store.UsersStore
}

func (u *UsersController) AddUser(c echo.Context) error {
	var req store.AddUserModel

	if err := c.Bind(&req); err != nil {
		return err
	}

	createdUser, err := u.UsersStore.AddUser(req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, createdUser)
}

func (u *UsersController) GetUserByID(c echo.Context) error {
	var req struct {
		ID string `json:"id"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := u.UsersStore.GetUserByID(req.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UsersController) GetUserByAPIKey(c echo.Context) error {
	var req struct {
		APIKey string `json:"apikey"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := u.UsersStore.GetUserByAPIKey(req.APIKey)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}
