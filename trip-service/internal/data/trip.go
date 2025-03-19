package data

import (
	"context"
	"ryde/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripStore struct {
	collection *mongo.Collection
}

func NewTripStore(db *mongo.Database) *TripStore {
	return &TripStore{
		collection: db.Collection("trips"),
	}
}

func (s *TripStore) NewTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) {
	result, err := s.collection.InsertOne(ctx, trip)
	if err != nil {
		return nil, err
	}
	trip.ID = result.InsertedID.(primitive.ObjectID)
	return trip, nil
}
