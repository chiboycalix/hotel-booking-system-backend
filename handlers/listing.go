package handlers

import (
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

type CreateListingDTO struct {
	RoomPrice   int64  `json:"roomPrice" bson:"roomPrice" validate:"required"`
	Location    string `json:"location" bson:"location" validate:"required"`
	RoomName    string `json:"roomName" bson:"roomName" validate:"required"`
	RoomBedType string `json:"roomBedType" bson:"roomBedType" validate:"required"`
	RoomImage   string `json:"roomImage" bson:"roomImage"`
}
type GetListingDTO struct {
	ID          string `json:"id" bson:"_id"`
	RoomPrice   int64  `json:"roomPrice" bson:"roomPrice" validate:"required"`
	Location    string `json:"location" bson:"location" validate:"required"`
	RoomName    string `json:"roomName" bson:"roomName" validate:"required"`
	RoomBedType string `json:"roomBedType" bson:"roomBedType" validate:"required"`
	RoomImage   string `json:"roomImage" bson:"roomImage"`
}
type UpdateListingDTO struct {
	RoomPrice   int64  `json:"roomPrice" bson:"roomPrice"`
	Location    string `json:"location" bson:"location"`
	RoomName    string `json:"roomName" bson:"roomName"`
	RoomBedType string `json:"roomBedType" bson:"roomBedType"`
	RoomImage   string `json:"roomImage" bson:"roomImage"`
}

var LISTING_MODEL = "listings"

func CreateListing(c *fiber.Ctx) error {
	var listingCollection = common.GetDBCollection(LISTING_MODEL)
	var createListing CreateListingDTO
	if err := c.BodyParser(&createListing); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": "Please provide request body"}})
	}

	if validationErr := validate.Struct(&createListing); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	formHeader, err := c.FormFile("roomImage")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	uploadUrl, err := utils.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	createListing.RoomImage = uploadUrl
	result, err := listingCollection.InsertOne(c.Context(), createListing)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusCreated).JSON(responses.APIResponse{Status: http.StatusCreated, Message: "Listing created successfully", Data: &fiber.Map{"listing": result}})
}

func GetListing(c *fiber.Ctx) error {
	listingCollection := common.GetDBCollection(LISTING_MODEL)
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	listing := GetListingDTO{}

	err = listingCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&listing)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "Listing fetched successfully", Data: &fiber.Map{"listing": listing}})
}

func GetAllListings(c *fiber.Ctx) error {
	var listingCollection = common.GetDBCollection(LISTING_MODEL)
	listings := make([]GetListingDTO, 0)
	cursor, err := listingCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		listing := GetListingDTO{}
		err := cursor.Decode(&listing)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		listings = append(listings, listing)
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "Listings fetched successfully", Data: &fiber.Map{"listings": listings}})
}

func UpdateListing(c *fiber.Ctx) error {
	listingCollection := common.GetDBCollection(LISTING_MODEL)
	b := new(UpdateListingDTO)
	if err := c.BodyParser(b); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid body"}})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	listing := GetListingDTO{}

	err = listingCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&listing)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	if b.Location == "" {
		b.Location = listing.Location
	}
	if b.RoomBedType == "" {
		b.RoomBedType = listing.RoomBedType
	}
	if b.RoomName == "" {
		b.RoomName = listing.RoomName
	}
	if b.RoomPrice == 0 {
		b.RoomPrice = listing.RoomPrice
	}
	formHeader, err := c.FormFile("roomImage")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	uploadUrl, err := utils.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	b.RoomImage = uploadUrl
	result, err := listingCollection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update listing", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "Listing update was successful", Data: &fiber.Map{"listing": result}})
}

func DeleteListing(c *fiber.Ctx) error {
	id := c.Params("id")
	listingCollection := common.GetDBCollection(LISTING_MODEL)
	if id == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	result, err := listingCollection.DeleteOne(c.Context(), bson.M{"_id": objectId})

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Fail to delete listing"}})
	}
	return c.Status(http.StatusOK).JSON(responses.APIResponse{Status: http.StatusOK, Message: "Listing deleted successfully", Data: &fiber.Map{"data": result}})
}
