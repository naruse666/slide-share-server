package http

import (
	"slide-share/service/slides/usecase"

	"github.com/labstack/echo/v4"
)

type ISlideController interface {
	GetNewestSlideGroup(c echo.Context) error
	GetSlideGroup(c echo.Context) error
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
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, slideGroup)
}

func (sc *SlideController) GetSlideGroup(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideGroup, err := sc.su.GetSlideGroup(slideGroupID)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, slideGroup)
}

func (sc *SlideController) GetSlide(c echo.Context) error {
	slideGroupID := c.Param("slide_group_id")
	slideID := c.Param("slide_id")
	slide, err := sc.su.GetSlide(slideGroupID, slideID)
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return c.JSON(200, slide)
}
