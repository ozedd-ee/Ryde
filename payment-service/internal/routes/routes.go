package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine, controller *controller.PaymentController) {
	paymentGroup := router.Group("/api/v1/payment")

	paymentGroup.GET("/get-payment", controller.GetPayment) // Get payment

	paymentGroup.POST("/charge-card", controller.ChargeCard) // Charge rider card for ride payment

	paymentGroup.POST("/add-driver-account", controller.AddDriverAccount) // Add driver account

	paymentGroup.POST("/paystack-callback", controller.PaystackCallbackHandler) // Handle Paystack callbacks. To be called by Paystack only.

	paymentGroup.POST("/add-payment-method", controller.AddPaymentMethod) // Add rider payment method
}
