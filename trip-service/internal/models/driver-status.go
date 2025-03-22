package models

type StatusUpdate struct {
	DriverID string `json:"id"`
	Status   string `json:"status"`
}
