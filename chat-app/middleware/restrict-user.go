package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func RestrictUser(c *fiber.Ctx) error {
	secretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
		})
		return nil
	}

	var authToken string = c.Cookies("token")
	token, err := jwt.ParseWithClaims(authToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error":  "JWT token cannot parsed",
		})
	}
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error":  "JWT Token is not valid",
		})
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": fiber.StatusUnauthorized,
			"error":  "JWT Token Claims Error",
		})
	}
	c.Locals("user", (*claims)["user"])
	return c.Next()
}
