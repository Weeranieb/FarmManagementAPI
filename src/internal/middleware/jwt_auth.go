package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

// JWTAuthMiddleware validates JWT tokens and sets user context
func JWTAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Public paths that don't require authentication
		publicPaths := []string{
			"/api/v1/auth/register",
			"/api/v1/auth/login",
			"/api/v1/auth/logout",
			// "/api/v1/user", // System setup user endpoint
			"/swagger",
			"/health",
		}

		// Check if the current path is public
		for _, path := range publicPaths {
			if strings.Contains(c.Path(), path) {
				return c.Next()
			}
		}

		// Try to get token from cookie first, then fallback to Authorization header
		var tokenString string

		// Check cookie first
		cookieToken := c.Cookies("jwt_token")
		if cookieToken != "" {
			tokenString = cookieToken
		} else {
			// Fallback to Authorization header
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
			}

			// Check if the Authorization header is in the format "Bearer <token>"
			tokenString = strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", 1))
			if tokenString == "" {
				return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
			}
		}

		// Parse and validate the JWT token
		secretKey := viper.GetString("authentication.jwt_secret")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Check the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return http.Error(c, errors.ErrAuthTokenInvalid.Code, "Invalid token claims")
		}

		// Set user context in context.Context
		ctx := c.UserContext()
		if ctx == nil {
			ctx = context.Background()
		}
		if userId, ok := claims["userId"].(float64); ok {
			ctx = context.WithValue(ctx, constants.UserIDKey, int(userId))
		}
		if username, ok := claims["username"].(string); ok {
			ctx = context.WithValue(ctx, constants.UsernameKey, username)
		}
		if userLevel, ok := claims["userLevel"].(float64); ok {
			ctx = context.WithValue(ctx, constants.UserLevelKey, int(userLevel))
		}
		if clientId, ok := claims["clientId"].(float64); ok {
			ctx = context.WithValue(ctx, constants.ClientIDKey, int(clientId))
		}

		// Update the context in the fiber context
		c.SetUserContext(ctx)

		// Token is valid, continue processing the request
		return c.Next()
	}
}
