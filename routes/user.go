package routes

import (
	"github.com/chiboycalix/code-snippet-manager/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {
	app.Post("/auth/register", handlers.RegisterUser)
	app.Post("/auth/login", handlers.LoginUser)
	app.Get("/auth/logout", handlers.LogoutUser)
	app.Get("/login", handlers.LoginUserPage)
	app.Get("/register", handlers.RegisterUserPage)
	// app.Get("/submit-login", handlers.SubmitLoginPage)
}
