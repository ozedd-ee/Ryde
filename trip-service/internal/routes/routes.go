package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

func TripRoutes(router *gin.Engine, controller *controller.TripController) {
	tripGroup := router.Group("/api/v1/trip")

	tripGroup.POST("/new-ride", controller.NewRideRequest) // Request a new ride

	tripGroup.POST("/start-trip", controller.StartTrip) // Pick-up rider origin

	tripGroup.POST("/end-trip", controller.EndTrip) // Drop-off rider at destination

	tripGroup.GET("/:id", controller.GetTripByID) // Retrieve trip details by ID
}
