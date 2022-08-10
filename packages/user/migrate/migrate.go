package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

// Request is for auth log in
type Request struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Repository struct {
	db *gorm.DB
}

var repo Repository

var allSchemas = []interface{}{
	&User{},
	&GalleryImage{},
	&GalleryImageTags{},
}

func init() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=require", os.Getenv("DATABASE_URL"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_NAME"), os.Getenv("DATABASE_PASSWORD"))
	repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func Main(in Request) (*Response, error) {
	doMigrations()
	return makeResponse(200, []byte(`{"message": "success"}`), nil), nil
}

func doMigrations() {
	for _, schema := range allSchemas {
		err := repo.db.Migrator().AutoMigrate(schema)
		if err != nil {
			panic(err)
		}
	}
}
