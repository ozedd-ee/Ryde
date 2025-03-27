package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"
	"ryde/utils"

	"github.com/gin-gonic/gin"
)

type TripController struct {
	TripService *service.TripService
}

func NewTripController(tripService *service.TripService) *TripController {
	return &TripController{
		TripService: tripService,
	}
}

func (s *TripController) NewRideRequest(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	var order *models.Order
	if err = c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	riderID := claims.DriverID
	tripBuffer, err := s.TripService.NewRideRequest(c.Request.Context(), riderID, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"tripBuffer": tripBuffer})
}

// Pick-up rider form origin
func (s *TripController) StartTrip(c *gin.Context) {
	token := c.Query("token")
	_, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	// driverID := claims.DriverID

	tripKey := c.Param("trip-key")
	// TODO: Verify caller is driver

	tripBuffer, err := s.TripService.StartTrip(c.Request.Context(), tripKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, tripBuffer)
}

// Drop-off rider at destination
func (s *TripController) EndTrip(c *gin.Context) {
	token := c.Query("token")
	_, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	// driverID := claims.DriverID

	tripKey := c.Param("trip-key")
	// TODO: Verify caller is driver
	updatedTrip, err := s.TripService.EndTrip(c.Request.Context(), tripKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, updatedTrip)
}

func (s *TripController) GetTripByID(c *gin.Context) {
	id := c.Param("id")
	trip, err := s.TripService.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, trip)
}

func (s *TripController) GetAllDriverTrips(c *gin.Context) {
	driverID := c.Param("id")

	trips, err := s.TripService.GetAllDriverTrips(c.Request.Context(), driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, trips)
}

func (s *TripController) GetPendingTrip(c *gin.Context) {
	tripKey := c.Param("trip-key")

	trip, err := s.TripService.GetPendingTrip(c.Request.Context(), tripKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, trip)
}
