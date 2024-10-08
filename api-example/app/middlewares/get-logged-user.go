package middlewares

import "C"
import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"strings"
)

// TODO: 12. Satirda arrayin 1. elemani yoksa ne yapilacak? Burasi kontrol edilmeli

func GetLoggedUser(c *fiber.Ctx) error {
	header := c.GetReqHeaders()
	bearerToken := strings.Split(header["Authorization"][0], " ")[1]
	if len(bearerToken) > 0 {
		token, _, err := new(jwt.Parser).ParseUnverified(bearerToken, jwt.MapClaims{})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Locals("session", claims["data"])
	}
	return c.Next()
}
