package handlers

import (
	"context"
	"net/http"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/chiboycalix/hotel-booking-system-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type createUserDTO struct {
	Email    string `json:"email,omitempty" bson:"email" validate:"required"`
	Password string `json:"password,omitempty" bson:"password" validate:"required"`
}
type loginDTO struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}

var validate = validator.New()

func RegisterUser(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	var u createUserDTO

	if err := c.BodyParser(&u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Error()}})
	}

	if validationErr := validate.Struct(&u); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	pass, hashErr := utils.HashPassword(u.Password)
	if hashErr != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Failed to hash password", Data: &fiber.Map{"error": hashErr.Error()}})
	}
	u.Password = pass
	result, err := userCollection.InsertOne(c.Context(), u)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.UserResonse{Status: http.StatusCreated, Message: "User created successfully", Data: &fiber.Map{"user": result}})
}

func LoginUser(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	var l loginDTO

	if err := c.BodyParser(&l); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid Body", Data: &fiber.Map{"error": err.Error()}})
	}

	if validationErr := validate.Struct(&l); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": l.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	if err := utils.CheckPasswordHash(result.Password, l.Password); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid credentials", Data: &fiber.Map{"error": "Invalid Email or Password"}})
	}

	jwt, err := utils.GenerateJWT(result.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to generate jwt", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Login successful", Data: &fiber.Map{"token": jwt}})
}
