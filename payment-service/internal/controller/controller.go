package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"
	"ryde/utils"

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
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	var DriverAccountRequest models.DriverAccountRequest
	if err := c.ShouldBindJSON(&DriverAccountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	}
	DriverAccountIDs, err := pc.PaymentService.AddDriverAccounts(c.Request.Context(), claims.UserID, &DriverAccountRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, DriverAccountIDs)
}

func (pc *PaymentController) AddPaymentMethod(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	email := c.Param("email")
	authURL, err := pc.PaymentService.AddPaymentMethod(c.Request.Context(), claims.UserID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.Redirect(http.StatusFound, authURL)
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

func (pc *PaymentController) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	payment, err := pc.PaymentService.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, payment)
}
