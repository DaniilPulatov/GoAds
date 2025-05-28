package middleware

import (
	"ads-service/pkg/utils"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (m *Middleware) UserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("YOU ARE HERE~!!!!")
		authHeader := c.GetHeader("Authorization")
		data := strings.Split(authHeader, " ")
		if len(data) != 2 || data[0] != "Bearer" {
			log.Println("Invalid Authorization header format")
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		rawtoken := data[1]

		token, err := jwt.ParseWithClaims(rawtoken, &utils.CustomClaims{},
			func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET_KEY")), nil
			})
		if err != nil {
			log.Println("#Err parsing token:", err)
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok || !token.Valid {
			log.Println("#Err validating token:", err)
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
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
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil // TODO: remember
		})
		if err != nil {
			log.Println("Error parsing token:", err)
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(*utils.CustomClaims)
		if !ok || !token.Valid {
			log.Println("Error validating token:")
			c.JSON(401, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}
		userID := claims.UserID

		isAdmin, err := m.authService.IsAdmin(c.Request.Context(), userID)
		if err != nil || !isAdmin {
			log.Println("User is not admin or error checking admin status:", err)
			c.JSON(403, gin.H{"error": "forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}

}
