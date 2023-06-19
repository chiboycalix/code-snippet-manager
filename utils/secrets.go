package utils

import (
	"errors"
	"net/http"
	"time"

	"github.com/chiboycalix/code-snippet-manager/configs"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type JWTClaim struct {
	ID string `json:"_id"`
	jwt.StandardClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hashed, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
}

func GenerateJWT(id string) (string, error) {
	claims := &JWTClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}
	tokenSecret := configs.EnvJWTSecret()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(tokenSecret))
}

func ValidateToken(signedToken string) (err error, claims *JWTClaim) {
	tokenSecret := configs.EnvJWTSecret()
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
	if !ok {
		fiber.NewError(http.StatusBadGateway, "Couldn't parse claims")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	return nil, claims
}

func GetUserIDFromToken(cookie string) (string, error) {
	if cookie == "" {
		return "", errors.New("unauthenticated")
	}
	err, claims := ValidateToken(cookie)
	if err != nil {
		return "", errors.New("unauthenticated")
	}
	return claims.ID, nil
}
