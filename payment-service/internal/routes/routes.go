package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

func PaymentRoutes(router *gin.Engine, controller *controller.PaymentController) {
	paymentGroup := router.Group("/api/v1/payment")

	paymentGroup.POST("charge-card", controller.ChargeCard) // Charge rider card for ride payment

	paymentGroup.POST("/add-driver-account", controller.AddDriverAccount) // Add driver account

	paymentGroup.POST("/paystack-callback", controller.PaystackCallbackHandler) // Handle Paystack callbacks. To called by Paystack only.
}
