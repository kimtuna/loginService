package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	services "github.com/kimtuna/goLogin/services"
	setup "github.com/kimtuna/goLogin/setup"
)

func main() {
	// 환경 변수 로드
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	setup.ConnectDataBase()

	r := gin.Default()
	// 구글 oauth
	r.GET("/auth/google/login", services.GoogleLogin)
	r.GET("/auth/google/callback", services.GoogleCallback)

	// 서비스 내부 로그인 회원가입
	public := r.Group("/api/auth")
	public.POST("/register", services.Register)
	public.POST("/login", services.Login)

	// 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	protected := r.Group("/api/auth")
	protected.GET("/user", services.UserInfo)

	r.Run(":" + port) // 서버 실행
}
