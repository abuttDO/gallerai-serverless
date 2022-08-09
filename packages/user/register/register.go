package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Request is for auth log in
type Request struct {
	Username string `json:"username"`
	Email    string `json:"email"`
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
	err := validateRequest(in)
	if err != nil {
		return &Response{
			StatusCode: 400,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, err
	}
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

func validateRequest(in Request) error {
	if in.Email == "" {
		return fmt.Errorf("email is required")
	}
	if in.Password == "" {
		return fmt.Errorf("password is required")
	}
	if in.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}

func createUser(in Request) error {
	var user User
	user.Email = in.Email
	// sha256 the password in to a string
	sha256Password := sha256.Sum256([]byte(in.Password))
	user.Password = fmt.Sprintf("%x", sha256Password)
	user.Username = in.Username
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return nil
}
