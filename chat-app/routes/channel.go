package routes

import "github.com/gofiber/fiber/v2"

func ChannelRoutes(app fiber.Router) {
	r := app.Group("/channel")
	r.Post("/create", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
	})
}
