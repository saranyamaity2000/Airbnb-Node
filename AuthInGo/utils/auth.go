package utils

import (
	env "AuthInGo/config/env"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var AuthSecretToken string

func init() {
	AuthSecretToken = env.GetString("JWT_SECRET", "TOKEN")
}

func HashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error hashing password:", err)
		return "", err
	}
	return string(hash), nil
}

func CheckPasswordHash(plainPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}

type UserClaims struct {
	UserId int64  `json:"id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewUserClaim(userId int64, email string) jwt.Claims {
	return &UserClaims{
		UserId:           userId,
		Email:            email,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
}

func ValidateAndExtractUserClaims(token string) (*UserClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (any, error) {
		// it was signed with jwt.SigningMethodHS256 , respective check needs to be done
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(AuthSecretToken), nil // Replace "TOKEN" with your actual secret key
	})
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("token claims are not of type UserClaims")
	}
	return claims, nil
}
