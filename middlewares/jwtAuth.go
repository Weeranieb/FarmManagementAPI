package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		paths := []string{
			"auth/register",
			"auth/login",
		}

		for _, path := range paths {
			if strings.Contains(c.Request.URL.Path, path) {
				c.Next()
				return
			}
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Check if the Authorization header is in the format "Bearer <token>"
		tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))

		// Parse and validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			secret := viper.GetString("authentication.jwt_secret")
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Set("userId", claims["userId"])
		c.Set("username", claims["username"])
		c.Set("userLevel", claims["userLevel"])
		c.Set("clientId", claims["clientId"])

		// Token is valid, continue processing the request
		c.Next()
	}
}
