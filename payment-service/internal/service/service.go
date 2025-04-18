package service

import (
	"context"
	"errors"
	"os"
	"ryde/internal/data"
	"ryde/internal/models"

	"github.com/rpip/paystack-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	PaymentStore   *data.PaymentStore
	PaystackClient *paystack.Client
}

func NewPaymentService(paymentStore *data.PaymentStore) *PaymentService {
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY") // Just for tests
	paystackClient := paystack.NewClient(secretKey, nil)

	return &PaymentService{
		PaymentStore:   paymentStore,
		PaystackClient: paystackClient,
	}
}

func (s *PaymentService) CreateSubAccount(ctx context.Context, driverID string, subaccount *models.SubAccountRequest) (*models.SubAccountID, error) {
	driver_id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, errors.New("invalid driver id format")
	}

	// Unpack request into Paystack's SubAccount type
	subAccountRequest := paystack.SubAccount{
		BusinessName:        subaccount.BusinessName,
		Description:         subaccount.Description,
		PrimaryContactName:  subaccount.PrimaryContactName,
		PrimaryContactEmail: subaccount.PrimaryContactEmail,
		PrimaryContactPhone: subaccount.PrimaryContactPhone,
		AccountNumber:       subaccount.AccountNumber,
		SettlementBank:      subaccount.SettlementBank,
		PercentageCharge:    subaccount.PercentageCharge,
	}

	// Create new SubAccount
	result, err := s.PaystackClient.SubAccount.Create(&subAccountRequest)
	if err != nil {
		return nil, err
	}

	//Store driver's new SubAccount ID and code
	subAccountID := models.SubAccountID{
		DriverID: driver_id,
		ID: result.ID,
		SubAccountCode: result.SubAccountCode,
	}
	if err := s.PaymentStore.StoreSubAccountID(ctx, &subAccountID); err != nil {
		return nil, err
	}

	// Return the new SubAccount's ID
	return &subAccountID, nil
}

func (s *PaymentService) GetSubAccountIDByDriverID(ctx context.Context, driverID string) (*models.SubAccountID, error) {
	return s.PaymentStore.GetSubAccountIDByDriverID(ctx, driverID)	
}
