package data

import (
	"context"
	"errors"
	"ryde/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverStore struct {
	DriverStore       *mongo.Collection
	DriverStatusCache *redis.Client
}

func NewDriverStore(db *mongo.Database, cache *redis.Client) *DriverStore {
	return &DriverStore{
		DriverStore:       db.Collection("drivers"),
		DriverStatusCache: cache,
	}
}

func (s *DriverStore) CreateDriver(ctx context.Context, driver *models.Driver) (*models.Driver, error) {
	result, err := s.DriverStore.InsertOne(ctx, driver)
	if err != nil {
		return nil, err
	}
	driver.ID = result.InsertedID.(primitive.ObjectID)
	return driver, nil
}

func (s *DriverStore) GetDriver(ctx context.Context, id string) (*models.Driver, error) {
	var driver models.Driver
	driverID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid driver id format")
	}
	filter := bson.M{"_id": driverID}

	if err := s.DriverStore.FindOne(ctx, filter).Decode(&driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (s *DriverStore) GetDriverByEmail(ctx context.Context, email string) (*models.Driver, error) {
	var driver models.Driver

	filter := bson.M{"email": email}
	if err := s.DriverStore.FindOne(ctx, filter).Decode(&driver); err != nil {
		return nil, err
	}
	return &driver, nil
}

func (s *DriverStore) SetStatusBusy(ctx context.Context, driverID string) error {
	if _, err := s.DriverStatusCache.Set(ctx, driverID, "busy", 24 * time.Hour).Result(); err != nil {
		return err
	}
	return nil
}

func (s *DriverStore) SetStatusAvailable(ctx context.Context, driverID string) error {
	if _, err := s.DriverStatusCache.Set(ctx, driverID, "available", 24 * time.Hour).Result(); err != nil {
		return err
	}
	return nil
}

func (s *DriverStore) SetStatusOffline(ctx context.Context, driverID string) error {
	if _, err := s.DriverStatusCache.Set(ctx, driverID, "offline", 24 * time.Hour).Result(); err != nil {
		return err
	}
	return nil
}
