package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/gofiber/fiber"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

type JWTClaim struct {
	Email string `json:"email,omitempty" validate:"required"`
	jwt.StandardClaims
}

func GenerateJWT(email string) (string, error) {
	claims := &JWTClaim{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
		},
	}
	tokenSecret := configs.EnvJWTSecret()
	fmt.Println(tokenSecret, "tokenSecret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateToken(signedToken string) (err error) {
	tokenSecret := configs.EnvJWTSecret()
	fmt.Println(tokenSecret, "tokenSecret")
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*JWTClaim)
	fmt.Println(claims, "claims")
	if !ok {
		fiber.NewError(http.StatusBadGateway, "Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}

	return

}
