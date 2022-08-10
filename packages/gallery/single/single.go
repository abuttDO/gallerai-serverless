package main

import (
	"fmt"
)

// Request is for auth log in
type Request struct {
	Token   string `json:"token"`
	ImageID string `json:"image_id"`
	Type    string `json:"type"`
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
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("invalid request type")), fmt.Errorf("invalid request type")), fmt.Errorf("invalid request type")
	}

	return makeResponse(200, nil, nil), nil
}
