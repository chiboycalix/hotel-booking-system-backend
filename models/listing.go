package models

type Listing struct {
	ID          string `json:"id" bson:"_id"`
	Location    string `json:"location" bson:"location"`
	RoomName    string `json:"roomName" bson:"roomName"`
	RoomPrice   int64  `json:"roomPrice" bson:"roomPrice"`
	RoomImage   string `json:"roomImage" bson:"roomImage"`
	RoomBedType string `json:"roomBedType" bson:"roomBedType"`
}
