package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	Services "github.com/kimtuna/goLogin/Services"
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

	public := r.Group("/api/auth")

	public.POST("/register", Services.Register)
	public.POST("/login", Services.Login)

	// 포트 설정
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	protected := r.Group("/api/auth")
	protected.GET("/user", Services.UserInfo)

	r.Run(":" + port) // 서버 실행
}
