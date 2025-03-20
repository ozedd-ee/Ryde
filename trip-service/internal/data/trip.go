package data

import (
	"context"
	"errors"
	"ryde/internal/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TripStore struct {
	Collection *mongo.Collection
}

func NewTripStore(db *mongo.Database) *TripStore {
	return &TripStore{
		Collection: db.Collection("trips"),
	}
}

func (s *TripStore) NewTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) {
	result, err := s.Collection.InsertOne(ctx, trip)
	if err != nil {
		return nil, err
	}
	trip.ID = result.InsertedID.(primitive.ObjectID)
	return trip, nil
}

func (s *TripStore) GetTripByID(ctx context.Context, tripID string) (*models.Trip, error) {
	var trip models.Trip
	id, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, errors.New("invalid trip ID format")
	}
	filter := bson.M{"_id": id}
	if err = s.Collection.FindOne(ctx, filter).Decode(&trip); err != nil {
		return nil, err
	}
	return &trip, nil
}

func (s *TripStore) GetTripByDriver(ctx context.Context, driverID string) (*models.Trip, error) {
	var trip models.Trip
	id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, errors.New("invalid driver ID format")
	}
	filter := bson.M{"driver_id": id}
	if err = s.Collection.FindOne(ctx, filter).Decode(&trip); err != nil {
		return nil, err
	}
	return &trip, nil
}

func (s *TripStore) StartTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	var updatedTrip models.Trip
	id, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	// Update the trip status and start time
	update := bson.M{"$set": bson.M{"Status": "en route", "StartTime": time.Now()}}
	// Modify options to return updated trip
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = s.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTrip); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("trip not found")
		}
		return nil, err
	}
	return &updatedTrip, nil
}

func (s *TripStore) EndTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	var updatedTrip models.Trip
	id, err := primitive.ObjectIDFromHex(tripID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": id}
	// Update the trip status and start time
	update := bson.M{"$set": bson.M{"Status": "completed", "EndTime": time.Now()}}
	// Modify options to return updated trip
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	if err = s.Collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedTrip); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("trip not found")
		}
		return nil, err
	}
	return &updatedTrip, nil
}
