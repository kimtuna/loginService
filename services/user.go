package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimtuna/goLogin/models"
	"github.com/kimtuna/goLogin/setup"
	"github.com/kimtuna/goLogin/token"
)

// 사용자 정보 핸들러
func UserInfo(c *gin.Context) {
	// Authorization 헤더에서 토큰 추출
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
		return
	}

	// Access Token 검증 및 클레임 추출
	claims, err := token.ValidateAccessToken(authHeader)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 데이터베이스에서 사용자 조회
	var user models.User
	if err := setup.DB.Where("email = ?", claims.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 해시 필드를 제외하고 사용자 정보 반환
	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"token": user.Token,
		// 다른 필요한 필드 추가
	})
}