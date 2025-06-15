package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"
	"ryde/utils"
	"strings"

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
		return
	}
	var DriverAccountRequest models.SubAccountRequest
	if err := c.ShouldBindJSON(&DriverAccountRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	DriverAccountIDs, err := pc.PaymentService.AddDriverAccount(c.Request.Context(), claims.UserID, &DriverAccountRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, DriverAccountIDs)
}

func (pc *PaymentController) AddPaymentMethod(c *gin.Context) {
	token := c.Query("token")
	claims, err := utils.ValidateJWT(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	email := c.Param("email")
	authURL, err := pc.PaymentService.AddPaymentMethod(c.Request.Context(), claims.UserID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.Redirect(http.StatusFound, authURL)
}

func (pc *PaymentController) PaystackCallbackHandler(c *gin.Context) {
	reference := c.Query("reference")
	if reference == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing transaction reference"})
		return
	}

	err := pc.PaymentService.PaystackCallbackHandler(c.Request.Context(), reference)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.Status(http.StatusOK)
}

func (pc *PaymentController) ChargeCard(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}
	token := parts[1]
	_, err := utils.ValidateChargeRequest(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized request"})
		return
	}

	var chargeRequest models.ChargeRequest
	if err := c.ShouldBindJSON(&chargeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	payment, err := pc.PaymentService.ChargeCard(c.Request.Context(), &chargeRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to charge card"})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (pc *PaymentController) GetPayment(c *gin.Context) {
	paymentID := c.Param("id")
	payment, err := pc.PaymentService.GetPayment(c.Request.Context(), paymentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, payment)
}
