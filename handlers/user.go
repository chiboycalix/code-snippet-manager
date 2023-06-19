package handlers

import (
	"context"
	"fmt"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/chiboycalix/code-snippet-manager/models"
	"github.com/chiboycalix/code-snippet-manager/utils"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func RegisterUser(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email, Password and Username are required",
		})
	}
	pass, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}
	user.Password = pass
	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	errr := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&user)
	if errr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}
	jwt, err := utils.GenerateJWT(user.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate jwt",
		})
	}

	utils.SetCookies(jwt, c)
	return c.Redirect("/")
}

func LoginUser(c *fiber.Ctx) error {
	type loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var user loginRequest

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and Password are required",
		})
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to login",
		})
	}

	if err := utils.CheckPasswordHash(result.Password, user.Password); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	jwt, err := utils.GenerateJWT(result.ID.Hex())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate jwt",
		})
	}
	utils.SetCookies(jwt, c)
	return c.Redirect("/")
}

func LogoutUser(c *fiber.Ctx) error {
	fmt.Println("Logout")
	c.ClearCookie("codeSnippetManagerJWT")
	return c.Render("login", fiber.Map{})
}

func LoginUserPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func RegisterUserPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{})
}

func SubmitLoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"message": "Login successful",
	})
}
