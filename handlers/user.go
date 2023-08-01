package handlers

import (
	"net/http"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type UsersDTO struct {
	ID          string `json:"id" bson:"_id"`
	Email       string `json:"email" bson:"email"`
	Role        string `json:"role" bson:"role"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	PhoneNumber int64  `json:"phoneNumber" bson:"phoneNumber"`
	Location    string `json:"location" bson:"location"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
}

func GetAllUsers(c *fiber.Ctx) error {
	coll := common.GetDBCollection("users")

	// find all users
	users := make([]UsersDTO, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		user := UsersDTO{}
		err := cursor.Decode(&user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		users = append(users, user)
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Users fetched successfully", Data: &fiber.Map{"users": users}})
}
