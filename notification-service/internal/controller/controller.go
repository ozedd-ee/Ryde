package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	NotificationService *service.NotificationService
}

func NewNotificationController(notificationService *service.NotificationService) *NotificationController {
	return &NotificationController{
		NotificationService: notificationService,
	}
}

func (s *NotificationController) SendRideRequest(c *gin.Context) {
	var rideRequest models.RideRequest

	if err := c.ShouldBindJSON(&rideRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	response, err := s.NotificationService.SendRideRequest(rideRequest.DriverID, rideRequest.Order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.String(http.StatusOK, "UTF-8", response)
}
