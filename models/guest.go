package models

type Guest struct {
	ID        string `json:"id" bson:"_id"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	Phone     string `json:"phone" bson:"phone"`
	RoomID    string `json:"roomId" bson:"roomId"`
}
