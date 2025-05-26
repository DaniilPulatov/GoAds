package middlware

import (
	"ads-service/internal/usecase/auth"
	"ads-service/internal/usecase/user"
	"ads-service/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Middleware struct {
	authService auth.AuthService
	userService user.UserAdvertisementService
}

func NewMiddleware(authService auth.AuthService, userService user.UserAdvertisementService) *Middleware {
	return &Middleware{
		authService: authService,
		userService: userService,
	}
}

func (m *Middleware) AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		data := strings.Split(authHeader, " ")
		if len(data) != 2 || data[0] != "Bearer" {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		rawtoken := data[1]

		// TODO: Secret key should be stored securely, not hardcoded
		token, err := jwt.ParseWithClaims(rawtoken, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			return utils.SecretKey, nil // Replace with your secret key
		})
		if err != nil {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*utils.CustomClaims)
		if  !ok || !token.Valid {
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		userID := claims.UserID

		isAdmin, err := m.authService.IsAdmin(c.Request.Context(), userID)
		if err != nil || !isAdmin {
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
	
}
