package service

import (
	"errors"
	"ryde/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createTrip(riderID string, request *models.RideRequest) (*models.Trip, error) {
	driver_id, err := primitive.ObjectIDFromHex(request.DriverID)
	if err != nil {
		return nil, errors.New("invalid driver id format")
	}
	rider_id, err := primitive.ObjectIDFromHex(riderID)
	if err != nil {
		return nil, errors.New("invalid rider id format")
	}
	trip := models.Trip{
		DriverID:    driver_id,
		RiderID:     rider_id,
		Origin:      request.Order.Origin,
		Destination: request.Order.Destination,
		Status:      "initiated",
	}
	return &trip, nil
}
