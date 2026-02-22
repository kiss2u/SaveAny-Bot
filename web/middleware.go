package web

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

// AuthConfig holds authentication configuration
type AuthConfig struct {
	Username string
	Password string
}

// NewAuthMiddleware creates basic auth middleware if credentials are configured
func NewAuthMiddleware(cfg *AuthConfig) fiber.Handler {
	if cfg.Username == "" && cfg.Password == "" {
		// No authentication configured
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	// Create basic auth middleware
	auth := basicauth.New(basicauth.Config{
		Users: map[string]string{
			cfg.Username: cfg.Password,
		},
		Realm: "SaveAny-Bot",
		Authorizer: func(user, pass string) bool {
			return user == cfg.Username && pass == cfg.Password
		},
		Unauthorized: func(c *fiber.Ctx) error {
			c.Set("WWW-Authenticate", `Basic realm="SaveAny-Bot"`)
			return c.Status(401).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		},
	})

	return auth
}

// APIKeyAuth middleware for API endpoints
func APIKeyAuth(apiKey string) fiber.Handler {
	if apiKey == "" {
		return func(c *fiber.Ctx) error {
			return c.Next()
		}
	}

	return func(c *fiber.Ctx) error {
		// Check API key in header
		key := c.Get("X-API-Key")
		if key == "" {
			// Check API key in query param
			key = c.Query("api_key")
		}

		if key == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "API key required",
			})
		}

		if key != apiKey {
			return c.Status(403).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		return c.Next()
	}
}

// CORSConfig for development
func NewCORSConfig() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}

		return c.Next()
	}
}

// ParseBasicAuth parses Basic Auth header
func ParseBasicAuth(authHeader string) (username, password string, ok bool) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[6:])
	if err != nil {
		return
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return
	}

	username = parts[0]
	password = parts[1]
	ok = true
	return
}
