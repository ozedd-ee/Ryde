package main

import (
	"log"
	"net/http"
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

	// Initialize database
	db := db.Init()

	// Initialize all layers
	tripStore := data.NewTripStore(db)
	tripService := service.NewTripService(tripStore)
	tripController := controller.NewTripController(tripService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "trip service is live"})
	})

	// Define trip routes
	routes.TripRoutes(router, tripController)

	//start server
	if err := router.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err.Error())
	}
}
