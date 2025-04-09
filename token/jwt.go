package token

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"github.com/kimtuna/goLogin/models"
	"github.com/kimtuna/goLogin/setup"
	"gorm.io/gorm"
)

func init() {
	godotenv.Load()
}

var accessKey = []byte(os.Getenv("ACCESS_KEY"))
var refreshKey = []byte(os.Getenv("REFRESH_KEY"))

// Claims 구조체 정의
type Claims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Hash  string `json:"hash"`
	jwt.RegisteredClaims
}

// Access Token 생성
func GenerateAccessToken(email string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // 15분 만료
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessKey)
}

// Refresh Token 생성
func GenerateRefreshToken(email, hash string) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) // 7일 만료
	claims := &Claims{
		Email: email,
		Hash:  hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshKey)
}

// 데이터베이스에서 토큰을 확인하는 함수
func checkTokenInDatabase(tokenString string) error {
	// setup 패키지의 IsTokenInDatabase 함수 사용
	if !setup.IsTokenInDatabase(tokenString) {
		return errors.New("유효하지 않은 토큰입니다")
	}
	return nil
}

// "Bearer " 접두사 제거 함수
func DeleteBearer(authHeader string) string {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		log.Printf("Invalid token format: %s", authHeader)
		return ""
	}
	return tokenString
}

// Access Token 검증
func ValidateAccessToken(authHeader string) (*Claims, error) {
	// "Bearer " 접두사 제거
	tokenString := DeleteBearer(authHeader)
	if tokenString == "" {
		return nil, errors.New("Invalid token format")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return accessKey, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("유효하지 않은 액세스 토큰입니다")
	}

	// 만료 시간 검증
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("액세스 토큰이 만료되었습니다")
	}

	return claims, nil
}

func GetUserInfoFromToken(r *http.Request) (*Claims, error) {
	// Authorization 헤더에서 토큰 추출
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, errors.New("토큰이 제공되지 않았습니다")
	}

	// 토큰 유효성 검증
	claims, err := ValidateAccessToken(authHeader)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

// Access Token과 Refresh Token 생성
func GenerateTokens(db *gorm.DB, name, email string) (string, string, error) {
	// Access Token 생성
	accessTokenExpiration := time.Now().Add(15 * time.Minute) // 15분 만료
	accessClaims := &Claims{
		Name:  name,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExpiration),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(accessKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token 생성
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7일 만료
		},
	})
	refreshTokenString, err := refreshToken.SignedString(refreshKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token 데이터베이스에 저장
	refreshTokenRecord := models.RefreshToken{
		Token:     refreshTokenString,
		Email:     email,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	if err := db.Create(&refreshTokenRecord).Error; err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// Refresh Token 검증 및 새로운 Access Token 생성
func ValidateAndRefreshTokens(refreshTokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return refreshKey, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("유효하지 않은 리프레시 토큰입니다")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return "", errors.New("리프레시 토큰이 만료되었습니다")
	}

	// 새로운 Access Token 생성
	newAccessToken, err := GenerateAccessToken(claims.Email)
	if err != nil {
		return "", errors.New("새로운 액세스 토큰을 생성할 수 없습니다")
	}

	return newAccessToken, nil
}
