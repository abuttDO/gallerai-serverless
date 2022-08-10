package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type Repository struct {
	db *gorm.DB
}

// Request details what is being asked to list
type Request struct {
	Username string `json:"username"`
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
	return makeResponse(200, []byte(`{"message": "success"}`), nil), nil
}
