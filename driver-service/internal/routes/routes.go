package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

// Define and group routes for driver-related operations
func DriverRoutes(router *gin.Engine, controller *controller.DriverController) {
	driverGroup := router.Group("api/v1/drivers")

	//Driver sign-up:   /drivers/sign-up
	driverGroup.POST("/sign-up", controller.CreateDriver)

	// Driver login: /drivers/login
	driverGroup.POST("/login", controller.Login)

	// Add vehicle: /drivers/:id/add-vehicle
	driverGroup.POST("/:id/add-vehicle", controller.AddVehicle)

	// Set driver status to 'available'
	driverGroup.PUT("/set-status-available", controller.SetStatusAvailable)

	// Set driver status to 'offline'
	driverGroup.PUT("/set-status-offline", controller.SetStatusOffline)

	// Get driver: /drivers/:id
	driverGroup.GET("/:id", controller.GetDriver)
	
	//Get driver's vehicle:  /drivers/:id/vehicles
	driverGroup.GET("/:id/vehicles", controller.GetVehicle)  
}
