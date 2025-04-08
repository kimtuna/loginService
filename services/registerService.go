package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimtuna/goLogin/security"
	s "github.com/kimtuna/goLogin/setup"
	"github.com/kimtuna/goLogin/token"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Email string `gorm:"primaryKey"`
	Token string
	Hash  string
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 데이터베이스에서 사용자 Email이 이미 존재하는지 확인
	var existingUser User
	if err := s.DB.First(&existingUser, "email = ?", req.Email).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// 비밀번호 해시
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// 토큰 생성
	accessToken, refreshToken, err := token.GenerateTokens(s.DB, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// 사용자 정보 데이터베이스에 저장
	newUser := User{
		Email: req.Email,
		Token: refreshToken,
		Hash:  hashedPassword,
	}

	if err := s.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "User registered successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
