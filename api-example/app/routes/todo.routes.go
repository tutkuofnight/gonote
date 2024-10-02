package routes

import (
	"api_example/app/middlewares"
	"api_example/app/services"
	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(app fiber.Router) {
	r := app.Group("/todo").Use(middlewares.GetLoggedUser)
	r.Get("/list", services.ListTodos)
	r.Get("/:id", services.GetTodo)
	r.Post("/add", services.AddTodo)
	r.Put("/:id/update", services.UpdateTodo)
	r.Delete("/:id/delete", services.DeleteTodo)
	r.Get("/test", services.Test)
}
