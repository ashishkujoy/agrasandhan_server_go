package models

type User struct {
	ID    string   `json:"id" bson:"id"`
	Name  string   `json:"name" bson:"name"`
	Email string   `json:"email" bson:"email"`
	Roles []string `json:"roles" bson:"roles"`
}
