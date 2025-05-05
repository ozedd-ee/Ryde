package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"ryde/internal/models"
	"ryde/utils"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TripStore struct {
	Collection *mongo.Collection
	Cache      *redis.Client
}

func NewTripStore(db *mongo.Database, cache *redis.Client) *TripStore {
	return &TripStore{
		Collection: db.Collection("trips"),
		Cache:      cache,
	}
}

// ------ cache-only methods ------
func (s *TripStore) CacheTrip(ctx context.Context, trip *models.Trip) (*models.TripBuffer, error) {
	jsonTrip, err := json.Marshal(trip)
	if err != nil {
		return nil, err
	}
	// Generate temporary id for trip
	key := utils.GenerateMD5Key(trip.DriverID.String(), trip.RiderID.String())

	_, err = s.Cache.Set(ctx, key, jsonTrip, 24*time.Hour).Result()
	if err != nil {
		return nil, err
	}
	return &models.TripBuffer{Key: key, Trip: *trip}, nil
}

func (s *TripStore) StartTrip(ctx context.Context, tripKey string) (*models.TripBuffer, error) {
	var tripBuffer models.TripBuffer

	jsonTripBuffer, err := s.Cache.Get(ctx, tripKey).Result()
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(jsonTripBuffer), &tripBuffer)

	// Update the trip status and start time
	tripBuffer.Trip.Status = "en route"
	tripBuffer.Trip.StartTime = time.Now()

	// Update cache
	jsonTrip, err := json.Marshal(tripBuffer)
	if err != nil {
		return nil, err
	}
	_, err = s.Cache.Set(ctx, tripKey, jsonTrip, 24*time.Hour).Result()
	if err != nil {
		return nil, err
	}
	driverID := tripBuffer.Trip.DriverID.String()
	// Publish driver status update for driver service
	s.publishDriverStatusUpdate(ctx, driverID, "busy")

	return &tripBuffer, nil
}

func (s *TripStore) GetPendingTrip(ctx context.Context, tripKey string) (*models.Trip, error) {
	var tripBuffer models.TripBuffer

	jsonTripBuffer, err := s.Cache.Get(ctx, tripKey).Result()
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(jsonTripBuffer), &tripBuffer)
	return &tripBuffer.Trip, nil
}

// ------ Trip database methods ------
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

func (s *TripStore) GetAllDriverTrips(ctx context.Context, driverID string) ([]models.Trip, error) {
	id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, errors.New("invalid driver ID format")
	}
	filter := bson.M{"driver_id": id}
	cursor, err := s.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var trips []models.Trip
	if err := cursor.All(ctx, &trips); err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *TripStore) GetAllRiderTrips(ctx context.Context, riderID string) ([]models.Trip, error) {
	id, err := primitive.ObjectIDFromHex(riderID)
	if err != nil {
		return nil, errors.New("invalid rider ID format")
	}
	filter := bson.M{"rider_id": id}
	cursor, err := s.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var trips []models.Trip
	if err := cursor.All(ctx, &trips); err != nil {
		return nil, err
	}
	return trips, nil
}

func (s *TripStore) EndTrip(ctx context.Context, tripKey string) (*models.Trip, error) {
	var tripBuffer models.TripBuffer

	jsonTripBuffer, err := s.Cache.Get(ctx, tripKey).Result()
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(jsonTripBuffer), &tripBuffer)

	trip := tripBuffer.Trip
	// Update the trip status and start time
	trip.Status = "completed"
	trip.StartTime = time.Now()

	// Store trip in database
	s.newTrip(ctx, &trip)

	// Publish driver status update for driver service
	s.publishDriverStatusUpdate(ctx, trip.DriverID.String(), "available")

	// Remove trip buffer from cache
	s.Cache.Del(ctx, tripKey)
	return &trip, nil
}

// --------- Internal functions ---------
func (s *TripStore) newTrip(ctx context.Context, trip *models.Trip) (*models.Trip, error) {
	result, err := s.Collection.InsertOne(ctx, trip)
	if err != nil {
		return nil, err
	}
	trip.ID = result.InsertedID.(primitive.ObjectID)
	return trip, nil
}

func (s *TripStore) publishDriverStatusUpdate(ctx context.Context, driverID, status string) {
	payload := models.StatusUpdate{
		DriverID: driverID,
		Status:   status,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	}
	if err := s.Cache.Publish(ctx, "driver_status", jsonPayload).Err(); err != nil {
		fmt.Println("Error publishing status update:", err)
	}
}
