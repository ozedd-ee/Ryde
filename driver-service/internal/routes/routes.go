package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

// Define and group routes for driver-related operations
func DriverRoutes(router *gin.Engine, controller *controller.DriverController) {
	driverGroup := router.Group("api/v1/drivers")

	//Driver signup:   /drivers/signup
	driverGroup.POST("/signup", controller.CreateDriver)

	// Driver login: /drivers/login
	driverGroup.POST("/login", controller.Login)

	// Add vehicle: /drivers/:id/addvehicle
	driverGroup.POST("/:id/addvehicle", controller.AddVehicle)

	// Set driver status to 'busy'
	driverGroup.POST("/set-status-busy", controller.SetStatusBusy)

	// Set driver status to 'available'
	driverGroup.POST("/set-status-available", controller.SetStatusAvailable)
	
	// Set driver status to 'offline'
	driverGroup.POST("/set-status-offline", controller.SetStatusOffline)

	// Get driver: /drivers/:id
	driverGroup.GET("/:id", controller.GetDriver)
	
	//Get driver's vrhicle:  /drivers/:id/vehiches
	driverGroup.GET("/:id/vehicles", controller.GetVehicle)  
}
