package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kimtuna/goLogin/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDataBase() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// MySQL 관련 환경 변수 설정
	DbUser := os.Getenv("DB_USER")         // 사용자 이름
	DbPassword := os.Getenv("DB_PASSWORD") // 비밀번호
	DbHost := os.Getenv("DB_HOST")         // 외부 서버의 도메인 이름
	DbPort := os.Getenv("DB_PORT")         // 외부 서버의 포트 번호
	DbName := "logindb"                    // 데이터베이스 이름

	// DSN(Data Source Name) 생성
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Println("Cannot connect to database")
		log.Fatal("connection error:", err)
	} else {
		fmt.Println("We are connected to the database")
	}

	// User 및 RefreshToken 모델에 맞게 테이블을 자동으로 마이그레이션합니다.
	DB.AutoMigrate(&models.User{}, &models.RefreshToken{})
}

// User 모델 정의
type User struct {
	Email string `gorm:"primaryKey"`
	Token string
	Hash  string
}

// 데이터베이스에서 토큰이 존재하는지 확인하는 함수
func IsTokenInDatabase(tokenString string) bool {
	var user User
	// 토큰이 데이터베이스에 존재하는지 확인
	if err := DB.Where("token = ?", tokenString).First(&user).Error; err != nil {
		return false
	}
	return true
}
