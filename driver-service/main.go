package main

import (
	"log"
	"ryde/cache"
	"ryde/db"
	"ryde/internal/controller"
	"ryde/internal/data"
	"ryde/internal/routes"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize gin
	router := gin.Default()
	
	// Initialize database and cache
	db := db.Init()
	cache := cache.Init()

	// Initialize all layers
	driverStore := data.NewDriverStore(db, cache)
	vehicleStore := data.NewVehicleStore(db)
	driverService := service.NewDriverService(driverStore, vehicleStore)
	driverController := controller.NewDriverController(driverService)

	// define driver routes
	routes.DriverRoutes(router, driverController)

	// start server
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
