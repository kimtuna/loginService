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
	// 쿠키에서 액세스 토큰 추출
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token missing"})
		return
	}

	// 액세스 토큰 검증 및 클레임 추출
	claims, err := token.ValidateAccessToken(accessToken)
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

	// 클레임 정보를 포함하여 사용자 정보 반환
	c.JSON(http.StatusOK, gin.H{
		"email": claims.Email,
		"id":    claims.ID,
		"name":  claims.Name,
	})
}
