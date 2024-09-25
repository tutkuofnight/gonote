package main

import (
	"api_example/app/types"
	"github.com/gofiber/fiber/v2"
	"log"
	"time"
)

func main() {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/user/save", func(c *fiber.Ctx) error {
		var body types.User
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		body.Todos[0].Date = time.Now()
		return c.Status(fiber.StatusOK).JSON(body)
	})
	log.Fatal(app.Listen(":3000"))
}
