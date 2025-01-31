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
	userStore := data.NewUserStore(db)
	userService := service.NewUserService(userStore)
	UserController := controller.NewUserController(userService)

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "user service is live"})
	})

	// Define user routes
	routes.UserRoutes(router, UserController)

	//start server
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err.Error())
	}
}
