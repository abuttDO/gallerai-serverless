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

type Token struct {
	Token string `json:"token"`
}

func Main(in Request) (*Response, error) {
	err := validateRequest(in)
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err.Error()), err), err
	}
	token := uuid.New().String()
	tokenResponse, err := json.Marshal(Token{token})
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err.Error()), err), err
	}
	return makeResponse(200, tokenResponse, nil), nil
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
