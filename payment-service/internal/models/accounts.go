package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Structure for request used to create SubAccount and TransferRecipient for driver
type DriverAccountRequest struct {
	BusinessName        string  `json:"business_name"`
	Description         string  `json:"description"`
	PrimaryContactName  string  `json:"primary_contact_name"`
	PrimaryContactEmail string  `json:"primary_contact_email"`
	PrimaryContactPhone string  `json:"primary_contact_phone"`
	PercentageCharge    float32 `json:"percentage_charge"`
	SettlementBank      string  `json:"settlement_bank"`
	AccountNumber       string  `json:"account_number"`
	BankCode            string  `json:"bank_code"`
}

// Merge DriverID with SubAccount and TransferRecipient
type DriverAccountIDs struct {
	DriverID            primitive.ObjectID `json:"driver_id" bson:"driver_id"`
	SubAccountID        int                `json:"subaccount_id" bson:"subaccount_id"`
	TransferRecipientID int                `json:"transfer_recipient_id" bson:"transfer_recipient_id"`
	SubAccountCode      string             `json:"subaccount_code" bson:"subaccount_code"`
	RecipientCode       string             `json:"recipient_code" bson:"recipient_code"`
}
