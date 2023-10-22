package entity

import "gorm.io/gorm"

// 유저의 정보를 저장하는 테이블입니다
type UserEntity struct {
	gorm.Model
	Uuid                 string `gorm:"type:varchar(36);unique"`
	KakaoId              *int
	KakaoEmail           *string
	KakaoNickname        *string
	KakaoProfileImgUrl   *string
	KakaoThumbnailImgUrl *string
}

func (UserEntity) TableName() string {
	return "user"
}
