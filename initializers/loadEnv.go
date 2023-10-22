package initializers

import (
	"github.com/joho/godotenv"
)

// 환경변수를 load하는 파일입니다
func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		panic("Failed to load env")

	}
}
