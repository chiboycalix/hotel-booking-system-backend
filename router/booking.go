package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func BookingsRoutes(app *fiber.App) {
	bookingGroup := app.Group("/bookings")
	bookingGroup.Post("/", handlers.CreateBooking)
	bookingGroup.Get("/", handlers.GetAllBookings)
	bookingGroup.Put("/:id", handlers.UpdateBooking)
	bookingGroup.Delete("/:id", handlers.DeleteBooking)
}
