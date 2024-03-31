package http

import (
	"log"
	"net/http"
	"os"
	"slide-share/model"
	"slide-share/service/slides/usecase"
	"slide-share/utils"

	"github.com/labstack/echo/v4"
)

type ISlideController interface {
	GetNewestSlideGroup(c echo.Context) error
	GetSlideGroup(c echo.Context) error
	GetSlideGroups(c echo.Context) error
	CreateSlideGroup(c echo.Context) error
	GetSlide(c echo.Context) error
}

type SlideController struct {
	su usecase.ISlideUsecase
}

func NewSlideController(su usecase.ISlideUsecase) ISlideController {
	return &SlideController{su: su}
}

func (sc *SlideController) GetNewestSlideGroup(c echo.Context) error {
	slideGroup, err := sc.su.GetNewestSlideGroup()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroup)
}

func (sc *SlideController) GetSlideGroup(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideGroup, err := sc.su.GetSlideGroup(slideGroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroup)
}

func (sc *SlideController) GetSlideGroups(c echo.Context) error {
	slideGroups, err := sc.su.GetSlideGroups()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slideGroups)
}

func (sc *SlideController) CreateSlideGroup(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")
	secret := os.Getenv("AUTH_SECRET")
	payload, err := utils.VerifyAndGetUserClaims(authToken, secret)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, err.Error())
	}
	if payload.Role == "user" || payload.Role == "" {
		return c.JSON(http.StatusForbidden, "Forbidden")
	}

	slideGroup := model.SlideGroup{}
	if err := c.Bind(&slideGroup); err != nil {
		log.Println(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	slideGroupID := c.Param("slide_group_id")
	slideGroup.ID = slideGroupID

	DriveID, err := sc.su.CreateSlideGroup(&slideGroup)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, DriveID)
}

func (sc *SlideController) GetSlide(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideID := c.Param("slide_id")
	slide, err := sc.su.GetSlide(slideGroupID, slideID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, slide)
}
