package http

import "github.com/labstack/echo/v4"

func UserRoutes(g *echo.Group, uc IUserController) {
	g.GET("", uc.GetUsers)
	g.PUT("", uc.UpdateUser)
}

func NewUserRouter(e *echo.Echo, uc IUserController) {
	g := e.Group("/user")
	UserRoutes(g, uc)
}
