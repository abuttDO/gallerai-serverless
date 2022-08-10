package main

import (
	"gorm.io/gorm"
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

func Main(in Request) (*Response, error) {
	doMigrations()
	return makeResponse(200, []byte(`{"message": "success"}`), nil), nil
}

func doMigrations() {
	repo.db = initDatabase()
	for _, schema := range allSchemas {
		err := repo.db.Migrator().AutoMigrate(schema)
		if err != nil {
			panic(err)
		}
	}
}
