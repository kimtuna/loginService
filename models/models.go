package models

import "time"

// User 모델 정의
type User struct {
	ID                    uint      `gorm:"primaryKey;autoIncrement"`
	Name                  string    `gorm:"type:varchar(255)"`
	Email                 string    `gorm:"type:varchar(255);uniqueIndex"`
	Hash                  string
	RefreshTokenExpiresAt time.Time
}

// RefreshToken 모델 정의
type RefreshToken struct {
	Token     string `gorm:"primaryKey"`
	Email     string
	ExpiresAt time.Time
}

type OAuthUser struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	Email        string `gorm:"type:varchar(255);uniqueIndex"`
	GoogleID     string `gorm:"type:varchar(255);uniqueIndex"`
	RefreshToken string `gorm:"type:text"` // 리프레시 토큰 저장
	CreatedAt    time.Time
}
