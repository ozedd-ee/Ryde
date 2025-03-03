package controller

import (
	"net/http"
	"ryde/internal/service"
	"ryde/utils"

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

func (s *NotificationController) UpdateFCMToken(c *gin.Context) {
	jwtToken := c.Query("token")
	claims, err := utils.ValidateJWT(jwtToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}
	ownerID := claims.OwnerID

	// Firebase Cloud Messaging token
	var FCMToken string
	if err := c.ShouldBindJSON(&FCMToken); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	err = s.NotificationService.UpdateFCMToken(c.Request.Context(), ownerID, FCMToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.Writer.WriteHeader(http.StatusOK)
}
