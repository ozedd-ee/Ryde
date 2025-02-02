package data

import (
	"context"
	"ryde/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VehicleStore struct {
	VehicleStore *mongo.Collection
}

func NewVehicleStore(db *mongo.Database) *VehicleStore {
	return &VehicleStore{
		VehicleStore: db.Collection("vehicles"),
	}
}

func (s *VehicleStore) AddVehicle(ctx context.Context, vehicle *models.Vehicle) (*models.Vehicle, error) {
	result, err := s.VehicleStore.InsertOne(ctx, vehicle)
	if err != nil {
		return nil, err
	}
	vehicle.ID = result.InsertedID.(primitive.ObjectID)
	return vehicle, nil
}

// Get vehicle by registration number
func (s *VehicleStore) GetVehicleByRegNum(ctx context.Context, regNum string) (*models.Vehicle, error) {
	var vehicle *models.Vehicle

	filter := bson.M{"regnum": regNum}
	if err := s.VehicleStore.FindOne(ctx, filter).Decode(vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}

func (s *VehicleStore) GetVehicleByDriverID(ctx context.Context, id string) (*models.Vehicle, error) {
	var vehicle *models.Vehicle
	driverID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_driverID": driverID}
	if err := s.VehicleStore.FindOne(ctx, filter).Decode(vehicle); err != nil {
		return nil, err
	}
	return vehicle, nil
}
