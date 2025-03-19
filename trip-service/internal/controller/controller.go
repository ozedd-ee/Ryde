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
	trip, err := s.TripService.NewRideRequest(c.Request.Context(), riderID, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"trip": trip})
}
