package routes

import (
	"github.com/chiboycalix/code-snippet-manager/handlers"
	"github.com/chiboycalix/code-snippet-manager/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SnippetRoute(app *fiber.App) {
	app.Get("/", middlewares.JWTMiddleware(), handlers.GetAllSnippets)
	app.Post("/snippets", middlewares.JWTMiddleware(), handlers.CreateSnippet)
	app.Post("/snippets/:id", handlers.DeleteSnippet)
	app.Put("/snippets/:id", handlers.UpdateSnippet)
}
