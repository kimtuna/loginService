package models

import "time"

// User 모델 정의
type User struct {
	ID                    uint   `gorm:"primaryKey;autoIncrement"`
	Email                 string `gorm:"type:varchar(255);uniqueIndex"`
	Hash                  string
	RefreshTokenExpiresAt time.Time
}

// RefreshToken 모델 정의
type RefreshToken struct {
	Token     string `gorm:"primaryKey"`
	Email     string
	ExpiresAt time.Time
}
