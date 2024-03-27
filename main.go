package main

import (
	"fmt"
	"os"
	"slide-share/infrastructure/repository/firebase"
	"slide-share/lib"
	auth_http "slide-share/service/auth/http"
	auth_usecase "slide-share/service/auth/usecase"
	user_http "slide-share/service/user/http"
	user_usecase "slide-share/service/user/usecase"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	firestore, _ := lib.InitFirebase()
	userRepository := firebase.NewUserRepository(firestore)

	authUsecase := auth_usecase.NewAuthUsecase(userRepository)
	userUsecase := user_usecase.NewUserUsecase(userRepository)

	authController := auth_http.NewAuthController(authUsecase)
	userController := user_http.NewUserController(userUsecase)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FE_URL")},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))

	auth_http.NewAuthRouter(e, authController)
	user_http.NewUserRouter(e, userController)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", 8080)))
}
