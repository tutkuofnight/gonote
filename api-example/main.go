package main

import (
	"api_example/app/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
)

func main() {
	app := fiber.New()

	app.Static("/static", "./static")
	app.Use(logger.New())

	//routes.TodoRoutes(app)
	routes.UserRoutes(app)
	app.Get("/stack", func(c *fiber.Ctx) error {
		return c.JSON(c.App().Stack())
	})
	log.Fatal(app.Listen(":3000"))
}
