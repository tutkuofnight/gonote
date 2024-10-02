package services

import (
	"api_example/app/config"
	"api_example/app/helper"
	"api_example/app/repository"
	"api_example/app/types"
	"fmt"
	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v4"
	"path/filepath"
	"time"
)

func Register(c *fiber.Ctx) error {
	var user types.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	var findedUser types.User
	db.Model(&findedUser).Where("username = ?", user.Username).First(&findedUser)
	fmt.Println(findedUser)
	if findedUser.Username == user.Username {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "This user already exists",
		})
	}
	pass, _ := helper.HashPassword(user.Password)
	user.Password = pass
	db.Create(&user)
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
	var users []types.User
	db.Find(&users)
	findedUser, credentialsError = repository.FindByCredentials(users, user.Username, user.Password)
	if credentialsError != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Login failed",
		})
	}
	exp := time.Now().Add(24 * time.Hour)
	claims := jtoken.MapClaims{
		"data": map[string]interface{}{
			"Id":       findedUser.Id,
			"Username": findedUser.Username,
		},
		"exp": exp.Unix(),
	}

	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Secretkey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "token"
	cookie.Value = t
	cookie.Expires = exp
	c.Cookie(cookie)

	// Return the token
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": t,
	})
}

func ChangeProfileImage(c *fiber.Ctx) error {
	file, err := c.FormFile("profileImage")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot find file..")
	}
	dst := filepath.Join("static/profile", file.Filename)
	if err := c.SaveFile(file, dst); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("File is not saved :/")
	}
	user := c.Locals("session").(map[string]interface{})
	var findedUser types.User
	db.First(&findedUser, user["Id"].(string))
	db.Model(&findedUser).Update("ProfileImage", file.Filename)
	//db.Save(&findedUser)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": findedUser,
	})
}
func Logout(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Succesfully Logout",
	})
}

func GetSession(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"session": c.Locals("session"),
	})
}
