package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

// Define and group routes for notification-related functions
func NotificationRoutes(router *gin.Engine, controller *controller.NotificationController) {
	//group notification routes under /notifications path
	notificationGroup := router.Group("/api/v1/notifications")

	// notifications/send-ride-request
	notificationGroup.POST("/send-ride-request", controller.SendRideRequest)
}
