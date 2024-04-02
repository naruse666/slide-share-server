package http

import "github.com/labstack/echo/v4"

func SlideRoutes(g *echo.Group, sc ISlideController) {
	g.GET("", func(c echo.Context) error {
		if page := c.QueryParam("page"); page != "" {
			return sc.GetSlideGroupByPage(c)
		}
		return sc.GetSlideGroups(c)
	})

	g.POST("/upload/slides", sc.UploadSlideBySlidesURL)
	g.POST("/upload/pdf", sc.UploadSlideByPDF)

	g.GET("/newest", sc.GetNewestSlideGroup)

	g.GET("/:slide_group_id", sc.GetSlideGroup)
	g.POST("/:slide_group_id", sc.CreateSlideGroup)

	g.GET("/:slide_group_id/:slide_id", sc.GetSlide)
	g.PUT("/:slide_group_id/:slide_id", sc.UpdateSlide)
}

func NewSlideRouter(e *echo.Echo, sc ISlideController) {
	g := e.Group("/slides")
	SlideRoutes(g, sc)
}
