package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/google/uuid"
	"os"
	"time"
)

// Request is for auth log in
type Request struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
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
		updatePassword(in.UsernameOrEmail, in.Password)
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
	body := fmt.Sprintf("Hi %s,\n\nYou recently requested to reset your password for your account.\n\nPlease click the following link to reset your password:\n\n%sreset-password/%s\n\nIf you did not request a password reset, please ignore this email.\n\nThanks,\nThe Go-Auth Team", os.Getenv("APP_RUL"), forgotPassword.UsernameOrEmail, forgotPassword.Token)
	fmt.Printf("%s\n", body)
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

func updatePassword(usernameOrEmail string, password string) {
	var user User
	repo.db.Where("username = ? OR email = ?", usernameOrEmail).First(&user)
	sha256Password := sha256.Sum256([]byte(password))
	user.Password = fmt.Sprintf("%x", sha256Password)
	repo.db.Save(&user)
}
