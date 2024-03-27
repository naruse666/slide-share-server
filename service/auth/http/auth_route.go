package http

import "github.com/labstack/echo/v4"

func AuthRoutes(g *echo.Group, ac IAuthController) {
	g.POST("/signin", ac.SignIn)
}

func NewAuthRouter(e *echo.Echo, ac IAuthController) {
	g := e.Group("/auth")
	AuthRoutes(g, ac)
}
