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

func (pc *PaymentController) NewSubAccount(c *gin.Context) {
	driverID := c.Param("driver-id")
	var subAccountRequest models.SubAccountRequest
	if err := c.ShouldBindJSON(&subAccountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	subAccountID, err := pc.PaymentService.CreateSubAccount(c.Request.Context(), driverID, &subAccountRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, subAccountID)
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
