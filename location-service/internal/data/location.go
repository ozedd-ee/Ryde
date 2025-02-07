package data

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type DataStore struct {
	RedisClient *redis.Client
}

func NewDataStore(redisClient *redis.Client) *DataStore {
	return &DataStore{
		RedisClient: redisClient,
	}
}

func (s *DataStore) UpdateDriverLocation(ctx context.Context, driverID string, lat, lon float64) error {
	_, err := s.RedisClient.GeoAdd(ctx, "drivers:location", &redis.GeoLocation{
		Name: driverID,
		Latitude: lat,
		Longitude: lon,
	}).Result()
	if err != nil {
		return err
	}
	// Set time to live to 10 minutes
	s.RedisClient.Expire(ctx, "drivers:location", 600*time.Second)
	return nil
}

func (s *DataStore) FindNearbyDrivers(ctx context.Context, lat, lon, radius float64) ([]string, error) {
	drivers, err := s.RedisClient.GeoRadius(ctx, "drivers:location", lon, lat, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit: "km",
		Sort: "ASC", // Closest first 
		Count: 10, // Limit result to 10 closest drivers
	}).Result()

	if err != nil {
		return nil, err
	}

	var driverIDs []string
	for _, d := range drivers {
		driverIDs = append(driverIDs, d.Name)
	}
	return driverIDs, nil
}

// For determining drop-off location at end of ride
func (s *DataStore) GetDriverLocation(ctx context.Context, driverID string) (float64, float64, error) {
	position, err := s.RedisClient.GeoPos(ctx, "drivers:location",driverID).Result()
	if err != nil {
		return 0, 0, err
	}

	// Check if location exists
	if len(position) == 0 || position[0] == nil {
		return 0, 0, errors.New("driver location not found")
	}

	return position[0].Latitude, position[0].Longitude, nil
}
