package controller

import (
	"net/http"
	"ryde/internal/models"
	"ryde/internal/service"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		UserService: userService,
	}
}

func (s *UserController) CreateUser(c *gin.Context) {
	var user models.User

    // Bind input to struct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
        return
	}
    newUser, err := s.UserService.SignUp(c.Request.Context(), &user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
        return
    }
    // Respond with new user
    c.JSON(http.StatusOK, newUser)
}

func (s *UserController) Login(c *gin.Context) {
    var request struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error":"invalid request"})
        return
    }

    token, err := s.UserService.Login(c.Request.Context(), request.Email, request.Password)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unable to generate token"})
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (s *UserController) GetUser(c *gin.Context) {
	id := c.Param("id")

	// call service layer to search for user
	user, err := s.UserService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return
	}

    c.JSON(http.StatusOK, user)
}
