package main

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Request is for user log in
type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Response returns back the http code, type of data, and the presigned url to the user.
type Response struct {
	// StatusCode is the http code that will be returned back to the user.
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