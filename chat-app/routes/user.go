package routes

import (
	"chat-app/helper"
	"chat-app/middleware"
	"chat-app/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func UserRotues(app fiber.Router) {
	r := app.Group("/user")
	r.Post("/register", func(ctx *fiber.Ctx) error {
		var user types.User
		if err := ctx.BodyParser(&user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validate := validator.New()
		if err := validate.Struct(&user); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := helper.RenderValidationErrors(validationErrors)
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": errorMessages,
				})
			}
		}

		var findedUser types.User
		db.Model(&types.User{}).Where("username = ?", user.Username).First(&findedUser)
		if findedUser.Username != "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username already exist",
			})
		}
		hashedPassword, _ := helper.HashPassword(user.Password)
		user.Password = hashedPassword
		err := db.Create(&user).Error
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "User not saved",
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User saved",
			"user":    user,
		})
	})
	r.Post("/login", func(ctx *fiber.Ctx) error {
		var user types.User
		if err := ctx.BodyParser(&user); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "data cannot read",
			})
		}
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := helper.RenderValidationErrors(validationErrors)
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": errorMessages,
				})
			}
		}
		var findedUser types.User
		err := db.Model(&types.User{}).Where("username = ?", user.Username).First(&findedUser).Error
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "User not found",
			})
		}
		passwordMatch := helper.MatchPassword(findedUser.Password, user.Password)
		if !passwordMatch {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "password does not match",
			})
		}

		claims := jwt.MapClaims{
			"user": map[string]interface{}{
				"id":       findedUser.Id,
				"username": findedUser.Username,
			},
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "Token signing failed",
			})
		}
		return ctx.JSON(fiber.Map{"token": t})
	})
	r.Post("/logout", middleware.AuthMiddleware(), func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Succesfully Logout",
		})
	})
	r.Get("/session", middleware.RestrictUser, func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": ctx.Locals("user"),
		})
	})

}
