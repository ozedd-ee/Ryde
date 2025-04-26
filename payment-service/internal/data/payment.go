package data

import (
	"context"
	"errors"
	"ryde/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PaymentStore struct {
	PaymentCollection *mongo.Collection
	SubAccountCollection *mongo.Collection
	PaymentMethodCollection *mongo.Collection
}

func NewPaymentStore(db *mongo.Database) *PaymentStore {
	return &PaymentStore{
		PaymentCollection: db.Collection("payments"),
		SubAccountCollection: db.Collection("subAccounts"),
		PaymentMethodCollection: db.Collection("paymentMethods"),
	}
}

func (s *PaymentStore) NewPayment(ctx context.Context, payment *models.Payment) (*models.Payment, error) {
	result, err := s.PaymentCollection.InsertOne(ctx, payment)
	if err != nil {
		return nil, err
	}
	payment.PaymentID = result.InsertedID.(primitive.ObjectID)
	return payment, nil
}

func (s *PaymentStore) SaveRiderPaymentMethod(ctx context.Context, PaymentMethod *models.PaymentMethod) error {
	_, err := s.PaymentMethodCollection.InsertOne(ctx, PaymentMethod)
	if err != nil {
		return err
	}
	return nil
}

func (s *PaymentStore) GetPayment(ctx context.Context, paymentID string) (*models.Payment, error) {
	var payment models.Payment
	id, err := primitive.ObjectIDFromHex(paymentID)
	if err != nil {
		return nil, errors.New("invalid trip id format")
	}
	filter := bson.M{"payment_id": id}
	if err := s.PaymentCollection.FindOne(ctx, filter).Decode(&payment); err != nil {
		return nil, err
	}
	return &payment, nil
}
