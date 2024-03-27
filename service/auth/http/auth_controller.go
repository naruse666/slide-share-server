package http

import (
	"net/http"
	"slide-share/model"
	"slide-share/service/auth/usecase"

	"github.com/labstack/echo/v4"
)

type IAuthController interface {
	SignIn(c echo.Context) error
}

type authController struct {
	au usecase.IAuthUsecase
}

func NewAuthController(au usecase.IAuthUsecase) IAuthController {
	return &authController{au: au}
}

func (ac *authController) SignIn(c echo.Context) error {
	user := model.SignInUser{}
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	userRes, err := ac.au.SignIn(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, userRes)
}
