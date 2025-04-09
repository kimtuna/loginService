package services

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kimtuna/goLogin/models"
	"github.com/kimtuna/goLogin/setup"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Error exchanging token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// 자동으로 액세스 토큰을 갱신하는 TokenSource 생성
	tokenSource := googleOauthConfig.TokenSource(context.Background(), token)
	client := oauth2.NewClient(context.Background(), tokenSource)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// OAuthUser 정보 저장
	var oauthUser models.OAuthUser
	if err := setup.DB.Where("email = ?", userInfo.Email).First(&oauthUser).Error; err != nil {
		// 새로운 사용자라면 데이터베이스에 추가
		oauthUser = models.OAuthUser{
			Email:        userInfo.Email,
			GoogleID:     userInfo.Id,
			RefreshToken: token.RefreshToken, // 리프레시 토큰 저장
		}
		setup.DB.Create(&oauthUser)
	} else {
		// 기존 사용자라면 리프레시 토큰 업데이트
		oauthUser.RefreshToken = token.RefreshToken
		setup.DB.Save(&oauthUser)
	}

	c.JSON(http.StatusOK, gin.H{
		"email": userInfo.Email,
		// 다른 필요한 필드 추가
	})
}
