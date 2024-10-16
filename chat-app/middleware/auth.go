package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwt "github.com/gofiber/jwt/v3"
	"os"
)

func AuthMiddleware() fiber.Handler {
	return jwt.New(jwt.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET_KEY")),
	})
}
