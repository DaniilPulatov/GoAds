package auth

import (
	"ads-service/internal/domain/entities"
	"ads-service/internal/usecase/auth"
	"log"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.UserAuthService
}

func NewAuthHandler(authService auth.UserAuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var user entities.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	if err := h.authService.Register(c.Request.Context(), user); err != nil {
		log.Println("problem in handler")
		c.JSON(500, gin.H{"error": "failed to register user: " + err.Error()})
		return
	}
	c.JSON(201, gin.H{"message": "user registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(400, gin.H{"error": "invalid request body: " + err.Error()})
		return
	}

	access, refresh, err := h.authService.Login(c.Request.Context(), loginReq.Phone, loginReq.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "failed to login: " + err.Error()})
		return
	}

	c.JSON(200, LoginResponse{access, refresh})
}
