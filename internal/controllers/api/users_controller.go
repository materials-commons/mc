package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/materials-commons/mc/internal/store"
	"github.com/materials-commons/mc/internal/store/model"
)

type UsersController struct {
	UsersStore *store.UsersStore
}

func (u *UsersController) AddUser(c echo.Context) error {
	var req model.AddUserModel

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

// Login verifies the users password and returns the APIKey
func (u *UsersController) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return err
	}

	user, err := u.UsersStore.GetAndVerifyUser(req.Email, req.Password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, user)
}

func (u *UsersController) ResetAPIKey(c echo.Context) error {
	return nil
}

func (u *UsersController) ResetPassword(c echo.Context) error {
	return nil
}

func (u *UsersController) ModifyUser(c echo.Context) error {
	return nil
}
