package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role int

const (
	Creator  Role = iota
	Analytic Role = iota
)

// User swagger: model User
type User struct {
	ID           primitive.ObjectID `bson:"_id" json:"id"`
	Login        string             `bson:"login" json:"login" validate:"required" example:"test123""`
	Password     string             `bson:"password" json:"password" validate:"required" example:"qwerty"`
	Email        string             `bson:"email" json:"email"`
	FirstName    string             `bson:"first_name" json:"first_name"`
	LastName     string             `bson:"last_name" json:"last_name"`
	Role         Role               `bson:"role" json:"role"`
	CreationDate uint64             `bson:"creation_date" json:"creationDate"`
}
