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
	UpdateUser(c echo.Context) error
}

type userController struct {
	uc usecase.IUserUsecase
}

func NewUserController(uc usecase.IUserUsecase) IUserController {
	return &userController{uc: uc}
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
