package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

type PaymentController struct {
	PaymentService *service.PaymentService
}

func NewPaymentController(paymentService *service.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

func (pc *PaymentController) AddDriverAccount(c *gin.Context) {
	driverID := c.Param("driver-id")
	var DriverAccountRequest models.DriverAccountRequest
	if err := c.ShouldBindJSON(&DriverAccountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	DriverAccountIDs, err := pc.PaymentService.AddDriverAccounts(c.Request.Context(), driverID, &DriverAccountRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, DriverAccountIDs)
}

func (pc *PaymentController) PaystackCallbackHandler(c *gin.Context) {
	reference := c.Query("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing transaction reference"})
	}

	err := pc.PaymentService.PaystackCallbackHandler(c.Request.Context(), reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}

	c.Status(http.StatusOK)
}

// TODO: Add authentication
func (pc *PaymentController) ChargeCard(c *gin.Context) {
	var chargeRequest models.ChargeRequest
	if err := c.ShouldBindJSON(&chargeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	}

	payment, err := pc.PaymentService.ChargeCard(c.Request.Context(), &chargeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to charge card"})
	}
	c.JSON(http.StatusOK, payment)
}
