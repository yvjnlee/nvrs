package middleware

import (
	"nvrs-gateway/config"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTMiddleware returns middleware for securing routes
func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  config.JWTSecret,
		TokenLookup: "header:Authorization", // Looks for tokens in the Authorization header
		AuthScheme:  "Bearer",
	})
}

// GenerateToken generates a JWT token for an agent
func GenerateToken(agentID int, agentName string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"agent_id":   agentID, // Include agent ID in the JWT claims
		"agent_name": agentName,
	})
	return token.SignedString(config.JWTSecret)
}
