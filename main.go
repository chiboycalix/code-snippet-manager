package main

import (
	"fmt"
	"log"
	"os"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})

	app.Use(cors.New())
	app.Static("/", "./public", fiber.Static{})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// database connection
	configs.ConnectDatabase()

	// routes
	routes.SnippetRoute(app)
	routes.UserRoute(app)

	port := os.Getenv("PORT")

	fmt.Println("Server started on http://localhost:" + port)
	if port == "" {
		port = "4000"
	}
	log.Fatal(app.Listen(":" + port))
}
