package services

import (
	"api_example/app/database"
	. "api_example/app/types"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"time"
)

var db = database.GetConnection()

func AddTodo(c *fiber.Ctx) error {
	var todo Todo
	userSession := c.Locals("session").(map[string]interface{})
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	todo.Date = time.Now()
	uId, _ := uuid.Parse(userSession["Id"].(string))
	todo.UserId = uId
	if result := db.Create(&todo); result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": result.Error,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo Added Succesfully",
	})
}

func UpdateTodo(c *fiber.Ctx) error {
	queryValue := c.Params("id")
	//return c.SendStatus(fiber.StatusOK)
	var body Todo
	var todo Todo
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	db.Model(&todo).Update("name", body.Name)
	db.First(&todo, queryValue)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"updated todo": todo,
	})
}

func DeleteTodo(c *fiber.Ctx) error {
	var todo Todo
	queryValue := c.Params("id")
	//uid, _ := uuid.Parse(queryValue)
	db.Delete(&todo, queryValue)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Todo successfully deleted",
	})
}

func ListTodos(c *fiber.Ctx) error {
	var todos []Todo
	userSession := c.Locals("session").(map[string]interface{})

	uId := userSession["Id"].(string)
	db.Where("user_id = ?", uId).Find(&todos)

	if count := len(todos); count == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "There is no todo here",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"todos": todos,
	})
}

func GetTodo(c *fiber.Ctx) error {
	var todo Todo
	todoId := c.Params("id")
	if err := db.First(&todo, todoId).Error; err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No todo found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"todo": todo,
	})
}
