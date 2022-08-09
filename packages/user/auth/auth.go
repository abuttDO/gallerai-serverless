package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

// Request is for auth log in
type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response returns back the http code, type of data, and the presigned url to the auth.
type Response struct {
	// StatusCode is the http code that will be returned back to the auth.
	StatusCode int `json:"statusCode,omitempty"`
	// Headers is the information about the type of data being returned back.
	Headers map[string]string `json:"headers,omitempty"`
	// Body will contain the token.
	Body string `json:"body,omitempty"`
}

type Token struct {
	Token string `json:"token"`
}

func Main(in Request) (*Response, error) {
	token := uuid.New().String()
	tokenResponse, err := json.Marshal(Token{token})
	if err != nil {
		return &Response{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, err
	}
	return &Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(tokenResponse),
	}, nil
}

func getUserByEmail(email string) (user *User) {
	repo.db.First(&user)
	return user
}

func validatePassword(password string, passwordHash string) bool {
	// sha256 the password
	sha256Password := sha256.Sum256([]byte(password))
	passwordToCheck := fmt.Sprintf("%x", sha256Password)
	return passwordToCheck == passwordHash
}

func issueSignedJwtToken(user *User) (string, error) {
	return "", nil
}
