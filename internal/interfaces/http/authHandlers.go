package http

import (
	"log"
	"net/http"
	"salaryAdvance/internal/entity"
	"salaryAdvance/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	AuthService *services.AuthService
}


func (h *AuthHandlers) Register(c *gin.Context) {
	var authUser entity.AuthUser
	if err := c.ShouldBindJSON(&authUser); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	user.Username = authUser.Username
	user.PasswordHash = authUser.Password
	user.Role = entity.Uploader

	err := h.AuthService.Register(user)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}


func (h *AuthHandlers) RegisterAdmin(c *gin.Context) {
	var authUser entity.AuthUser
	if err := c.ShouldBindJSON(&authUser); err != nil {
		log.Printf("Error binding to JSON: %v", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	user.Username = authUser.Username
	user.PasswordHash = authUser.Password
	user.Role = entity.Admin

	err := h.AuthService.RegisterAdmin(user)
	if err != nil {
		log.Printf("Error registering admin: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Admin registered successfully"})
}



func (h *AuthHandlers) Login(c *gin.Context) {
	var authUser entity.AuthUser
	if err := c.ShouldBindJSON(&authUser); err != nil {
		log.Printf("Error binding to JSON: %v", err)
		c.JSON(400, gin.H{"Error": err.Error()})
		return
	}
	token, err := h.AuthService.Login(authUser.Username, authUser.Password, c.ClientIP()+":"+authUser.Username)

	if err != nil {
		log.Printf("Internal Service Error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Login successful")
	c.JSON(http.StatusOK, gin.H{"token": token})
}



func (h *AuthHandlers) Logout(c *gin.Context) {
	rawToken, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing auth token"})
		return
	}
	token, ok := rawToken.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth token"})
		return
	}

	if err := h.AuthService.Logout(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
