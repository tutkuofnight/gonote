package routes

import (
	"chat-app/helper"
	"chat-app/middleware"
	"chat-app/types"
	"os"
	"path/filepath"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func UserRotues(app fiber.Router) {
	r := app.Group("/user")

	r.Post("/register", func(ctx *fiber.Ctx) error {
		var body types.RegisterDto
		var user types.User
		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		validate := validator.New()
		if err := validate.Struct(&body); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := helper.RenderValidationErrors(validationErrors)
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": errorMessages,
				})
			}
		}

		var findedUser types.User
		db.Model(&types.User{}).Where("username = ?", body.Username).First(&findedUser)
		if findedUser.Username != "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "username already exist",
			})
		}
		hashedPassword, _ := helper.HashPassword(body.Password)
		
		user.Username = body.Username
		user.Password = hashedPassword

		err := db.Create(&user).Error
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "User not saved",
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "User successfully registered!",
		})
	})

	r.Post("/login", func(ctx *fiber.Ctx) error {
		var body types.LoginDto
		if err := ctx.BodyParser(&body); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   err.Error(),
				"message": "data cannot read",
			})
		}
		validate := validator.New()
		if err := validate.Struct(body); err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := helper.RenderValidationErrors(validationErrors)
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  fiber.StatusBadRequest,
					"message": errorMessages,
				})
			}
		}
		var findedUser types.User
		err := db.Model(&types.User{}).Where("username = ?", body.Username).Preload("Channels").First(&findedUser).Error
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  err.Error(),
				"message": "User not found",
			})
		}
		passwordMatch := helper.MatchPassword(findedUser.Password, body.Password)
		if !passwordMatch {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "password does not match",
			})
		}
		tokenExp := time.Now().Add(time.Hour * 24).Unix()
		claims := jwt.MapClaims{
			"user": map[string]interface{}{
				"id":       findedUser.Id,
				"username": findedUser.Username,
			},
			"exp": tokenExp,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Token signing failed",
			})
		}
		return ctx.JSON(fiber.Map{
			"status": fiber.StatusOK,
			"token":  t,
			"exp":    tokenExp,
			"user": fiber.Map{
				"username":     findedUser.Username,
				"id":           findedUser.Id,
				"profileImage": findedUser.ProfileImage,
			},
		})
	})

	r.Post("/avatar", middleware.AuthMiddleware(), func(ctx *fiber.Ctx) error {
		file, err := ctx.FormFile("avatar")
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{})
		}
		dst := filepath.Join("static/avatars", file.Filename)
		if err := ctx.SaveFile(file, dst); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": "Avatar not saved",
			})
		}
		user := ctx.Locals("session").(map[string]interface{})
		var findedUser types.UserDto
		db.First(&findedUser, user["Id"].(string))
		db.Model(&findedUser).Update("avatar", file.Filename)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"message": "Avatar successfully saved!",
		})

	})

	r.Post("/logout", middleware.AuthMiddleware(), func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Succesfully Logout",
		})
	})

	r.Get("/session", middleware.RestrictUser, func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(map[string]interface{})
		var findedUser types.User
		// db.First(&findedUser, int(user["id"].(float64))).Preload("Channels")
		db.Model(&types.User{}).Preload("Channels").Find(&findedUser, int(user["id"].(float64)))
		if len(findedUser.Password) > 0 {
			findedUser.Password = ""
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"user": findedUser,
		})
	})
}
