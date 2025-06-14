package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *TripService) NewRideRequest(ctx context.Context, riderID, email string, request *models.Order) (*models.TripBuffer, error) {
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
				chargeRequest := models.ChargeRequest{
					Email:  email,
					RideID: trip.ID.String(),
					To:     trip.DriverID.String(),
					Amount: 50000, // Arbitrary number for now
				}
				payment, err := comms.SendChargeRequest(riderID, &chargeRequest)
				if err != nil {
					return nil, fmt.Errorf("error charging card: %v", err)
				}
				fmt.Print(payment.PaymentID)
				return s.TripStore.CacheTrip(ctx, trip)
			} else {
				continue
			}
		}
	}
	return nil, errors.New("sorry, all nearby drivers are currently busy")
}

func (s *TripService) StartTrip(ctx context.Context, tripKey string) (*models.TripBuffer, error) {
	return s.TripStore.StartTrip(ctx, tripKey)
}

func (s *TripService) EndTrip(ctx context.Context, tripKey string) (*models.Trip, error) {
	return s.TripStore.EndTrip(ctx, tripKey)
}

func (s *TripService) GetPendingTrip(ctx context.Context, tripKey string) (*models.Trip, error) {
	return s.TripStore.GetPendingTrip(ctx, tripKey)
}

func (s *TripService) GetTripByID(ctx context.Context, tripID string) (*models.Trip, error) {
	return s.TripStore.GetTripByID(ctx, tripID)
}

func (s *TripService) GetAllDriverTrips(ctx context.Context, driverID string) ([]models.Trip, error) {
	return s.TripStore.GetAllDriverTrips(ctx, driverID)
}

func (s *TripService) GetAllRiderTrips(ctx context.Context, riderID string) ([]models.Trip, error) {
	return s.TripStore.GetAllRiderTrips(ctx, riderID)
}
