package models

import "time"

type Post struct {
	ID          uint       `json:"id" gorm:"primaryKey;column:id"`
	Title       string     `json:"title" gorm:"column:title"`
	Content     string     `json:"content" gorm:"column:content;type:longtext"`
	Thumbnail   *string    `json:"thumbnail" gorm:"column:thumbnail"`
	Published   bool       `json:"published" gorm:"column:published"`
	PublishedAt *time.Time `json:"publishedAt" gorm:"column:publishedAt"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt" gorm:"column:updatedAt"`
}

func (Post) TableName() string {
	return "post"
}