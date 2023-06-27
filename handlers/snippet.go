package handlers

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/models"
	"github.com/chiboycalix/code-snippet-manager/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var snippetCollection *mongo.Collection = configs.GetCollection(configs.DB, "snippets")

type ViewData struct {
	Token string
}

func GetAllSnippets(c *fiber.Ctx) error {
	_id, _ := utils.GetUserIDFromToken(c.Cookies("codeSnippetManagerJWT"))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var snippets []models.Snippet
	defer cancel()

	ownerId, _ := primitive.ObjectIDFromHex(_id)
	filter := bson.M{"owner": ownerId}
	results, err := snippetCollection.Find(ctx, filter)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't find any snippets",
		})
	}
	// reading from the db
	defer results.Close(ctx)
	for results.Next(ctx) {
		var singleSnippet models.Snippet
		if err := results.Decode(&singleSnippet); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Couldn't find any snippets",
			})
		}
		snippets = append(snippets, singleSnippet)
	}

	return c.Render("index", fiber.Map{
		"Snippets": snippets,
		"Theme":    "monokai",
	})
}

func CreateSnippet(c *fiber.Ctx) error {
	_id, _ := utils.GetUserIDFromToken(c.Cookies("codeSnippetManagerJWT"))
	snippet := new(models.Snippet)
	if err := c.BodyParser(snippet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Couldn't create snippet",
		})
	}

	if snippet.Description == "" || snippet.Snippet == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Description and snippet are required",
		})
	}
	snippet.Owner, _ = primitive.ObjectIDFromHex(_id)
	_, err := snippetCollection.InsertOne(context.Background(), snippet)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't create snippet",
		})
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't delete snippet",
		})
	}
	return c.Redirect("/")
}

func UpdateSnippet(c *fiber.Ctx) error {
	idParam := c.Params("id")
	match := regexp.MustCompile(`ObjectID\(%22(.*?)%22\)`).FindStringSubmatch(idParam)
	snippetId, _ := primitive.ObjectIDFromHex(match[1])
	snippet := new(models.Snippet)
	if err := c.BodyParser(snippet); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Couldn't update snippet",
		})
	}

	_, err := snippetCollection.UpdateOne(context.Background(), bson.M{"_id": snippetId}, bson.M{"$set": snippet})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Couldn't update snippet",
		})
	}

	return c.Redirect("/")
}
