package service

import (
	"context"
	"errors"
	"ryde/internal/comms"
	"ryde/internal/data"
	"ryde/internal/models"
	"time"
)

type TripService struct {
	TripStore *data.TripStore
}

func NewTripService(tripStore *data.TripStore) *TripService {
	return &TripService{
		TripStore: tripStore,
	}
}

func (s *TripService) NewRideRequest(ctx context.Context, riderID string, request *models.Order) (*models.Trip, error) {
	drivers, err := comms.FindNearbyDrivers(request)
	if err != nil {
		return nil, err
	}
	for _, driver := range drivers {
		timeout := time.After(10 * time.Second)
		select {
		case <-timeout:
			// Skip to next driver if no response after 10 seconds
			continue
		default:
			request, err := comms.NotifyDriver(driver, request)
			if err != nil {
				continue // Skip to next driver if error occurs during notification
			}
			status := request.Order.Status
			if status == "accepted" {
				trip, err := createTrip(riderID, request)
				if err != nil {
					return nil, err
				}
				return s.TripStore.NewTrip(ctx, trip)
			} else {
				continue
			}
		}
	}
	return nil, errors.New("sorry, all nearby drivers are currently busy")
}

func (s *TripService) GetTripByID(ctx context.Context, tripID string) (*models.Trip, error) {
	return s.TripStore.GetTripByID(ctx, tripID)
}

func (s *TripService) GetTripByDriver(ctx context.Context, driverID string) (*models.Trip, error) {
	return s.TripStore.GetTripByDriver(ctx, driverID)
}

func (s *TripService) StartTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	return s.TripStore.StartTrip(ctx, tripID)
}

func (s *TripService) EndTrip(ctx context.Context, tripID string) (*models.Trip, error) {
	return s.TripStore.EndTrip(ctx, tripID)
}
