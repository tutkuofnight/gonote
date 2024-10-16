package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
)

func RestrictUser(c *fiber.Ctx) error {
	secretKey, ok := os.LookupEnv("JWT_SECRET_KEY")
	if !ok {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
		})
	}
	header := c.GetReqHeaders()
	bearerToken := strings.Split(header["Authorization"][0], " ")[1]
	token, err := jwt.ParseWithClaims(bearerToken, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
		})
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "Invalid token",
		})
	}
	c.Locals("user", (*claims)["user"])
	return c.Next()
}
