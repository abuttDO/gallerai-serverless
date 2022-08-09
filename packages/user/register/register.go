package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
		return nil, err
	}
	token := uuid.New().String()
	tokenResponse, err := json.Marshal(Token{token})
	if err != nil {
		return nil, err
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
