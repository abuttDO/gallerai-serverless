package main

import (
	"gorm.io/gorm"
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

type GalleryImage struct {
	ID          int       `gorm:"primary_key" json:"id"`
	Title       string    `gorm:"column:title" json:"title"`
	Description string    `gorm:"column:description" json:"description"`
	UserID      string    `gorm:"column:user_id" json:"user_id"`
	URL         string    `gorm:"column:url" json:"url"`
	Thumbnail   string    `gorm:"column:thumbnail" json:"thumbnail"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type GalleryImageTags struct {
	ID      int    `gorm:"primary_key" json:"id"`
	ImageID string `gorm:"column:image_id" json:"image_id"`
	Tag     string `gorm:"column:tag" json:"tag"`
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

type ForgotPassword struct {
	ID              int       `gorm:"primary_key" json:"id"`
	UsernameOrEmail string    `gorm:"column:username_or_email" json:"username_or_email"`
	Token           string    `gorm:"column:token" json:"token"`
	Expiry          time.Time `gorm:"column:expiry" json:"expiry"`
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
}

type Reaction struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	UserID    int    `gorm:"column:user_id" json:"user_id"`
	ImageID   int    `gorm:"column:image_id" json:"image_id"`
	Reaction  string `gorm:"column:reaction" json:"reaction"`
	CreatedAt string `gorm:"column:created_at" json:"created_at"`
}
