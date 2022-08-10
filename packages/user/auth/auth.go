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

type Token struct {
	Token string `json:"token"`
}

func Main(in Request) (*Response, error) {
	token := uuid.New().String()
	tokenResponse, err := json.Marshal(Token{token})
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err.Error()), err), err
	}
	return makeResponse(200, tokenResponse, nil), nil
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
