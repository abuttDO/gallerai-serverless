package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type User struct {
	ID        int       `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"column:username" json:"username"`
	Email     string    `gorm:"column:email" json:"email"`
	Password  string    `gorm:"column:password" json:"password"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

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

type Repository struct {
	db *gorm.DB
}

var repo Repository

func init() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func Main(in Request) (*Response, error) {
	err := repo.db.Migrator().AutoMigrate(&User{})
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
		Body: "All done",
	}, nil
}
