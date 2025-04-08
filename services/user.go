package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kimtuna/goLogin/models"
	"github.com/kimtuna/goLogin/setup"
)

// 사용자 정보 핸들러
func UserInfo(c *gin.Context) {
	email := c.Param("email")

	var user models.User
	if err := setup.DB.Where("email = ?", email).First(&user).Error; err != nil {
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
