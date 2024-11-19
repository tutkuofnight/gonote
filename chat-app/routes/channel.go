package routes

import (
	database "chat-app/db"
	"chat-app/helper"
	"chat-app/middleware"
	"chat-app/repository"
	"chat-app/types"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var db = database.GetDB()
var validate *validator.Validate

var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

func ChannelRoutes(app fiber.Router) {
	r := app.Group("/channel").Use(middleware.RestrictUser)
	r.Post("/create", func(ctx *fiber.Ctx) error {
		channelId := uuid.New().String()
		validate = validator.New(validator.WithRequiredStructEnabled())
		user := ctx.Locals("user").(map[string]interface{})
		userUsername := user["username"].(string)
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

		channel.Id = channelId
		channel.AuthorUsername = userUsername
		
		result := db.Create(&channel)
		if result.Error != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": result.Error,
			})
		}

		if  err := repository.AddUserChannels(userId, channelId); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err,
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Channel created for " + userUsername,
		})
	})
	r.Get("/all", func(ctx *fiber.Ctx) error {
		var channels []types.Channel
		db.Model(&types.Channel{}).Preload("Users").Find(&channels)
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"channels": channels,
		})
	})
	r.Get("/:id", func(ctx *fiber.Ctx) error {
		var channel types.Channel
		result := db.Model(&types.Channel{}).Where("id = ?", ctx.Params("id")).Preload("Users").First(&channel)
		if result.Error != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": result.Error,
			})
		}
		messages, err := repository.GetChannelMessages(ctx.Params("id"))
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": err,
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"channel":  channel,
			"messages": messages,
		})
	})

	r.Post("/invite/create", func(ctx *fiber.Ctx) error {
		var inviteBody types.Invite

		if err := ctx.BodyParser(&inviteBody); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Body cannot parsed",
			})
		}

		key := fmt.Sprintf("invite:%s", inviteBody.ChannelId)
		err := rdb.Set(ctx.Context(), key, inviteBody.ChannelId, inviteBody.Exp).Err()

		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Data cannot saved",
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Invite link created",
		})

	})

}
