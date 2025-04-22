package service

import (
	"context"
	"errors"
	"os"
	"ryde/internal/data"
	"ryde/internal/models"
	"time"

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

func (s *PaymentService) AddPaymentMethod(ctx context.Context, riderID, email string) (string, error) {
	txRequest := paystack.TransactionRequest{
		CallbackURL: "",
		Amount: 100, // 1 Naira
		Email: email,
		Currency: "NGN",
		Metadata: map[string]any{"rider_id": riderID},
	}
	resp, err := s.PaystackClient.Transaction.Initialize(&txRequest)
	if err != nil {
		return "", err
	}
	data, ok := resp["data"].(map[string]any)
	if !ok {
		return "", errors.New("unexpected response format from Paystack")
	}
	authURL, ok := data["authorization_url"].(string)
	if !ok {
		return "", errors.New("authorization url not found in response")
	}
	return authURL, nil
}

// To be triggered in response to 'AddPayMentMethod'
func (s *PaymentService) PaystackCallbackHandler(ctx context.Context, reference string) error {
	transaction, err := s.PaystackClient.Transaction.Verify(reference)
	if err != nil {
		return errors.New("failed to verify transaction")
	}

    if transaction.Status != "success" {
        return errors.New("transaction not successful:" +  transaction.Status)
    }
    auth := transaction.Authorization
    if auth.AuthorizationCode == "" {
        return errors.New("no authorization code found")
    }

	paymentMethod := models.PaymentMethod{
		Email : transaction.Customer.Email,
		AuthCode : auth.AuthorizationCode,
		CardType : auth.CardType,
		Last4 : auth.Last4,
		ExpMonth : auth.ExpMonth,
		ExpYear : auth.ExpYear,
		Bank : auth.Bank,
	}
    return s.PaymentStore.SaveRiderPaymentMethod(ctx, &paymentMethod)
}

func (s *PaymentService) ChargeCard(ctx context.Context, chargeRequest *models.ChargeRequest) (*models.Payment, error) {
	authorizationCode, err := s.PaymentStore.GetAuthorizationCodeByEmail(ctx, chargeRequest.Email)
	if err != nil {
		return nil, err
	}

	chargeReq := paystack.TransactionRequest{
		AuthorizationCode: authorizationCode, //Use saved auth code
		Email:             chargeRequest.Email,
		Amount:            chargeRequest.Amount,
		Metadata:          map[string]any{"ride_id": chargeRequest.RideID},
	}

	tx, err := s.PaystackClient.Transaction.ChargeAuthorization(&chargeReq)
	if err != nil || tx.Status != "Success" {
		return nil, err
	}
	payment := models.Payment{
		PaystackID: tx.ID,
		TripID: tx.Metadata,  // Why is Metadata a string?
		TransactionRef: tx.Reference,
		TransactionTime: time.Now(), // Why is tx.CreatedAt a string? Use time module, for now.
		Amount: tx.Amount,
	}
	newPayment, err := s.PaymentStore.NewPayment(ctx, &payment)
	if err != nil {
		return nil, err
	}
	
	return newPayment, nil
}

func (s *PaymentService) GetSubAccountIDByDriverID(ctx context.Context, driverID string) (*models.SubAccountID, error) {
	return s.PaymentStore.GetSubAccountIDByDriverID(ctx, driverID)	
}
