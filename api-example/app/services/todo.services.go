package services

import (
	"api_example/app/types"
	"github.com/gofiber/fiber/v2"
)

func AddTodo(c *fiber.Ctx) error {
	var todo types.Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// Save db here..
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo Added Succesfully",
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest)
}

func DeleteTodo(c *fiber.Ctx) error {

}

func ListTodos(c *fiber.Ctx) error {

}

func GetTodo(c *fiber.Ctx) error {

}
