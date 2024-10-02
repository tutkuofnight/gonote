package routes

import (
	"api_example/app/config"
	. "api_example/app/middlewares"
	"api_example/app/services"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	r := app.Group("/user")
	//r.Get("/")
	r.Post("/register", services.Register)
	r.Post("/login", services.Login)
	r.Post("/logout", AuthMiddleware(config.Secretkey), GetLoggedUser, services.Logout)
	r.Get("/get-session", AuthMiddleware(config.Secretkey), GetLoggedUser, services.GetSession)
	r.Post("/profile/image", AuthMiddleware(config.Secretkey), GetLoggedUser, services.ChangeProfileImage)
}
