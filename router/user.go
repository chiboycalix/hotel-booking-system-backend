package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	userGroup := app.Group("/users")
	userGroup.Get("/", handlers.GetAllUsers)
	userGroup.Get("/:id", handlers.GetUser)
	userGroup.Put("/:id", handlers.UpdateUser)
	userGroup.Delete("/:id", handlers.DeleteUser)
}
