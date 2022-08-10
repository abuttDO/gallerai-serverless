package main

import (
	"crypto/sha256"
	"fmt"
)

// Request is for auth log in
type Request struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

func Main(in Request) (*Response, error) {
	user := getUserByEmail(in.UsernameOrEmail)
	if user == nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("user not found")), fmt.Errorf("user not found")), fmt.Errorf("user not found")
	}
	if !validatePassword(in.Password, user.Password) {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("invalid password")), fmt.Errorf("invalid password")), fmt.Errorf("invalid password")
	}
	token, err := issueSignedJwtToken(user)
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err.Error()), err), err
	}
	var tokenResponse Token
	tokenResponse.Token = token
	return makeResponse(200, tokenResponse, nil), nil
}

func getUserByEmail(search string) (user *User) {
	repo.db = initDatabase()
	repo.db.Where("username = ? OR email = ?", search, search).First(&user)
	return user
}

func validatePassword(password string, passwordHash string) bool {
	// sha256 the password
	sha256Password := sha256.Sum256([]byte(password))
	passwordToCheck := fmt.Sprintf("%x", sha256Password)
	return passwordToCheck == passwordHash
}
