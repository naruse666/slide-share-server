package http

import (
	"net/http"
	"slide-share/service/speaker/usecase"

	"github.com/labstack/echo/v4"
)

type ISpeakerController interface {
	GetSpeakerByID(c echo.Context) error
	GetSpeakerList(c echo.Context) error
}

type speakerController struct {
	sc usecase.ISpeakerUsecase
}

func NewSpeakerController(sc usecase.ISpeakerUsecase) ISpeakerController {
	return &speakerController{sc: sc}
}

func (sc *speakerController) GetSpeakerByID(c echo.Context) error {
	speakerID := c.Param("speaker_id")
	speaker, err := sc.sc.GetSpeakerByID(speakerID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, speaker)
}

func (sc *speakerController) GetSpeakerList(c echo.Context) error {
	speakers, err := sc.sc.GetSpeakerList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, speakers)
}
