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
	AccountStore   *data.AccountStore
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

func (s *PaymentService) AddDriverAccounts(ctx context.Context, driverID string, driverAccounts *models.DriverAccountRequest) (*models.DriverAccountIDs, error) {
	driver_id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, errors.New("invalid driver id format")
	}

	// Unpack request into Paystack's SubAccount type
	subAccountRequest := paystack.SubAccount{
		BusinessName:        driverAccounts.BusinessName,
		Description:         driverAccounts.Description,
		PrimaryContactName:  driverAccounts.PrimaryContactName,
		PrimaryContactEmail: driverAccounts.PrimaryContactEmail,
		PrimaryContactPhone: driverAccounts.PrimaryContactPhone,
		AccountNumber:       driverAccounts.AccountNumber,
		SettlementBank:      driverAccounts.SettlementBank,
		PercentageCharge:    driverAccounts.PercentageCharge,
	}

	// Unpack request into Paystack's TransferRecipient type
	transferRecipientRequest := paystack.TransferRecipient{
		Type:          "nuban",
		Name:          driverAccounts.PrimaryContactName,
		AccountNumber: driverAccounts.AccountNumber,
		BankCode:      driverAccounts.BankCode,
		Currency:      "NGN",
	}

	// Create new SubAccount
	subAccount, err := s.PaystackClient.SubAccount.Create(&subAccountRequest)
	if err != nil {
		return nil, err
	}

	// Create new TransferRecipient
	transferRecipient, err := s.PaystackClient.Transfer.CreateRecipient(&transferRecipientRequest)
	if err != nil {
		return nil, err
	}

	DriverAccountIDs := models.DriverAccountIDs{
		DriverID:            driver_id,
		SubAccountID:        subAccount.ID,
		SubAccountCode:      subAccount.SubAccountCode,
		TransferRecipientID: transferRecipient.ID,
		RecipientCode:       transferRecipient.RecipientCode,
	}
	// Store DriverAccountIDs
	if err := s.AccountStore.StoreDriverAccountIDs(ctx, &DriverAccountIDs); err != nil {
		return nil, err
	}

	// Return the new DriverAccountIDs
	return &DriverAccountIDs, nil
}

func (s *PaymentService) AddPaymentMethod(ctx context.Context, riderID, email string) (string, error) {
	txRequest := paystack.TransactionRequest{
		CallbackURL: "",
		Amount:      100, // 1 Naira
		Email:       email,
		Currency:    "NGN",
		Metadata:    map[string]any{"rider_id": riderID},
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
		return errors.New("transaction not successful:" + transaction.Status)
	}
	auth := transaction.Authorization
	if auth.AuthorizationCode == "" {
		return errors.New("no authorization code found")
	}

	paymentMethod := models.PaymentMethod{
		Email:    transaction.Customer.Email,
		AuthCode: auth.AuthorizationCode,
		CardType: auth.CardType,
		Last4:    auth.Last4,
		ExpMonth: auth.ExpMonth,
		ExpYear:  auth.ExpYear,
		Bank:     auth.Bank,
	}
	return s.PaymentStore.SaveRiderPaymentMethod(ctx, &paymentMethod)
}

func (s *PaymentService) ChargeCard(ctx context.Context, chargeRequest *models.ChargeRequest) (*models.Payment, error) {
	authorizationCode, err := s.PaymentStore.GetAuthorizationCodeByEmail(ctx, chargeRequest.Email)
	if err != nil {
		return nil, err
	}

	driverAccountIDs, err := s.PaymentStore.GetDriverAccountIDsByDriverID(ctx, chargeRequest.To)
	if err != nil {
		return nil, err
	}
	chargeReq := paystack.TransactionRequest{
		AuthorizationCode: authorizationCode, //Use saved auth code
		Email:             chargeRequest.Email,
		Amount:            chargeRequest.Amount,
		SubAccount:        driverAccountIDs.SubAccountCode,
		Bearer:            "subaccount",
		Metadata:          map[string]any{"ride_id": chargeRequest.RideID},
	}

	tx, err := s.PaystackClient.Transaction.ChargeAuthorization(&chargeReq)
	if err != nil || tx.Status != "success" {
		return nil, err
	}
	payment := models.Payment{
		PaystackID:      tx.ID,
		TripID:          tx.Metadata, // Why is Metadata a string?
		TransactionRef:  tx.Reference,
		TransactionTime: time.Now(), // Why is tx.CreatedAt a string?
		Amount:          tx.Amount,
	}
	newPayment, err := s.PaymentStore.NewPayment(ctx, &payment)
	if err != nil {
		return nil, err
	}
	return newPayment, nil
}

func (s *PaymentService) GetSubAccountIDByDriverID(ctx context.Context, driverID string) (*models.DriverAccountIDs, error) {
	return s.PaymentStore.GetDriverAccountIDsByDriverID(ctx, driverID)
}
