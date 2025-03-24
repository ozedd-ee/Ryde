package data

import (
	"context"
	"errors"
	"ryde/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
)

type LocationStore struct {
	RedisClient *redis.Client
}

func NewLocationStore(redisClient *redis.Client) *LocationStore {
	return &LocationStore{
		RedisClient: redisClient,
	}
}

func (s *LocationStore) UpdateDriverLocation(ctx context.Context, driverID string, lat, lon float64) error {
	_, err := s.RedisClient.GeoAdd(ctx, "drivers:location", &redis.GeoLocation{
		Name:      driverID,
		Latitude:  lat,
		Longitude: lon,
	}).Result()
	if err != nil {
		return err
	}
	// Set time to live to 10 minutes
	s.RedisClient.Expire(ctx, "drivers:location", 600*time.Second)
	return nil
}

func (s *LocationStore) FindNearbyDrivers(ctx context.Context, lat, lon, radius float64) ([]string, error) {
	drivers, err := s.RedisClient.GeoRadius(ctx, "drivers:location", lon, lat, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   "km",
		Sort:   "ASC", // Closest first
		Count:  10,    // Limit result to 10 closest drivers
	}).Result()

	if err != nil {
		return nil, err
	}

	var driverIDs []string
	for _, d := range drivers {
		if status := s.RedisClient.Get(ctx, d.Name).String(); status == "available" {
			driverIDs = append(driverIDs, d.Name)
		}
	}
	return driverIDs, nil
}

// For determining drop-off location at end of ride
func (s *LocationStore) GetDriverLocation(ctx context.Context, driverID string) (*models.Location, error) {
	var location models.Location

	position, err := s.RedisClient.GeoPos(ctx, "drivers:location", driverID).Result()
	if err != nil {
		return nil, err
	}

	// Check if location exists
	if len(position) == 0 || position[0] == nil {
		return nil, errors.New("driver location not found")
	}

	location.DriverID = driverID
	location.Latitude = position[0].Latitude
	location.Longitude = position[0].Longitude

	return &location, nil
}
