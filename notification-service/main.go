package main

import (
	"log"
	"ryde/cache"
	"ryde/internal/controller"
	"ryde/internal/data"
	"ryde/internal/routes"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize router
	router := gin.Default()

	// Initialize Redis client
	cache := cache.Init()

	// Initialize all three layers
	tokenStore := data.NewTokenStore(cache)
	notificationService := service.NewNotificationService(tokenStore)
	notificationController := controller.NewNotificationController(notificationService)

	// Define and group routes
	routes.NotificationRoutes(router, notificationController)

	// Start server
	if err := router.Run(":8083"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
