package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/chiboycalix/hotel-booking-system-backend/utils"
)

type createUserDTO struct {
	Email       string `json:"email,omitempty"    bson:"email"       validate:"required"`
	Password    string `json:"password,omitempty" bson:"password"    validate:"required"`
	Role        string `json:"role"               bson:"role"`
	FirstName   string `json:"firstName"          bson:"firstName"`
	LastName    string `json:"lastName"           bson:"lastName"`
	PhoneNumber int64  `json:"phoneNumber"        bson:"phoneNumber"`
	Location    string `json:"location"           bson:"location"`
	DateOfBirth string `json:"dateOfBirth"        bson:"dateOfBirth"`
	IsVerified  bool   `json:"isVerified"         bson:"isVerified"`
	IsAdmin     bool   `json:"isAdmin"            bson:"isAdmin"`
}
type loginDTO struct {
	Email    string `json:"email,omitempty"    validate:"required"`
	Password string `json:"password,omitempty" validate:"required"`
}
type forgotPasswordDTO struct {
	Email string `json:"email,omitempty" validate:"required"`
}
type resetPasswordDTO struct {
	Email      string `json:"email,omitempty"    validate:"required"`
	Password   string `json:"password,omitempty" validate:"required"`
	IsVerified bool   `json:"isVerified"                             bson:"isVerified"`
}

type verifyUserDTO struct {
	Email      string `json:"email,omitempty" validate:"required"`
	IsVerified bool   `json:"isVerified"                          bson:"isVerified"`
}

var (
	validate     = validator.New()
	ssoOAuth     *oauth2.Config
	RandomString = "random-text"
)

func RegisterUser(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	var u createUserDTO
	u.IsVerified = true
	u.Role = "GUEST"
	if err := c.BodyParser(&u); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": "Please provide request body"}})
	}
	// set is verified to true for newly registered user

	if validationErr := validate.Struct(&u); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	var r models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&r)
	if err == nil {
		return c.Status(http.StatusConflict).
			JSON(responses.APIResponse{Status: http.StatusConflict, Message: "Already Exist", Data: &fiber.Map{"error": "User with this Email already exist"}})
	}

	pass, hashErr := utils.HashPassword(u.Password)
	if hashErr != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Failed to hash password", Data: &fiber.Map{"error": hashErr.Error()}})
	}
	u.Password = pass
	u.IsAdmin = false
	result, err := userCollection.InsertOne(c.Context(), u)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	// convert interface to string
	jwt, err := utils.GenerateJWT(fmt.Sprint(result.InsertedID), u.IsAdmin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to generate jwt", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusCreated).
		JSON(responses.APIResponse{Status: http.StatusCreated, Message: "User created successfully", Data: &fiber.Map{"user": result, "token": jwt}})
}

func LoginUser(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	var l loginDTO

	if err := c.BodyParser(&l); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid Body", Data: &fiber.Map{"error": err.Error()}})
	}

	if validationErr := validate.Struct(&l); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": l.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	if err := utils.CheckPasswordHash(result.Password, l.Password); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid credentials", Data: &fiber.Map{"error": "Invalid Email or Password"}})
	}

	if !result.IsVerified {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Not Verified", Data: &fiber.Map{"error": "User is not Verified"}})
	}

	jwt, err := utils.GenerateJWT(result.ID, result.IsAdmin)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to generate jwt", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Login successful", Data: &fiber.Map{"token": jwt}})
}

func ForgetPassword(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	var f forgotPasswordDTO
	if err := c.BodyParser(&f); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	if validationErr := validate.Struct(&f); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": f.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": "User with this email not found"}})
	}

	// send email
	err = utils.SendMailService(result, "templates/forget-password.html", "Forget Password")
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Error sending mail", Data: &fiber.Map{"error": "Something went wrong, error sending mail"}})
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Please check your mail for further instructions"})
}

func ResetPassword(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	var r resetPasswordDTO
	if err := c.BodyParser(&r); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	if validationErr := validate.Struct(&r); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}
	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": r.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": "User not found"}})
	}

	objectId, err := primitive.ObjectIDFromHex(result.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	// set it to false so that the user must use the verify link to change it back to true before he/she can login
	r.IsVerified = false
	updateReq, err := userCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": objectId},
		bson.M{"$set": r},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update user", Data: &fiber.Map{"error": err.Error()}})
	}

	// send passowrd changed email
	err = utils.SendMailService(result, "templates/password-changed.html", "Password Changed")
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Error sending mail", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Password reset was successful", Data: &fiber.Map{"user": updateReq}})
}

func VerifyAccount(c *fiber.Ctx) error {
	userCollection := common.GetDBCollection(USERS_MODEL)
	var v verifyUserDTO

	if err := c.BodyParser(&v); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid body", Data: &fiber.Map{"error": err.Error()}})
	}
	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": v.Email}).Decode(&result)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "User not found", Data: &fiber.Map{"error": err.Error()}})
	}

	objectId, err := primitive.ObjectIDFromHex(result.ID)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	v.IsVerified = true
	updateReq, err := userCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": objectId},
		bson.M{"$set": v},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update user", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Your account has been verified", Data: &fiber.Map{"user": updateReq}})
}

func SignInWithGoogle(c *fiber.Ctx) error {
	ssoOAuth = &oauth2.Config{
		RedirectURL:  common.GoogleRedirectURI(),
		ClientID:     common.GoogleClientID(),
		ClientSecret: common.GoogleClientSecret(),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	url := ssoOAuth.AuthCodeURL(RandomString)
	fmt.Println(url)
	return c.Redirect(url, 302)
}

func GoogleCallback(c *fiber.Ctx) error {
	state := c.Params("state")
	code := c.Params("code")
	data, err := getUserData(state, code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Error getting User data", Data: &fiber.Map{"error": "Something went wrong, getting User data"}})
	}
	fmt.Println(data, "data")
	return nil
}

func getUserData(state, code string) ([]byte, error) {
	if state != RandomString {
		fmt.Println("Invalid state")
		return nil, errors.New("invalid state")
	}
	token, err := ssoOAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	response, err := http.Get(
		"https://googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken,
	)
	if err != nil {
		return nil, err
	}
	fmt.Println(response)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
