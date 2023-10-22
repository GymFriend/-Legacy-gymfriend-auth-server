package dto

import (
	"gym_friend_auth_server/entity"
	"gym_friend_auth_server/model"
	"time"
)

type UserDto struct {
	Uuid                 string       `json:"uuid"`
	KakaoId              *int         `json:"kakaoId"`
	KakaoEmail           *string      `json:"kakaoEmail"`
	KakaoNickname        *string      `json:"kakaoNickname"`
	KakaoProfileImgUrl   *string      `json:"kakaoProfileImgUrl"`
	KakaoThumbnailImgUrl *string      `json:"kakaoThumbnailImgUrl"`
	CreatedAt            time.Time    `json:"createdAt"`
	Token                model.Tokens `json:"token"`
}

func UserDtoFromEntity(entity entity.UserEntity, tokens model.Tokens) *UserDto {
	return &UserDto{
		Uuid:                 entity.Uuid,
		KakaoId:              entity.KakaoId,
		KakaoEmail:           entity.KakaoEmail,
		KakaoNickname:        entity.KakaoNickname,
		KakaoProfileImgUrl:   entity.KakaoProfileImgUrl,
		KakaoThumbnailImgUrl: entity.KakaoThumbnailImgUrl,
		CreatedAt:            entity.CreatedAt,
		Token:                tokens,
	}
}

type LogoutDto struct {
	Uuid string `json:"uuid" bind:"required"`
}
