package http

import "github.com/labstack/echo/v4"

func SlideRoutes(g *echo.Group, sc ISlideController) {
	g.GET("", func(c echo.Context) error {
		if page := c.QueryParam("page"); page != "" {
			// page パラメーターが存在する場合に呼び出すメソッド
			// return sc.GetSlidePage(c)
		}
		return sc.GetNewestSlideGroup(c)
	})
	g.GET("/:slide_group_id", sc.GetSlideGroup)
	g.GET("/:slide_group_id/:slide_id", sc.GetSlide)
}

func NewSlideRouter(e *echo.Echo, sc ISlideController) {
	g := e.Group("/slides")
	SlideRoutes(g, sc)
}
