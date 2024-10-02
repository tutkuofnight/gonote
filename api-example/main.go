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

	routes.UserRoutes(app)
	routes.TodoRoutes(app)

	app.Get("/stack", func(c *fiber.Ctx) error {
		return c.JSON(c.App().Stack())
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		//db := database.GetConnection()
		//var users []types.User
		//db.Find(&users)
		//return c.JSON(fiber.Map{
		//	"users": users,
		//})
		res := c.GetReqHeaders()
		return c.JSON(fiber.Map{
			"res": res,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
