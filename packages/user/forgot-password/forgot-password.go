package main

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

// Request is for auth log in
type Request struct {
	UsernameOrEmail string `json:"username_or_email"`
	Type            string `json:"type"`
}

const (
	RequestTypeGet  = "GET"
	RequestTypePost = "POST"
)

var (
	TokenExpiry       = time.Now().Add(time.Hour * 24)
	InvalidTokenError = fmt.Errorf("invalid token")
)

func Main(in Request) (response *Response, err error) {
	switch in.Type {
	case RequestTypeGet:
		fmt.Printf("GET %s\n", in.UsernameOrEmail)
		valid := validateForgotPasswordToken(in.UsernameOrEmail)
		if !valid {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, InvalidTokenError), InvalidTokenError), InvalidTokenError
		}
		return makeResponse(200, nil, nil), nil
	case RequestTypePost:
		fmt.Printf("POST %s\n", in.UsernameOrEmail)
		token := createForgotPasswordToken(in.UsernameOrEmail)
		err = sendForgotPasswordEmail(token)
		if err != nil {
			return makeResponse(500, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		return makeResponse(200, nil, nil), nil
	}
	return makeResponse(200, nil, nil), nil
}

func createForgotPasswordToken(usernameOrEmail string) (forgotPassword ForgotPassword) {
	forgotPassword.UsernameOrEmail = usernameOrEmail
	forgotPassword.Token = uuid.New().String()
	forgotPassword.Expiry = TokenExpiry
	return
}

func sendForgotPasswordEmail(forgotPassword ForgotPassword) error {
	return nil
}

func validateForgotPasswordToken(token string) (valid bool) {
	repo.db = initDatabase()
	var forgotPassword ForgotPassword
	repo.db.Where("token = ?", token).First(&forgotPassword)
	if forgotPassword.Token == "" {
		return false
	}
	if forgotPassword.Expiry.Before(time.Now()) {
		return false
	}
	return
}
