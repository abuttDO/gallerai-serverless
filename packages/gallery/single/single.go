package main

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/disintegration/imaging"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// Request is for auth log in
type Request struct {
	Token       string `json:"token"`
	ImageID     string `json:"image_id"`
	Description string `json:"description"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Filename    string `json:"filename"`
	Duration    string `json:"duration"`
}

var (
	key, secret, bucket, region string
	// ErrNoFilename will return an error if no filename is provided by the user.
	ErrNoFilename = errors.New("no filename provided")
	// ErrNoRequest will return an error if no request type is provided by the user.
	ErrNoRequest = errors.New("no request type provided")
	// ErrNoDuration will return an error if no duration is provided by the user.
	ErrNoDuration = errors.New("no duration provided")
	// ErrNegativeDuration will return an error if a negative duration is provided by the user.
	ErrNegativeDuration = errors.New("negative duration provided")
)

const (
	RequestTypeGet    = "GET"
	RequestTypePost   = "POST"
	RequestTypePut    = "PUT"
	RequestTypeDelete = "DELETE"
)

func init() {
	key = os.Getenv("SPACES_KEY")
	if key == "" {
		panic("no key provided")
	}
	secret = os.Getenv("SPACES_SECRET")
	if secret == "" {
		panic("no secret provided")
	}
	bucket = os.Getenv("BUCKET_NAME")
	if bucket == "" {
		panic("no bucket provided")
	}
	region = os.Getenv("REGION")
	if region == "" {
		panic("no region provided")
	}
}

/*
Flow for saving a new image:
- Frontend requests a signed URL for a new image. This is the POST
- Frontend uploads the image to the URL
- Frontend sends this URL to the backend with the image ID. This is the PUT
- Backend creates a thumbnail for the image, updates the record
*/

func Main(in Request) (response *Response, err error) {
	user, valid, err := validateToken(in.Token)
	if !valid && err != nil {
		return makeResponse(401, fmt.Sprintf(`{"error": "%s"}`, InvalidTokenError), InvalidTokenError), InvalidTokenError
	}
	if in.Filename == "" {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoFilename
	}

	duration, err := time.ParseDuration(in.Duration)
	if err != nil {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNoDuration
	}
	if duration < 0 {
		return &Response{StatusCode: http.StatusBadRequest}, ErrNegativeDuration
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
	}
	switch in.Type {
	case RequestTypeGet:
		image, err := getImageByImageID(in.ImageID)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		return makeResponse(200, image, nil), nil
	case RequestTypePost:
		url, err := uploadURL(sess, bucket, in.Filename, duration)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		return makeResponse(200, fmt.Sprintf(`{"url": "%s"}`, url), nil), nil
	case RequestTypePut:
		err = updateImageWithUrl(in.ImageID, in.URL)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		thumbnail, err := createThumbnails(in.ImageID)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		err = updateImageWithThumbnailUrl(in.ImageID, thumbnail)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		return makeResponse(200, fmt.Sprintf(`{"url": "%s", "thumbnail": "%s"}`, in.URL, thumbnail), nil), nil
	case RequestTypeDelete:
		fmt.Printf("DELETE %s\n", in.ImageID)
		err := deleteImageByUserNameAndImageID(user, in.ImageID)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
		err = deleteImageFromS3(in.ImageID)
		if err != nil {
			return makeResponse(400, fmt.Sprintf(`{"error": "%s"}`, err), err), err
		}
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

func updateImageWithThumbnailUrl(imageID string, url string) error {
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	image.Thumbnail = url
	repo.db.Save(&image)
	return nil
}

func createThumbnails(imageID string) (string, error) {
	// Save the original image locally from image.URL
	// Create a thumbnail for the image
	// Save the thumbnail locally
	// Upload the thumbnail to S3
	// Update the image record with the thumbnail URL
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	originalImage, err := imaging.Open(image.URL)
	if err != nil {
		return "", err
	}
	thumbnail := imaging.Thumbnail(originalImage, 200, 200, imaging.Lanczos)
	thumbnailFilename := fmt.Sprintf("%s-thumbnail.jpg", imageID)
	err = imaging.Save(thumbnail, thumbnailFilename)
	if err != nil {
		return "", err
	}
	thumbnailURL, err := uploadFile(thumbnailFilename)
	if err != nil {
		return "", err
	}
	err = updateImageWithThumbnailUrl(imageID, thumbnailURL)
	if err != nil {
		return "", err
	}
	return thumbnailURL, nil
}

func uploadURL(sess *session.Session, bucket string, filename string, duration time.Duration) (string, error) {
	client := s3.New(sess)
	req, _ := client.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	url, err := req.Presign(duration)
	if err != nil {
		return "", err
	}
	return url, nil
}

func uploadFile(filename string) (string, error) {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return "", err
	}
	client := s3.New(sess)
	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   strings.NewReader(filename),
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s.digitaloceanspaces.com:443/%s", bucket, filename), nil
}

func getImageByImageID(imageID string) (GalleryImage, error) {
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	return image, nil
}

func deleteImageByUserNameAndImageID(username string, imageID string) (err error) {
	repo.db = initDatabase()
	var image GalleryImage
	repo.db.First(&image, imageID)
	// Get user by username
	var user User
	repo.db.First(&user, username)
	// Validate that the user owns the image
	if user.ID != image.UserID {
		return fmt.Errorf("user does not own image")
	}
	repo.db.Delete(&image)
	return nil
}

func deleteImageFromS3(imageID string) (err error) {
	var image GalleryImage
	repo.db = initDatabase()
	repo.db.First(&image, imageID)

	err = doDeleteFromS3GivenFilename(image.URL)
	if err != nil {
		return err
	}
	err = doDeleteFromS3GivenFilename(image.Thumbnail)
	if err != nil {
		return err
	}
	return nil
}

func doDeleteFromS3GivenFilename(filename string) error {
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
		Endpoint:    aws.String(fmt.Sprintf("%s.digitaloceanspaces.com:443", region)),
		Region:      aws.String(region),
	}
	sess, err := session.NewSession(config)
	if err != nil {
		return err
	}
	client := s3.New(sess)
	_, err = client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return err
	}
	return nil
}
