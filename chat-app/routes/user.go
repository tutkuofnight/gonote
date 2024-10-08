package routes

import (
	"github.com/gofiber/fiber/v2"
)

func UserRotues(app fiber.Router) {
	r := app.Group("/user")
	r.Post("/register" , func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
	})
	r.Post("/login" , func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
	})
	r.Get("/session" , func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{})
	})
}
