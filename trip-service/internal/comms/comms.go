package comms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ryde/internal/models"
	"ryde/utils"
)

func NotifyDriver(driverID string, request *models.Order) (*models.RideRequest, error) {
	// API call to notification service
	notification := models.RideRequest{
		DriverID: driverID,
		Order:    *request,
	}
	payload, err := json.Marshal(notification)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}
	response, err := http.Post("http://localhost:8083/api/v1/notifications/send-ride-request", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	notificationResponse := models.RideRequest{
		DriverID: driverID,
		Order:    *request,
	}
	err = json.Unmarshal(body, &notificationResponse)
	if err != nil {
		return nil, fmt.Errorf("invalid notification response format: %v", err)
	}
	return &notificationResponse, nil
}

func FindNearbyDrivers(request *models.Order) ([]string, error) {
	var payload struct {
		Latitude  float64
		Longitude float64
		Radius    float64
	}

	payload.Latitude = request.Origin.Latitude
	payload.Longitude = request.Origin.Longitude
	payload.Radius = 50 // 50km

	requestBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	// TODO:switch to gRPC after MVP
	response, err := http.Post("http://localhost:8082/api/v1/location/nearbyDrivers", "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	var drivers []string
	err = json.Unmarshal(body, &drivers)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	return drivers, nil
}

func SendChargeRequest(riderID string, request *models.ChargeRequest) (*models.Payment, error) {
	token, err := utils.GeneratePaymentJWT(riderID, request.Email)
	if err != nil {
		return nil, fmt.Errorf("error generating token: %v", err)
	}

	payload, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshalling request: %v", err)
	}

	client := http.DefaultClient
	req, err := http.NewRequest("POST", "http://localhost:8085/api/v1/payment/charge-card", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}
	payment := &models.Payment{}
	err = json.Unmarshal(body, payment)
	if err != nil {
		return nil, fmt.Errorf("unexpected response type: %v", err)
	}
	return payment, nil
}
