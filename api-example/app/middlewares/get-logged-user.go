package middlewares

import "C"
import (
	"api_example/app/config"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
)

func GetLoggedUser(c *fiber.Ctx) error {
	fmt.Println(c.Cookies("token"))
	token, err := jtoken.Parse(c.Cookies("token"), func(token *jtoken.Token) (interface{}, error) {
		return []byte(config.Secretkey), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims := token.Claims.(jtoken.MapClaims)
	data := claims["data"].(map[string]interface{})
	c.Locals("session", data)
	return c.Next()
}
