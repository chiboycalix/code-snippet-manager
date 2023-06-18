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
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var snippetCollection *mongo.Collection = configs.GetCollection(configs.DB, "snippets")

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func GetAllSnippets(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")

	// Check if the header value is present and in the expected format
	if authorization == "" || len(authorization) < 7 || authorization[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or missing authorization header",
		})
	}

	// Extract the token value
	token := authorization[7:]
	err := utils.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired JWT",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var snippets []models.Snippet
	defer cancel()

	results, err := snippetCollection.Find(ctx, bson.M{})
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
		// "Theme": "sunburst",
	})
}

func CreateSnippet(c *fiber.Ctx) error {
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
