package handlers

import (
	"net/http"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	IsVerified  bool   `json:"isVerified" bson:"isVerified"`
}

type UpdateUserDTO struct {
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	PhoneNumber int64  `json:"phoneNumber" bson:"phoneNumber"`
	Location    string `json:"location" bson:"location"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
}

const USERS_MODEL = "users"

func GetAllUsers(c *fiber.Ctx) error {
	coll := common.GetDBCollection(USERS_MODEL)

	// find all users
	users := make([]UsersDTO, 0)
	cursor, err := coll.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
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
func GetUser(c *fiber.Ctx) error {
	coll := common.GetDBCollection(USERS_MODEL)

	// find the user
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	user := UsersDTO{}

	err = coll.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "User fetched successfully", Data: &fiber.Map{"user": user}})
}
func UpdateUser(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	b := new(UpdateUserDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid body"}})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	result, err := userCollection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to update user", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "User update was successful", Data: &fiber.Map{"user": result}})
}
func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	userCollection := common.GetDBCollection(USERS_MODEL)
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	result, err := userCollection.DeleteOne(c.Context(), bson.M{"_id": objectId})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Fail to delete user"}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "User deleted successfully", Data: &fiber.Map{"data": result}})
}
