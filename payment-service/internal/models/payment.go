package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentID       primitive.ObjectID `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	TripID          string             `json:"trip_id" bson:"trip_id"`
	TransactionRef  string             `json:"tx_ref" bson:"tx_ref"`
	TransactionTime time.Time          `json:"tx_time" bson:"tx_time"`
}
