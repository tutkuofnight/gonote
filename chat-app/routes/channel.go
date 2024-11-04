package routes

import (
	database "chat-app/db"
	"chat-app/helper"
	"chat-app/middleware"
	"chat-app/types"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var db = database.GetDB()
var validate *validator.Validate

func ChannelRoutes(app fiber.Router) {
	r := app.Group("/channel").Use(middleware.RestrictUser)
	r.Post("/create", func(ctx *fiber.Ctx) error {
		validate = validator.New(validator.WithRequiredStructEnabled())
		user := ctx.Locals("user").(map[string]interface{})
		userId := int(user["id"].(float64))

		var channel types.Channel
		if err := ctx.BodyParser(&channel); err != nil {
			return fiber.ErrBadRequest
		}

		err := validate.Struct(&channel)
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorMessages := helper.RenderValidationErrors(validationErrors)
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": errorMessages,
				})
			}
		}

		channel.AdminId = userId
		result := db.Create(&channel)
		if result.Error != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": result.Error,
			})
		}
		return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Channel created for" + user["username"].(string),
		})
	})
	r.Get("/all", func (ctx *fiber.Ctx) error  {
		var channels []types.Channel
		db.Model(&types.Channel{}).Preload("Users").Find(&channels)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"channels": channels,
		})
	})
	r.Get("/:id", func(ctx *fiber.Ctx) error {
		var channel types.Channel
		result := db.Model(&types.Channel{}).Where("id = ?", ctx.Params("id")).Preload("Messages").Preload("Users").First(&channel)
		if result.Error != nil {
			ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  fiber.StatusBadRequest,
				"message": result.Error,
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  fiber.StatusOK,
			"channel": channel,
		})
	})
}
