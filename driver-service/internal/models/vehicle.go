package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vehicle struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	DriverID primitive.ObjectID `json:"driverID" bson:"_driverID"`
	Brand    string             `json:"brand" bson:"brand"`
	Model    string             `json:"model" bson:"model"`
	Year     string             `json:"year" bson:"year"`
	Color    string             `json:"color" bson:"color"`
	RegNum   string				`json:"regnum" bson:"regnum"`
}
