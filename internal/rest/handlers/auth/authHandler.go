package auth

import (
	"ads-service/internal/domain/entities"
	"log"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with phone number and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body entities.User true "User registration payload"
// @Success 201 {object} map[string]string "user registered successfully"
// @Failure 400 {object} map[string]string "invalid request body"
// @Failure 500 {object} map[string]string "failed to register user"
// @Router /auth/register [post]
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

// Login godoc
// @Summary User login
// @Description Authenticate user and return access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body AuthRequest true "User login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string "invalid request body"
// @Failure 401 {object} map[string]string "failed to login"
// @Router /auth/login [post]
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
