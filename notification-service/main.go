package main

import (
	"log"
	"ryde/internal/controller"
	"ryde/internal/routes"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	mqttBroker := "tcp://localhost:1883"
	// Initialize router
	router := gin.Default()

	// Initialize both layers
	notificationService := service.NewNotificationService(mqttBroker)
	notificationController := controller.NewNotificationController(notificationService)

	// Define and group routes
	routes.NotificationRoutes(router, notificationController)

	// Start server
	if err := router.Run(":8083"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
