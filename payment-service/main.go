package main

import (
	"log"
	"net/http"

	"ryde/db"
	"ryde/internal/controller"
	"ryde/internal/data"
	"ryde/internal/routes"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize gin
	router := gin.Default()

	// Initialize database
	db := db.Init()

	// Initialize all layers
	paymentStore := data.NewPaymentStore(db)
	paymentService := service.NewPaymentService(paymentStore)
	tripController := controller.NewPaymentController(paymentService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "payment service is live"})
	})

	// Define payment routes
	routes.PaymentRoutes(router, tripController)

	//start server
	if err := router.Run(":8084"); err != nil {
		log.Fatal("Failed to start server:", err.Error())
	}
}
