package main

import (
	"gym_friend_auth_server/entity"
	"gym_friend_auth_server/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DBConnection()
}

// 데이터베이스 migration 파일입니다
// 데이터베이스의 스키마가 변경되면 동기화 및 데이터 보존을 위해 "반드시" 실행시켜줘야 합니다
func main() {
	initializers.DB.AutoMigrate(&entity.UserEntity{})
}
