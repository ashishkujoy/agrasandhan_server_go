package models

type UserRole int

type User struct {
	ID    string   `json:"id" bson:"id"`
	Name  string   `json:"name" bson:"name"`
	Email string   `json:"email" bson:"email"`
	Role  UserRole `json:"role" bson:"role"`
}
