package models

type Room struct {
	ID                string   `json:"id" bson:"_id"`
	RoomImage         string   `json:"roomImage" bson:"roomImage"`
	RoomName          string   `json:"roomName" bson:"roomName"`                   // Deluxe, Suite, etc.
	RoomFacilities    []string `json:"roomFacilities" bson:"roomFacilities"`       // Wifi, AC, TV, etc.
	RoomBookingStatus string   `json:"roomBookingStatus" bson:"roomBookingStatus"` // Available, Booked, etc.
	RoomFloor         int64    `json:"roomFloor" bson:"roomFloor"`                 // 1, 2, 3, etc.
	RoomBlock         string   `json:"roomBlock" bson:"roomBlock"`                 // A, B, C, etc.
	RoomNumber        int64    `json:"roomNumber" bson:"roomNumber"`               // 101, 102, 103, etc.
	RoomCategory      string   `json:"roomCategory" bson:"roomCategory"`           // Single, Double, Triple, etc.
}
