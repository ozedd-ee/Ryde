package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Trip struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	RiderID     primitive.ObjectID `json:"rider_id" bson:"rider_id"`
	DriverID    primitive.ObjectID `json:"driver_id" bson:"driver_id"`
	Origin      Location           `json:"origin" bson:"origin"`
	Destination Location           `json:"destination" bson:"destination"`
	Status      string             `json:"status" bson:"status"` // can be initiated, en route or completed
	StartTime   time.Time          `json:"start_time" bson:"start_time"`
	EndTime     time.Time          `json:"end_time" bson:"end_time"`
}

// For temporary caching of trips
type TripBuffer struct {
	Key string `json:"key"`
	Trip Trip `json:"trip"`
}
