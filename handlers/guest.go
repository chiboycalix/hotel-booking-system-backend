package handlers

import "github.com/gofiber/fiber/v2"

var GUEST_MODEL = "guests"

type CreateGuestDTO struct {
	FirstName string `json:"firstName" bson:"firstName" validate:"required"`
	LastName  string `json:"lastName" bson:"lastName" validate:"required"`
	Email     string `json:"email" bson:"email" validate:"required,email"`
	Phone     string `json:"phone" bson:"phone" validate:"required"`
	RoomID    string `json:"roomId" bson:"roomId"`
}

func CreateGuest(c *fiber.Ctx) error {
	return nil
}
