package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	userGroup := app.Group("/users")
	userGroup.Get("/", handlers.GetAllUsers)
}
