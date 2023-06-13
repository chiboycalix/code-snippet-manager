package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
		return fiber.NewError(http.StatusBadRequest, "Invalid request")
	}

	if user.Email == "" || user.Password == "" || user.Username == "" {
		return fiber.NewError(http.StatusBadRequest, "Email, Username and Password are required")
	}
	pass, hashErr := utils.HashPassword(user.Password)
	if hashErr != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to hash password")
	}
	user.Password = pass
	_, err := userCollection.InsertOne(context.Background(), user)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to save user")
	}

	return c.Redirect("/")
}

func LoginUser(c *fiber.Ctx) error {
	user := new(models.User)
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := c.BodyParser(user); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request")
	}
	fmt.Println(user.Email)
	fmt.Println(user.Password)
	if user.Email == "" || user.Password == "" {
		return fiber.NewError(http.StatusBadRequest, "Email and Password are required")
	}

	var result models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&result)
	if err != nil {
		return fiber.NewError(http.StatusNotFound, "Failed to logins")
	}

	if err := utils.CheckPasswordHash(result.Password, user.Password); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Failed to login")
	}

	jwt, err := utils.GenerateJWT(result.Email)
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Failed to generate JWT token")
	}
	fmt.Println(jwt)

	return c.Redirect("/")
}

func LogoutUser(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func LoginUserPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}

func RegisterUserPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{})
}
