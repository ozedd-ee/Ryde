package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	PaymentID       primitive.ObjectID `json:"payment_id,omitempty" bson:"payment_id,omitempty"`
	PaystackID      int                `json:"paystack_id" bson:"paystack_id"`
	TripID          string             `json:"trip_id" bson:"trip_id"`
	TransactionRef  string             `json:"tx_ref" bson:"tx_ref"`
	TransactionTime time.Time          `json:"tx_time" bson:"tx_time"`
	Amount          float32            `json:"amount" bson:"amount"`
}

type PaymentMethod struct {
	Email    string `json:"email" bson:"email"`
	AuthCode string `json:"auth_code" bson:"auth_code"`
	CardType string `json:"card_type" bson:"card_type"`
	Last4    string `json:"last4" bson:"last4"`
	ExpMonth string `json:"exp_month" bson:"exp_month"`
	ExpYear  string `json:"exp_year" bson:"exp_year"`
	Bank     string `json:"bank" bson:"bank"`
}

type ChargeRequest struct {
	Email  string  `json:"email"` // Rider email
	RideID string  `json:"ride_id"`
	To     string  `json:"to"` // Driver ID
	Amount float32 `json:"amount"`
}
