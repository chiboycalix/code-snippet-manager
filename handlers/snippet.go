package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var snippetCollection *mongo.Collection = configs.GetCollection(configs.DB, "snippets")

func GetAllSnippets(c *fiber.Ctx) error {
	cursor, err := snippetCollection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(context.Background())

	var snippets []models.Snippet

	for cursor.Next(context.Background()) {
		var snippet models.Snippet
		err := cursor.Decode(&snippet)
		if err != nil {
			log.Fatal(err)
		}
		snippets = append(snippets, snippet)
	}
	return c.Render("index", fiber.Map{
		"Snippets": snippets,
		"Theme":    "monokai",
		// "Theme": "sunburst",
	})
}

func CreateSnippet(c *fiber.Ctx) error {
	snippet := new(models.Snippet)
	if err := c.BodyParser(snippet); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request")
	}

	if snippet.Description == "" || snippet.Snippet == "" {
		return fiber.NewError(http.StatusBadRequest, "Name and code are required")
	}

	_, err := snippetCollection.InsertOne(context.Background(), snippet)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to save snippet")
	}

	return c.Redirect("/")
}

// delete snippet
func DeleteSnippet(c *fiber.Ctx) error {
	// Get the todo ID from the URL parameter
	idParam := c.Params("id")

	snippetId, _ := primitive.ObjectIDFromHex(idParam)
	res, err := snippetCollection.DeleteOne(context.Background(), bson.M{"id": snippetId})
	fmt.Println(err, "err")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Failed to delete snippetss")
	}
	fmt.Println(res, "res")
	return c.Redirect("/")
}
