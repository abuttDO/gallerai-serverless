package main

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func makeResponse(statusCode int, body interface{}, err error) *Response {
	bodyString, err := json.Marshal(body)
	if err != nil {
		fmt.Printf("error marshalling body: %s\n", err)
		return &Response{}
	}
	var response Response
	response.StatusCode = statusCode
	response.Headers = map[string]string{
		"Content-Type": "application/json",
	}
	response.Body = string(bodyString)
	return &response
}

func initDatabase() (db *gorm.DB) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	if os.Getenv("DATABASE_TLS") == "false" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
