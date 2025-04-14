package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type SubAccountRequest struct {
	BusinessName        string  `json:"business_name,omitempty"`
	Description         string  `json:"description,omitempty"`
	PrimaryContactName  string  `json:"primary_contact_name,omitempty"`
	PrimaryContactEmail string  `json:"primary_contact_email,omitempty"`
	PrimaryContactPhone string  `json:"primary_contact_phone,omitempty"`
	PercentageCharge    float32 `json:"percentage_charge,omitempty"`
	SettlementBank      string  `json:"settlement_bank,omitempty"`
	AccountNumber       string  `json:"account_number,omitempty"`
}

// Pair DriverID with SubAccount
type SubAccountID struct {
	DriverID       primitive.ObjectID `json:"driver_id" bson:"driver_id"`
	ID             int                `json:"id" bson:"id"` // SubAccount ID
	SubAccountCode string             `json:"subaccount_code" bson:"subaccount_code"`
}
