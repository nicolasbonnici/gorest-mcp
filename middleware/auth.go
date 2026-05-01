package middleware

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ValidateJWT(c *fiber.Ctx, log *slog.Logger) (string, string, error) {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("missing Authorization header")
	}

	// Extract Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", "", fmt.Errorf("invalid Authorization header format")
	}

	token := parts[1]
	if token == "" {
		return "", "", fmt.Errorf("missing token")
	}

	// Note: For v0.1.0, we rely on GoREST's JWT validation middleware
	// that should be configured before the MCP plugin endpoint.
	// The token validation is handled by GoREST core auth, and user context
	// is available in the Fiber context locals.

	// Extract user info from Fiber context (set by GoREST auth middleware)
	userID := c.Locals("user_id")
	if userID == nil {
		// If not available from GoREST middleware, we need to validate the token ourselves
		// For v0.1.0, we return an error indicating auth middleware is required
		return "", "", fmt.Errorf("user context not found - ensure GoREST auth middleware is enabled")
	}

	userEmail := c.Locals("user_email")
	if userEmail == nil {
		userEmail = ""
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", "", fmt.Errorf("invalid user_id type in context")
	}

	userEmailStr, _ := userEmail.(string)

	log.Debug("JWT validated successfully", "user_id", userIDStr, "email", userEmailStr)
	return userIDStr, userEmailStr, nil
}
