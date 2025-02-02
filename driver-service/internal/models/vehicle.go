package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vehicle struct {
	DriverID       primitive.ObjectID `json:"id" bson:"_id"`
	Brand     string             `json:"brand" bson:"brand"`
	Model    string             `json:"model" bson:"model"`
	Year string             `json:"year" bson:"year"`
	Color string             `json:"color" bson:"color"`
}
