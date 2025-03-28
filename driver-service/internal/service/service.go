package service

import (
	"context"
	"errors"
	"ryde/internal/data"
	"ryde/internal/models"
	"ryde/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DriverService struct {
	DriverStore  *data.DriverStore
	VehicleStore *data.VehicleStore
}

func NewDriverService(driverStore *data.DriverStore, vehicleStore *data.VehicleStore) *DriverService {
	return &DriverService{
		DriverStore:  driverStore,
		VehicleStore: vehicleStore,
	}
}

func (s *DriverService) SignUp(ctx context.Context, driver *models.Driver) (*models.Driver, error) {
	exists, err := s.DriverStore.GetDriverByEmail(ctx, driver.Email)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if exists != nil {
		return nil, errors.New("driver already exists")
	}

	hashedPassword, err := utils.HashPassword(driver.Password)
	if err != nil {
		return nil, err
	}
	driver.Password = hashedPassword
	return s.DriverStore.CreateDriver(ctx, driver)
}

func (s *DriverService) Login(ctx context.Context, email, password string) (string, error) {
	driver, err := s.DriverStore.GetDriverByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	err = utils.CheckPasswordMatch(password, driver.Password)
	if err != nil {
		return "", errors.New("incorrect password")
	}

	token, err := utils.GenerateJWT(driver.ID.Hex())
	if err != nil {
		return "", errors.New("unable to generate JWT")
	}
	return token, nil
}

func (s *DriverService) GetDriver(ctx context.Context, driverID string) (*models.Driver, error) {
	return s.DriverStore.GetDriver(ctx, driverID)
}

func (s *DriverService) AddVehicle(ctx context.Context, driverID string, vehicle *models.Vehicle) (*models.Vehicle, error) {
	driver_id, err := primitive.ObjectIDFromHex(driverID)
	if err != nil {
		return nil, err
	}
	vehicle.DriverID = driver_id

	exists, err := s.VehicleStore.GetVehicleByRegNum(ctx, vehicle.RegNum)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	} else if exists != nil {
		return nil, errors.New("vehicle already added")
	}

	return s.VehicleStore.AddVehicle(ctx, vehicle)
}

func (s *DriverService) GetVehicleDetails(ctx context.Context, driverID string) (*models.Vehicle, error) {
	vehicle, err := s.VehicleStore.GetVehicleByDriverID(ctx, driverID)
	if err != nil {
		return nil, errors.New("no vehicle added")
	}
	return vehicle, nil
}

func (s *DriverService) SetStatusAvailable(ctx context.Context, driverID string) error {
	return s.DriverStore.SetStatusAvailable(ctx, driverID)
}

func (s *DriverService) SetStatusOffline(ctx context.Context, driverID string) error {
	return s.DriverStore.SetStatusOffline(ctx, driverID)
}
