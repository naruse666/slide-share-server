package http

import "github.com/labstack/echo/v4"

func SpeakerRoutes(g *echo.Group, sc ISpeakerController) {
	g.GET("", sc.GetSpeakerList)
}

func NewSpeakerRouter(e *echo.Echo, sc ISpeakerController) {
	g := e.Group("/speaker")
	SpeakerRoutes(g, sc)
}
