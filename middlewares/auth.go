package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the token from the cookies or headers
		token := c.Cookies("codeSnippetManagerJWT")
		if token == "" {
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Redirect("/login")
			}
			token = authHeader[7:] // Remove "Bearer " from the beginning
		}
		// Attach the token to the request headers
		c.Request().Header.Set("Authorization", "Bearer "+token)

		// Continue processing the request
		return c.Next()
	}
}
