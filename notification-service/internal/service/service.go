package service

import (
	"encoding/json"
	"fmt"
	"log"
	"ryde/internal/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	clientID = "ryde-notification-service"
	qos      = 1
)

type NotificationService struct {
	mqttClient mqtt.Client
}

func NewNotificationService(mqttBroker string) *NotificationService {
	options := mqtt.NewClientOptions()
	options.AddBroker(mqttBroker)
	options.SetClientID(clientID)

	mqttClient := mqtt.NewClient(options)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("MQTT connection error: %v", token.Error())
	}
	return &NotificationService{
		mqttClient: mqttClient,
	}
}

// Publish trip request to driver
func (s *NotificationService) SendRideRequest(driverID string, order models.Order) (string, error) {
	topic := fmt.Sprintf("ride/request/%s", driverID)
	payload, err := json.Marshal(order)
	if err != nil {
		return "", err
	}
	token := s.mqttClient.Publish(topic, qos, false, payload)
	token.Wait()
	// Wait for and return driver's response
	return s.subscribeToDriverResponse(driverID), nil

}

func (s *NotificationService) subscribeToDriverResponse(driverID string) string {
	topic := fmt.Sprintf("ride/response/%s", driverID)
	var response string
	// Subscribe to driver's response topic and wait 10 seconds for a confirmation or decline
	s.mqttClient.Subscribe(topic, qos, func(client mqtt.Client, msg mqtt.Message) {
		for {
			select {
			// stop waiting for response after 10 seconds
			case <-time.After(10 * time.Second):
				return
			default:
				var responsePayload map[string]string
				json.Unmarshal(msg.Payload(), &responsePayload)

				if responsePayload["Status"] == "accepted" {
					response = "accepted"
				}
			}
		}

	})
	return response
}
