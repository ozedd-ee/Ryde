package data

import (
	"context"
	"ryde/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AccountStore struct {
	DriverAccountCollection *mongo.Collection
	PaymentMethodCollection *mongo.Collection
}

func NewAccountStore(db *mongo.Database) *AccountStore {
	return &AccountStore{
		DriverAccountCollection: db.Collection("driver_accounts"),
		PaymentMethodCollection: db.Collection("payment_methods"),
	}
}

func (s *AccountStore) StoreDriverAccountIDs(ctx context.Context, driverAccountIDs *models.DriverAccountIDs) error {
	_, err := s.DriverAccountCollection.InsertOne(ctx, driverAccountIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *PaymentStore) GetDriverAccountIDsByDriverID(ctx context.Context, driverID string) (*models.DriverAccountIDs, error) {
	var driverAccountIDs models.DriverAccountIDs
	driver_id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"driver_id": driver_id}
	if err := s.SubAccountCollection.FindOne(ctx, filter).Decode(&driverAccountIDs); err != nil {
		return nil, err
	}
	return &driverAccountIDs, nil
}

func (s *PaymentStore) GetAuthorizationCodeByEmail(ctx context.Context, email string) (string, error) {
	var paymentMethod *models.PaymentMethod
	filter := bson.M{"email": email}
	if err := s.PaymentMethodCollection.FindOne(ctx, filter).Decode(&paymentMethod); err != nil {
		return "", err
	}
	return paymentMethod.AuthCode, nil
}
