package main

import (
	"fmt"
	"log"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
)

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Use(cors.New())
	app.Static("/", "./public", fiber.Static{})
	// database connection
	configs.ConnectDatabase()
	routes.UserRoute(app)
	routes.SnippetRoute(app)

	// server port
	port := configs.EnvPort()

	fmt.Println("Server started on http://localhost:" + port)
	if port == "" {
		port = "4000"
	}
	log.Fatal(app.Listen(":" + port))
}
