package models

type FCMRequest struct {
	To   string `json:"to"`
	Data Order  `json:"data"`
}
