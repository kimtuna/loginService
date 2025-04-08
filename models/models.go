package models

import "time"

// User 모델 정의
type User struct {
	Email string `gorm:"primaryKey"`
	Token string
	Hash  string
}

// RefreshToken 모델 정의
type RefreshToken struct {
	Token     string `gorm:"primaryKey"`
	Email     string
	ExpiresAt time.Time
}
