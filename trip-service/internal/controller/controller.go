package controller

import (
	"errors"
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"
	"ryde/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
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
	trip, err := s.TripService.NewRideRequest(c.Request.Context(), riderID, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"trip": trip})
}

func (s *TripController) GetTripByID(c *gin.Context) {
	id := c.Param("id")
	trip, err := s.TripService.GetTripByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, trip)
}

// Pick-up rider form origin
func (s *TripController) StartTrip(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	driverID := claims.DriverID

	trip, err := s.TripService.GetTripByDriver(c.Request.Context(), driverID)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("trip not assigned to driver")})
	}

	updatedTrip, err := s.TripService.StartTrip(c.Request.Context(), trip.ID.String(), driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, updatedTrip)
}

// Drop-off rider at destination
func (s *TripController) EndTrip(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	driverID := claims.DriverID

	trip, err := s.TripService.GetTripByDriver(c.Request.Context(), driverID)
	if err != nil && err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	} else if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("trip not assigned to driver")})
	}

	updatedTrip, err := s.TripService.EndTrip(c.Request.Context(), trip.ID.String(), driverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, updatedTrip)
}
