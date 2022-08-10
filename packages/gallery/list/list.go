package main

import "fmt"

// Request details what is being asked to list
type Request struct {
	Username string `json:"username"`
}

func Main(in Request) (*Response, error) {
	images, err := getImagesByUsername(in.Username)
	if err != nil {
		return makeResponse(500, fmt.Sprintf(`{"error": "%s"}`, err), err), err
	}
	return makeResponse(200, images, nil), nil
}

func getImagesByUsername(username string) ([]GalleryImage, error) {
	repo.db = initDatabase()
	var images []GalleryImage
	repo.db.Where("username = ?", username).Order("id desc").Find(&images)
	return images, nil
}
