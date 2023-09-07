package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func RoomRoutes(app *fiber.App) {
	listingGroup := app.Group("/rooms")
	listingGroup.Post("/", handlers.CreateRoom)
	listingGroup.Get("/", handlers.GetAllRooms)
	listingGroup.Get("/:id", handlers.GetRoom)
	listingGroup.Put("/:id", handlers.UpdateRoom)
	listingGroup.Delete("/:id", handlers.DeleteRoom)
}
