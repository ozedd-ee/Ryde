package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Structure for request used to create SubAccount and TransferRecipient for driver
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
type DriverAccountIDs struct {
	DriverID       primitive.ObjectID `json:"driver_id" bson:"driver_id"`
	SubAccountID             int                `json:"subaccount_id" bson:"subaccount_id"`
	SubAccountCode string             `json:"subaccount_code" bson:"subaccount_code"`
}
