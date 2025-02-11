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
	locationStore := data.NewLocationStore(cache)
	locationService := service.NewLocationService(locationStore)
	locationController := controller.NewLocationController(locationService)

	// Define and group routes
	routes.LocationRoutes(router, locationController)

	// Start server
	if err := router.Run(":8082"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
