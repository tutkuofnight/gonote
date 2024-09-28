package services

import (
	"api_example/app/config"
	"api_example/app/helper"
	. "api_example/app/repository"
	"api_example/app/types"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"time"
)

var users []types.User

func Register(c *fiber.Ctx) error {
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	pass, _ := helper.HashPassword(user.Password)
	user.Password = pass
	users = append(users, user)
	// Save db here
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": user,
	})
}

func Login(c *fiber.Ctx) error {
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var credentialsError error
	var findedUser *types.User
	findedUser, credentialsError = FindByCredentials(users, user.Username, user.Password)
	if credentialsError != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Login failed",
		})
	}
	day := time.Hour * 24
	claims := jtoken.MapClaims{
		"Id":       findedUser.Id,
		"Username": findedUser.Username,
		"todos":    findedUser.Todos,
		"exp":      time.Now().Add(day * 1).Unix(),
	}
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Secretkey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Return the token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": t,
	})
}

func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Succesfully Logout",
	})
}

func GetSession(c *fiber.Ctx) error {
	user := c.Locals("user").(*jtoken.Token)
	claims := user.Claims.(jtoken.Claims)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"session": claims,
	})
}
