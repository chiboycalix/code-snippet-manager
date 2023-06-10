package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Snippet struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Snippet     string             `json:"snippet,omitempty" validate:"required"`
	Language    string             `json:"language,omitempty" validate:"required"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" validate:"required"`
}
