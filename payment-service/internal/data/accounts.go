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
}

func NewAccountStore(db *mongo.Database) *AccountStore {
	return &AccountStore{
		DriverAccountCollection: db.Collection("driver_accounts"),
	}
}

func (s *AccountStore) StoreDriverAccountIDs(ctx context.Context, driverAccountIDs *models.DriverAccountIDs) error {
	_, err := s.DriverAccountCollection.InsertOne(ctx, driverAccountIDs)
	if err != nil {
		return err
	}
	return nil
}

func (s *AccountStore) GetDriverAccountIDsByDriverID(ctx context.Context, driverID string) (*models.DriverAccountIDs, error) {
	var driverAccountIDs models.DriverAccountIDs
	driver_id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"driver_id": driver_id}
	if err := s.DriverAccountCollection.FindOne(ctx, filter).Decode(&driverAccountIDs); err != nil {
		return nil, err
	}
	return &driverAccountIDs, nil
}
