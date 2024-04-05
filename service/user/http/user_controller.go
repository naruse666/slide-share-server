package http

import (
	"net/http"
	"os"
	"slide-share/model"
	"slide-share/service/user/usecase"
	"slide-share/utils"

	"github.com/labstack/echo/v4"
)

type IUserController interface {
	GetUser(c echo.Context) error
	GetUsers(c echo.Context) error
	UpdateUser(c echo.Context) error
}

type userController struct {
	uc usecase.IUserUsecase
}

func NewUserController(uc usecase.IUserUsecase) IUserController {
	return &userController{uc: uc}
}

func (uc *userController) GetUser(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	id := c.Param("id")
	user, err := uc.uc.GetUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

func (uc *userController) GetUsers(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	users, err := uc.uc.GetUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

func (uc *userController) UpdateUser(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	user := model.User{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := uc.uc.UpdateUser(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}
