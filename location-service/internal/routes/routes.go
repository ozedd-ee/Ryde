package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

func LocationRoutes(router *gin.Engine, controller *controller.LocationController) {
	locationGroup := router.Group("/api/v1/location")
	websockets := router.Group("/api/v1/ws/location") // For websocket connections

	locationGroup.GET("/:id", controller.GetDriverLocation) // Get driver location by ID
	locationGroup.GET("/nearby-drivers", controller.FindNearbyDrivers) // Fetch nearby drivers for ride ordering

	websockets.GET("/update", controller.UpdateDriverLocation) // Transmit driver location updates from mobile app via websockets
}
