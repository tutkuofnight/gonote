package routes

import (
	"api_example/app/config"
	"api_example/app/middlewares"
	"api_example/app/services"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	r := app.Group("/user")
	//r.Get("/")
	r.Post("/register", services.Register)
	r.Post("/login", services.Login)
	r.Get("/get-session", middlewares.AuthMiddleware(config.Secretkey), services.GetSession)
}
