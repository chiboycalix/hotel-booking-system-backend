package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/chiboycalix/hotel-booking-system-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createUserDTO struct {
	Email      string `json:"email,omitempty" bson:"email" validate:"required"`
	Password   string `json:"password,omitempty" bson:"password" validate:"required"`
	IsVerified bool   `json:"isVerified" bson:"isVerified"`
}
type loginDTO struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}
type forgotPasswordDTO struct {
	Email string `json:"email,omitempty" validate:"required"`
}
type resetPasswordDTO struct {
	Email      string `json:"email,omitempty" validate:"required"`
	Password   string `json:"password,omitempty" validate:"required"`
	IsVerified bool   `json:"isVerified" bson:"isVerified"`
}

type verifyUserDTO struct {
	Email      string `json:"email,omitempty" validate:"required"`
	IsVerified bool   `json:"isVerified" bson:"isVerified"`
}

var validate = validator.New()

func RegisterUser(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	var u createUserDTO
	if err := c.BodyParser(&u); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Error()}})
	}
	// set is verified to true for newly registered user
	u.IsVerified = true
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
	// convert interface to string
	jwt, err := utils.GenerateJWT(fmt.Sprint(result.InsertedID))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to generate jwt", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResonse{Status: http.StatusCreated, Message: "User created successfully", Data: &fiber.Map{"user": result, "token": jwt}})
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

	if !result.IsVerified {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Not Verified", Data: &fiber.Map{"error": "User is not Verified"}})
	}

	jwt, err := utils.GenerateJWT(result.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to generate jwt", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Login successful", Data: &fiber.Map{"token": jwt}})
}

func ForgetPassword(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	var f forgotPasswordDTO
	if err := c.BodyParser(&f); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	if validationErr := validate.Struct(&f); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": f.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	// send email
	err = utils.SendMailService(result, "templates/forget-password.html")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Error sending mail", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Please check your mail for further instructions"})
}

func ResetPassword(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	// i think it will be better to get the email from url parameter rather than from the body
	// email := c.Params("email")
	var r resetPasswordDTO
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	if validationErr := validate.Struct(&r); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}
	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": r.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	objectId, err := primitive.ObjectIDFromHex(result.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	// set it to false so that the user must use the verify link to change it back to true before he/she can login
	r.IsVerified = false
	updateReq, err := userCollection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": r})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to update user", Data: &fiber.Map{"error": err.Error()}})
	}

	// send passowrd changed email
	err = utils.SendMailService(result, "templates/password-changed.html")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Error sending mail", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Password reset was successful", Data: &fiber.Map{"user": updateReq}})
}

func VerifyUser(c *fiber.Ctx) error {
	var userCollection = common.GetDBCollection("users")
	// i think it will be better to get the email from url parameter rather than from the body
	// email := c.Params("email")
	var v verifyUserDTO

	if err := c.BodyParser(&v); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": v.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	objectId, err := primitive.ObjectIDFromHex(result.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResonse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	v.IsVerified = true
	updateReq, err := userCollection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": v})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResonse{Status: http.StatusInternalServerError, Message: "Failed to update user", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResonse{Status: http.StatusOK, Message: "Password reset was successful", Data: &fiber.Map{"user": updateReq}})
}
