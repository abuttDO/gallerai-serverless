package main

import "time"

type GalleryImage struct {
	ID        int       `gorm:"primary_key" json:"id"`
	ImageID   string    `gorm:"column:image_id" json:"image_id"`
	UserID    string    `gorm:"column:user_id" json:"user_id"`
	URL       string    `gorm:"column:url" json:"url"`
	Thumbnail string    `gorm:"column:thumbnail" json:"thumbnail"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
