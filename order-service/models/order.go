package models

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Order struct {
	Origin      Location `json:"origin"`
	Destination Location `json:"destination"`
}
