package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
)

func AuthMiddleware(secret string) fiber.Handler {
	return jwt.New(jwt.Config{
		SigningKey: []byte(secret),
	})
}
