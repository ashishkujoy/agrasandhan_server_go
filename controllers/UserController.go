package controllers

import (
	"ashishkujoy/agrasandhan/services"
	"github.com/gin-gonic/gin"
)

type UserCreationRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  int    `json:"role"`
}

func AddUser(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody UserCreationRequest
		err := c.ShouldBindJSON(&reqBody)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		user, err := userService.CreateUser(reqBody.Name, reqBody.Email, reqBody.Role)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(201, user)
	}
}

func GetAllUsers(service *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := service.GetAllUsers()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, users)
	}
}
