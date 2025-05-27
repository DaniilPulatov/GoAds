package auth

import (
	"ads-service/internal/domain/entities"
	"log"

	"github.com/gin-gonic/gin"
)

func (h *AuthHandler) Register(c *gin.Context) {
	var user entities.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.userAuthService.Register(c.Request.Context(), &user); err != nil {
		log.Println("problem in handler")
		c.JSON(500, gin.H{"error": "failed to register user: " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq AuthRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	access, refresh, err := h.userAuthService.Login(c.Request.Context(), loginReq.Phone, loginReq.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "failed to login: " + err.Error()})
		return
	}

	c.JSON(200, LoginResponse{access, refresh})
}
