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
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file")
    }

    setup.ConnectDataBase()

    r := gin.Default()

    board := r.Group("/board")
    {
        api := board.Group("/api")
        {
            auth := api.Group("/auth")
            {
                // 구글 oauth
                auth.GET("/google/login", services.GoogleLogin)
                auth.GET("/google/callback", services.GoogleCallback)

                // 서비스 내부 로그인 회원가입
                auth.POST("/register", services.Register)
                auth.POST("/login", services.Login)

                auth.GET("/user", services.UserInfo)
            }
        }
    }

    // 포트 설정
    port := os.Getenv("PORT")
    if port == "" {
        log.Fatal("DOCKER_PORT environment variable not set")
    }

    r.Run(":" + port) // 서버 실행
}
