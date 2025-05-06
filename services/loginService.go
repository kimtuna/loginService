package services

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kimtuna/goLogin/models"
	"github.com/kimtuna/goLogin/setup"
	"github.com/kimtuna/goLogin/token"
	"golang.org/x/crypto/bcrypt"
)

// 로그인 핸들러
func Login(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// 데이터베이스에서 사용자 정보 조회
	var user models.User
	if err := setup.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// 비밀번호 검증
	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// JWT 토큰 생성
	accessToken, refreshToken, err := token.GenerateTokens(setup.DB, user.ID, "User Name", loginData.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// 환경 변수에서 도메인 가져오기
	domain := os.Getenv("DOMAIN")

	// 쿠키 설정
	c.SetCookie("access_token", accessToken, 3600, "/", domain, true, true)
	c.SetCookie("refresh_token", refreshToken, 7*24*3600, "/", domain, true, true)

	// 성공 응답
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}
