package main

import (
	"fmt"
	"net/http"
)

// Request is for auth log in
type Request struct {
	Token   string `json:"token"`
	ImageID string `json:"image_id"`
	Type    string `json:"type"`
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

const (
	RequestTypeGet    = "GET"
	RequestTypePost   = "POST"
	RequestTypePut    = "PUT"
	RequestTypeDelete = "DELETE"
)

func Main(in Request) (*Response, error) {
	switch in.Type {
	case RequestTypeGet:
		fmt.Printf("GET %s\n", in.ImageID)
	case RequestTypePost:
		fmt.Printf("POST %s\n", in.ImageID)
	case RequestTypePut:
		fmt.Printf("PUT %s\n", in.ImageID)
	case RequestTypeDelete:
		fmt.Printf("DELETE %s\n", in.ImageID)
	default:
		return &Response{StatusCode: http.StatusBadRequest}, fmt.Errorf("invalid request type")
	}

	return &Response{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: "",
	}, nil
}
