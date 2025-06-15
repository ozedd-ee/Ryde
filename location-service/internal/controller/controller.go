package controller

import (
	"net/http"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

type LocationController struct {
	LocationService *service.LocationService
}

func NewLocationController(locationService *service.LocationService) *LocationController {
	return &LocationController{
		LocationService: locationService,
	}
}

func (s *LocationController) UpdateDriverLocation(c *gin.Context) {
	s.LocationService.UpdateDriverLocation(c)
}

func (s *LocationController) GetDriverLocation(c *gin.Context) {
	driver_id := c.Param("id")
	location, err := s.LocationService.GetDriverLocation(c, driver_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "driver location not found"})
		return
	}
	c.JSON(http.StatusOK, location)
}

func (s *LocationController) FindNearbyDrivers(c *gin.Context) {
	var request struct {
		latitude  float64
		longitude float64
		radius    float64
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"invalid request"})
		return
	}
	drivers, err := s.LocationService.FindNearbyDrivers(c, request.latitude, request.longitude, request.radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error fetching nearby drivers"})
		return
	}
	c.JSON(http.StatusOK, drivers)
}
