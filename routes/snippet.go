package routes

import (
	"github.com/chiboycalix/code-snippet-manager/handlers"

	"github.com/gofiber/fiber/v2"
)

func SnippetRoute(app *fiber.App) {
	app.Get("/", handlers.GetAllSnippets)
	app.Post("/snippets", handlers.CreateSnippet)
}
