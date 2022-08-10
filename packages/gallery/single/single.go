package main

import (
	"fmt"
)

// Request is for auth log in
type Request struct {
	Token       string `json:"token"`
	ImageID     string `json:"image_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
}

const (
	RequestTypeGet    = "GET"
	RequestTypePost   = "POST"
	RequestTypePut    = "PUT"
	RequestTypeDelete = "DELETE"
)

/*
Flow for saving a new image:
- Frontend requests a signed URL for a new image. This is the POST
- Frontend uploads the image to the URL
- Frontend sends this URL to the backend with the image ID
- Backend creates a thumbnail for the image, updates the record
*/

func Main(in Request) (response *Response, err error) {
	switch in.Type {
	case RequestTypeGet:
		fmt.Printf("GET %s\n", in.ImageID)
	case RequestTypePost:
		fmt.Printf("POST %s\n", in.ImageID)
	case RequestTypePut:
		fmt.Printf("PUT %s\n", in.ImageID)
		err = updateImageWithUrl(in.ImageID, in.URL)
	case RequestTypeDelete:
		fmt.Printf("DELETE %s\n", in.ImageID)
	default:
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, fmt.Errorf("invalid request type")), fmt.Errorf("invalid request type")), fmt.Errorf("invalid request type")
	}

	return makeResponse(200, nil, nil), nil
}

func updateImageWithUrl(imageID string, url string) error {
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	image.URL = url
	repo.db.Save(&image)
	return nil
}

func createThumbnails(imageID string) error {
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	return nil
}
