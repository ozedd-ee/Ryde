package controller

import (
	"errors"
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"
	"ryde/utils"

	"github.com/gin-gonic/gin"
)

type DriverController struct {
	DriverService *service.DriverService
}

func NewDriverController(driverService *service.DriverService) *DriverController {
	return &DriverController{
		DriverService: driverService,
	}
}

func (s *DriverController) CreateDriver(c *gin.Context) {
	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	result, err := s.DriverService.SignUp(c.Request.Context(), &driver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, result)
}

func (s *DriverController) Login(c *gin.Context) {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}

	token, err := s.DriverService.Login(c.Request.Context(), request.Email, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to generate token"})
	}

	c.JSON(http.StatusOK, token)
}

func (s *DriverController) GetDriver(c *gin.Context) {
	id := c.Param("id")

	// call service layer to search for driver
	driver, err := s.DriverService.GetDriver(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "driver not found"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func (s *DriverController) AddVehicle(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("unauthorized")})
	}

	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
	newVehicle, err := s.DriverService.AddVehicle(c.Request.Context(), claims.DriverID, &vehicle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.JSON(http.StatusOK, newVehicle)
}

// Get vehicles by driver ID
func (s *DriverController) GetVehicle(c *gin.Context) {
	id := c.Param("id")

	// call service layer to search for vehicle
	vehicle, err := s.DriverService.GetVehicleDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "no vehicle found"})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

func (s *DriverController) SetStatusAvailable(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("unauthorized")})
	}
	err = s.DriverService.SetStatusAvailable(c.Request.Context(), claims.DriverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.Status(http.StatusOK)
}

func (s *DriverController) SetStatusOffline(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("unauthorized")})
	}
	err = s.DriverService.SetStatusOffline(c.Request.Context(), claims.DriverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.Status(http.StatusOK)
}
