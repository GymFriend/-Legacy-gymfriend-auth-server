package model

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type KakaoResponse struct {
	ID           int             `json:"id"`
	ConnectedAt  string          `json:"connected_at"`
	Properties   KakaoProperties `json:"properties"`
	KakaoAccount KakaoAccount    `json:"kakao_account"`
}

type KakaoTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type KakaoProperties struct {
	Nickname       string `json:"nickname"`
	ProfileImage   string `json:"profile_image"`
	ThumbnailImage string `json:"thumbnail_image"`
}

type KakaoAccount struct {
	Email string `json:"email"`
}

type AccessTokenClaims struct {
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"createdAT"`
	jwt.RegisteredClaims
}

type RefreshTokenClaims struct {
	Id        uint      `json:"id"`
	Uuid      string    `json:"uuid"`
	CreatedAt time.Time `json:"createdAT"`
	jwt.RegisteredClaims
}

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
