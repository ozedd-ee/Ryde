package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Order struct {
	Origin      Location `json:"origin"`
	Destination Location `json:"destination"`
	Status      string   `json:"status"` // can be pending, accepted or rejected
}

type RideRequest struct {
	DriverID string `json:"driver_id"`
	Order    Order  `json:"order"`
}

type ChargeRequest struct {
	Email  string  `json:"email"` // Rider email
	RideID string  `json:"ride_id"`
	To     string  `json:"to"` // Driver ID
	Amount float32 `json:"amount"`
}

type Payment struct {
	PaymentID       primitive.ObjectID `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	TransactionID   int                `json:"transaction_id" bson:"transaction_id"`
	TripID          string             `json:"trip_id" bson:"trip_id"`
	TransactionRef  string             `json:"tx_ref" bson:"tx_ref"`
	TransactionTime time.Time          `json:"tx_time" bson:"tx_time"`
	Amount          float32            `json:"amount" bson:"amount"`
}
