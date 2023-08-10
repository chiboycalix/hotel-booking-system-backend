package router

import (
	"github.com/chiboycalix/hotel-booking-system-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func ListingRoutes(app *fiber.App) {
	listingGroup := app.Group("/listings")
	listingGroup.Post("/", handlers.CreateListing)
	listingGroup.Get("/", handlers.GetAllListings)
	listingGroup.Put("/:id", handlers.UpdateListing)
	listingGroup.Delete("/:id", handlers.DeleteListing)
}
