package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"time"
)

// Request is for auth log in
type Request struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func Main(in Request) (*Response, error) {
	user := getUserByEmail(in.UsernameOrEmail)
	if user == nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("user not found")), fmt.Errorf("user not found")), fmt.Errorf("user not found")
	}
	if !validatePassword(in.Password, user.Password) {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("invalid password")), fmt.Errorf("invalid password")), fmt.Errorf("invalid password")
	}
	token, err := issueSignedJwtToken(user)
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err.Error()), err), err
	}
	return makeResponse(200, token, nil), nil
}

func getUserByEmail(search string) (user *User) {
	repo.db = initDatabase()
	repo.db.Where("username = ? OR email = ?", search, search).First(&user)
	return user
}

func validatePassword(password string, passwordHash string) bool {
	// sha256 the password
	sha256Password := sha256.Sum256([]byte(password))
	passwordToCheck := fmt.Sprintf("%x", sha256Password)
	return passwordToCheck == passwordHash
}

func issueSignedJwtToken(user *User) (string, error) {
	secret := []byte(os.Getenv("HMAC_SECRET"))
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": user.Username,
		"nbf":  time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
