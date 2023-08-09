package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	authGroup := app.Group("/auth")
	authGroup.Post("/register", handlers.RegisterUser)
	authGroup.Post("/login", handlers.LoginUser)
	authGroup.Post("/forget-password", handlers.ForgetPassword)
	authGroup.Post("/reset-password", handlers.ResetPassword)
	authGroup.Post("/verify-account", handlers.VerifyAccount)
	authGroup.Get("/google-signin", handlers.SignInWithGoogle)
	authGroup.Get("/callback", handlers.GoogleCallback)
}
