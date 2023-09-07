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

var ROOM_MODEL = "rooms"

type CreateRoomDTO struct {
	RoomImage         string   `json:"roomImage" bson:"roomImage"`
	RoomName          string   `json:"roomName" bson:"roomName" validate:"required"`                   // Deluxe, Suite, etc.
	RoomFacilities    []string `json:"roomFacilities" bson:"roomFacilities" validate:"required"`       // Wifi, AC, TV, etc.
	RoomBookingStatus string   `json:"roomBookingStatus" bson:"roomBookingStatus" validate:"required"` // Available, Booked, etc.
	RoomFloor         int64    `json:"roomFloor" bson:"roomFloor" validate:"required"`                 // 1, 2, 3, etc.
	RoomBlock         string   `json:"roomBlock" bson:"roomBlock" validate:"required"`                 // A, B, C, etc.
	RoomNumber        int64    `json:"roomNumber" bson:"roomNumber" validate:"required"`               // 101, 102, 103, etc.
	RoomCategory      string   `json:"roomCategory" bson:"roomCategory" validate:"required"`           // Single, Double, Triple, etc.
}

type UpdateRoomDTO struct {
	RoomImage         string   `json:"roomImage" bson:"roomImage"`
	RoomName          string   `json:"roomName" bson:"roomName"`                   // Deluxe, Suite, etc.
	RoomFacilities    []string `json:"roomFacilities" bson:"roomFacilities"`       // Wifi, AC, TV, etc.
	RoomBookingStatus string   `json:"roomBookingStatus" bson:"roomBookingStatus"` // Available, Booked, etc.
	RoomFloor         int64    `json:"roomFloor" bson:"roomFloor"`                 // 1, 2, 3, etc.
	RoomBlock         string   `json:"roomBlock" bson:"roomBlock"`                 // A, B, C, etc.
	RoomNumber        int64    `json:"roomNumber" bson:"roomNumber"`               // 101, 102, 103, etc.
	RoomCategory      string   `json:"roomCategory" bson:"roomCategory"`           // Single, Double, Triple, etc.
}

type GetRoomDTO struct {
	ID                string   `json:"id"          bson:"_id"`
	RoomImage         string   `json:"roomImage" bson:"roomImage"`
	RoomName          string   `json:"roomName" bson:"roomName" validate:"required"`                   // Deluxe, Suite, etc.
	RoomFacilities    []string `json:"roomFacilities" bson:"roomFacilities" validate:"required"`       // Wifi, AC, TV, etc.
	RoomBookingStatus string   `json:"roomBookingStatus" bson:"roomBookingStatus" validate:"required"` // Available, Booked, etc.
	RoomFloor         int64    `json:"roomFloor" bson:"roomFloor" validate:"required"`                 // 1, 2, 3, etc.
	RoomBlock         string   `json:"roomBlock" bson:"roomBlock" validate:"required"`                 // A, B, C, etc.
	RoomNumber        int64    `json:"roomNumber" bson:"roomNumber" validate:"required"`               // 101, 102, 103, etc.
	RoomCategory      string   `json:"roomCategory" bson:"roomCategory" validate:"required"`           // Single, Double, Triple, etc.
}

func CreateRoom(c *fiber.Ctx) error {
	roomCollection := common.GetDBCollection(ROOM_MODEL)
	var createRoomDto CreateRoomDTO
	if err := c.BodyParser(&createRoomDto); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": "Please provide request body"}})
	}
	var tokenString = c.Get("Authorization")
	isAdmin, err := utils.IsAdmin(tokenString, common.EnvJWTSecret(), c)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	if !isAdmin {
		return c.Status(http.StatusUnauthorized).
			JSON(responses.APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized", Data: &fiber.Map{"error": "Unauthorized"}})
	}
	if validationErr := validate.Struct(&createRoomDto); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	formHeader, err := c.FormFile("roomImage")
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	formFile, err := formHeader.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	uploadUrl, err := utils.NewMediaUpload().FileUpload(models.File{File: formFile})
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	createRoomDto.RoomImage = uploadUrl
	result, err := roomCollection.InsertOne(c.Context(), createRoomDto)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusCreated).
		JSON(responses.APIResponse{Status: http.StatusCreated, Message: "Room created successfully", Data: &fiber.Map{"room": result}})
}

func GetRoom(c *fiber.Ctx) error {
	roomCollection := common.GetDBCollection(ROOM_MODEL)
	var room models.Room
	if err := roomCollection.FindOne(c.Context(), models.Room{ID: c.Params("id")}).Decode(&room); err != nil {
		return c.Status(http.StatusNotFound).
			JSON(responses.APIResponse{Status: http.StatusNotFound, Message: "Room not found", Data: &fiber.Map{"error": err.Error()}})
	}
	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Room found", Data: &fiber.Map{"room": room}})
}

func GetAllRooms(c *fiber.Ctx) error {
	roomCollection := common.GetDBCollection(ROOM_MODEL)
	rooms := make([]GetRoomDTO, 0)
	cursor, err := roomCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		room := GetRoomDTO{}
		err := cursor.Decode(&room)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		rooms = append(rooms, room)
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Rooms fetched successfully", Data: &fiber.Map{"rooms": rooms}})
}

func UpdateRoom(c *fiber.Ctx) error {
	roomCollection := common.GetDBCollection(ROOM_MODEL)
	var b UpdateRoomDTO
	if err := c.BodyParser(&b); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid body"}})
	}
	var tokenString = c.Get("Authorization")
	isAdmin, err := utils.IsAdmin(tokenString, common.EnvJWTSecret(), c)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	if !isAdmin {
		return c.Status(http.StatusUnauthorized).
			JSON(responses.APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized", Data: &fiber.Map{"error": "Unauthorized"}})
	}
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	room := GetRoomDTO{}

	err = roomCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&room)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	if b.RoomBlock == "" {
		b.RoomBlock = room.RoomBlock
	}
	if b.RoomBookingStatus == "" {
		b.RoomBookingStatus = room.RoomBookingStatus
	}
	if b.RoomName == "" {
		b.RoomName = room.RoomName
	}
	if b.RoomNumber == 0 {
		b.RoomNumber = room.RoomNumber
	}
	if b.RoomFloor == 0 {
		b.RoomFloor = room.RoomFloor
	}
	if b.RoomFacilities == nil {
		b.RoomFacilities = room.RoomFacilities
	}

	if b.RoomImage != "" {
		var url models.Url
		url.Url = b.RoomImage
		uploadUrl, err := utils.NewMediaUpload().RemoteUpload(url)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		b.RoomImage = uploadUrl
	} else {
		formHeader, err := c.FormFile("roomImage")
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		formFile, err := formHeader.Open()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		// if image
		uploadUrl, err := utils.NewMediaUpload().FileUpload(models.File{File: formFile})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		b.RoomImage = uploadUrl
	}

	result, err := roomCollection.UpdateOne(
		c.Context(),
		bson.M{"_id": objectId},
		bson.M{"$set": b},
	)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update room", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Room update was successful", Data: &fiber.Map{"room": result}})
}

func DeleteRoom(c *fiber.Ctx) error {
	roomCollection := common.GetDBCollection(ROOM_MODEL)
	id := c.Params("id")
	var tokenString = c.Get("Authorization")
	isAdmin, err := utils.IsAdmin(tokenString, common.EnvJWTSecret(), c)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}
	if !isAdmin {
		return c.Status(http.StatusUnauthorized).
			JSON(responses.APIResponse{Status: http.StatusUnauthorized, Message: "Unauthorized", Data: &fiber.Map{"error": "Unauthorized"}})
	}
	if id == "" {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	result, err := roomCollection.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Fail to delete room"}})
	}
	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Room deleted successfully", Data: &fiber.Map{"data": result}})
}
