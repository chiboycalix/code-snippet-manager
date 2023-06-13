package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
}
