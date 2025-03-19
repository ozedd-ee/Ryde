package models

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Order struct {
	Origin      Location `json:"origin"`
	Destination Location `json:"destination"`
	Status      string   `json:"status"` // can be pending, accepted or rejected
}

type RideRequest struct {
	DriverID string `json:"driver_id"`
	Order    Order  `json:"order"`
}
