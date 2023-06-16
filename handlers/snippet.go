package handlers

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var snippetCollection *mongo.Collection = configs.GetCollection(configs.DB, "snippets")

func GetAllSnippets(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var snippets []models.Snippet
	defer cancel()

	results, err := snippetCollection.Find(ctx, bson.M{})
	if err != nil {
		fiber.NewError(http.StatusBadGateway, "StatusBadGateway")
	}

	// reading from the db
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSnippet models.Snippet
		if err := results.Decode(&singleSnippet); err != nil {
			fiber.NewError(http.StatusBadGateway, "StatusBadGateway")
		}
		snippets = append(snippets, singleSnippet)
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
	idParam := c.Params("id")
	match := regexp.MustCompile(`ObjectID\(%22(.*?)%22\)`).FindStringSubmatch(idParam)
	snippetId, _ := primitive.ObjectIDFromHex(match[1])
	_, err := snippetCollection.DeleteOne(context.Background(), bson.M{"_id": snippetId})
	fmt.Println(err, "err")
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Failed to delete snippetss")
	}
	return c.Redirect("/")
}
