package models

type Booking struct {
	ID          string `json:"id" bson:"_id"`
	RoomID      Room   `json:"roomId" bson:"roomId"`
	GuestID     User   `json:"guestId" bson:"guestId"`
	CheckIn     string `json:"checkIn" bson:"checkIn"`
	CheckOut    string `json:"checkOut" bson:"checkOut"`
	BookingDate string `json:"bookingDate" bson:"bookingDate"`
}
