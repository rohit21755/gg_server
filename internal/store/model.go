package store

import "time"

type User struct {
	ID       uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	// Posts     []Post
	CreatedAt time.Time
}
