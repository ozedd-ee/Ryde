package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"ryde/internal/data"
	"ryde/models"
)

type NotificationService struct {
	TokenStore *data.TokenStore
}

func NewNotificationService(tokenStore *data.TokenStore) *NotificationService {
	return &NotificationService{
		TokenStore: tokenStore,
	}
}

func (s *NotificationService) UpdateFCMToken(ctx context.Context, ownerID, token string) error {
	return s.TokenStore.UpdateFCMToken(ctx, ownerID, token)
}

func (s *NotificationService) NotifyDriver(ctx context.Context, driverID string, request models.Order) (string, error) {
	fcmServerKey := os.Getenv("TEST_FCM_SERVER_KEY")

	driverFcmToken, err := s.TokenStore.GetFCMToken(ctx, driverID)
	if err != nil {
		return "", err
	}
	notification := models.FCMRequest{
		To:   driverFcmToken,
		Data: request,
	}
	requestBody, err := json.Marshal(notification)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "key="+fcmServerKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response")
	}
	var resp models.Order
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", fmt.Errorf("error decoding JSON")
	}
	return resp.Status, nil
}
