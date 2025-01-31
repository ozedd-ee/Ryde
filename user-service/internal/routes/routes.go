package routes

import (
	"ryde/internal/controller"

	"github.com/gin-gonic/gin"
)

// Define and group routes for user-related functions
func UserRoutes(router *gin.Engine, controller *controller.UserController) {
	//group user routes under /users path
	userGroup := router.Group("/api/v1/users")

	// users/:id
	userGroup.GET("/:id", controller.GetUser)
	// users/signin
	userGroup.POST("/signin", controller.Login)
	// users/register
	userGroup.POST("/register", controller.CreateUser)
}
