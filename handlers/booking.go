package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/chiboycalix/hotel-booking-system-backend/common"
	"github.com/chiboycalix/hotel-booking-system-backend/models"
	"github.com/chiboycalix/hotel-booking-system-backend/responses"
	"github.com/chiboycalix/hotel-booking-system-backend/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var BOOKING_MODEL = "bookings"

type CreateBookingDTO struct {
	RoomID             string    `json:"roomId" bson:"roomId" validate:"required"`
	GuestID            string    `json:"guestId" bson:"guestId"`
	CheckIn            time.Time `json:"checkIn" bson:"checkIn" validate:"required"`
	CheckOut           time.Time `json:"checkOut" bson:"checkOut" validate:"required"`
	BookingDate        time.Time `json:"bookingDate" bson:"bookingDate"`
	BookingUpdatedDate time.Time `json:"bookingUpdatedDate" bson:"bookingUpdatedDate"`
}

type UpdateBookingDTO struct {
	RoomID             string    `json:"roomId" bson:"roomId" validate:"required"`
	CheckIn            time.Time `json:"checkIn" bson:"checkIn" validate:"required"`
	CheckOut           time.Time `json:"checkOut" bson:"checkOut" validate:"required"`
	BookingDate        time.Time `json:"bookingDate" bson:"bookingDate" validate:"required"`
	BookingUpdatedDate time.Time `json:"bookingUpdatedDate" bson:"bookingUpdatedDate"`
}

type GetBookingDTO struct {
	ID                 string    `json:"id" bson:"_id"`
	RoomID             string    `json:"roomId" bson:"roomId"`
	GuestID            string    `json:"guestId" bson:"guestId"`
	CheckIn            time.Time `json:"checkIn" bson:"checkIn"`
	CheckOut           time.Time `json:"checkOut" bson:"checkOut"`
	BookingDate        time.Time `json:"bookingDate" bson:"bookingDate"`
	BookingUpdatedDate time.Time `json:"bookingUpdatedDate" bson:"bookingUpdatedDate"`
}

func CreateBooking(c *fiber.Ctx) error {
	bookingCollection := common.GetDBCollection(BOOKING_MODEL)
	var createBookingDTO CreateBookingDTO
	var header = c.Get("Authorization")
	tokenString := strings.Split(header, " ")[1]

	if err := c.BodyParser(&createBookingDTO); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": "Please provide request body"}})
	}

	if validationErr := validate.Struct(&createBookingDTO); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	userId, err := utils.GetUserID(tokenString, common.EnvJWTSecret(), c)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	createBookingDTO.GuestID = userId
	createBookingDTO.BookingDate = time.Now()
	createBookingDTO.BookingUpdatedDate = time.Now()
	result, err := bookingCollection.InsertOne(c.Context(), createBookingDTO)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusCreated).
		JSON(responses.APIResponse{Status: http.StatusCreated, Message: "Bookings created successfully", Data: &fiber.Map{"booking": result}})
}

func GetAllBookings(c *fiber.Ctx) error {
	bookingCollection := common.GetDBCollection(BOOKING_MODEL)
	bookings := make([]GetBookingDTO, 0)
	cursor, err := bookingCollection.Find(c.Context(), bson.M{})
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	// iterate over the cursor
	for cursor.Next(c.Context()) {
		booking := GetBookingDTO{}
		err := cursor.Decode(&booking)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}
		bookings = append(bookings, booking)
	}

	// Iterate through bookings and populate user and room information
	populatedBooking := []map[string]interface{}{}

	for _, booking := range bookings {
		roomObjectId, _ := primitive.ObjectIDFromHex(booking.RoomID)
		guestObjectId, _ := primitive.ObjectIDFromHex(booking.GuestID)
		room := models.Room{}
		guest := models.User{}
		if err := common.GetDBCollection("rooms").FindOne(context.Background(), bson.M{"_id": roomObjectId}).Decode(&room); err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}

		if err := common.GetDBCollection("users").FindOne(context.Background(), bson.M{"_id": guestObjectId}).Decode(&guest); err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
		}

		// Combine booking, room and user information
		booking := map[string]interface{}{
			"booking": booking,
			"room":    room,
			"guest":   guest,
		}

		populatedBooking = append(populatedBooking, booking)
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "bookings fetched successfully", Data: &fiber.Map{"bookings": populatedBooking}})
}

func UpdateBooking(c *fiber.Ctx) error {
	bookingCollection := common.GetDBCollection(BOOKING_MODEL)
	var b UpdateBookingDTO
	if err := c.BodyParser(&b); err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Id is required"}})
	}
	if validationErr := validate.Struct(&b); validationErr != nil {
		for _, err := range validationErr.(validator.ValidationErrors) {
			return c.Status(http.StatusBadRequest).
				JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Invalid request", Data: &fiber.Map{"error": err.Field() + " is required"}})
		}
	}

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Invalid Id"}})
	}

	booking := GetBookingDTO{}

	err = bookingCollection.FindOne(c.Context(), bson.M{"_id": objectId}).Decode(&booking)
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Something went wrong", Data: &fiber.Map{"error": err.Error()}})
	}

	b.BookingUpdatedDate = time.Now()
	b.BookingDate = booking.BookingDate
	result, err := bookingCollection.UpdateOne(c.Context(), bson.M{"_id": objectId}, bson.M{"$set": b})
	if err != nil {
		return c.Status(http.StatusInternalServerError).
			JSON(responses.APIResponse{Status: http.StatusInternalServerError, Message: "Failed to update booking", Data: &fiber.Map{"error": err.Error()}})
	}

	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Booking update was successful", Data: &fiber.Map{"booking": result}})
}

func DeleteBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	bookingCollection := common.GetDBCollection(BOOKING_MODEL)
	var header = c.Get("Authorization")
	tokenString := strings.Split(header, " ")[1]
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

	result, err := bookingCollection.DeleteOne(c.Context(), bson.M{"_id": objectId})
	if err != nil {
		return c.Status(http.StatusBadRequest).
			JSON(responses.APIResponse{Status: http.StatusBadRequest, Message: "Something went wrong", Data: &fiber.Map{"error": "Fail to delete booking"}})
	}
	return c.Status(http.StatusOK).
		JSON(responses.APIResponse{Status: http.StatusOK, Message: "Booking deleted successfully", Data: &fiber.Map{"data": result}})
}
