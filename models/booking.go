package models

import "time"

type Booking struct {
	ID                 string    `json:"id" bson:"_id"`
	RoomID             Room      `json:"roomId" bson:"roomId"`
	GuestID            User      `json:"guestId" bson:"guestId"`
	CheckIn            time.Time `json:"checkIn" bson:"checkIn"`
	CheckOut           time.Time `json:"checkOut" bson:"checkOut"`
	BookingDate        time.Time `json:"bookingDate" bson:"bookingDate"`
	BookingUpdatedDate time.Time `json:"bookingUpdatedDate" bson:"bookingUpdatedDate"`
}
