package models

type User struct {
	ID          string `json:"id" bson:"_id"`
	Email       string `json:"email" bson:"email"`
	Password    string `json:"password" bson:"password"`
	Role        string `json:"role" bson:"role"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	PhoneNumber int64  `json:"phoneNumber" bson:"phoneNumber"`
	Location    string `json:"location" bson:"location"`
	DateOfBirth string `json:"dateOfBirth" bson:"dateOfBirth"`
}
