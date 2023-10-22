package crypto

import (
	"github.com/golang-jwt/jwt/v5"
	log "github.com/shyunku-libraries/go-logger"
	"gym_friend_auth_server/entity"
	"gym_friend_auth_server/initializers"
	"gym_friend_auth_server/model"
	"os"
	"strconv"
	"time"
)

// jwt관련 함수 집합입니다
func SaveTokens(uuid string, token string) error {
	if err := initializers.InMemoryDB.Set(uuid, token); err != nil {
		return err
	}

	log.Info("Token saved in in-memory DB / [uuid]:", uuid)
	return nil
}

func DeleteTokens(uuid string) error {
	if err := initializers.InMemoryDB.Del(uuid); err != nil {
		return err
	}

	log.Info("Token deleted in in-memory DB / [uuid]:", uuid)
	return nil
}

func GenerateTokens(e *entity.UserEntity) (*model.Tokens, error) {
	atKey := os.Getenv("JWT_ACCESS_SECRET")
	rtKey := os.Getenv("JWT_REFRESH_SECRET")
	atExp, _ := strconv.Atoi(os.Getenv("JWT_ACCESS_EXP_DATE"))
	rtExp, _ := strconv.Atoi(os.Getenv("JWT_REFRESH_EXP_DATE"))

	atClaims := model.AccessTokenClaims{
		Uuid:      e.Uuid,
		CreatedAt: e.CreatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(atExp))),
		},
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	atString, err := at.SignedString([]byte(atKey))
	if err != nil {
		return nil, err
	}

	rtClaims := model.RefreshTokenClaims{
		Id:        e.ID,
		Uuid:      e.Uuid,
		CreatedAt: e.CreatedAt,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(rtExp))),
		},
	}

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rtString, err := rt.SignedString([]byte(rtKey))
	if err != nil {
		return nil, err
	}

	token := model.Tokens{
		AccessToken:  atString,
		RefreshToken: rtString,
	}

	log.Info("Auth token generated / [uuid]:", e.Uuid)
	return &token, nil
}

func ValidateAccessToken(at string) (*model.AccessTokenClaims, error) {
	claim := &model.AccessTokenClaims{}

	_, err := jwt.ParseWithClaims(at, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return claim, err
	}
	return claim, nil
}

func ValidateRefreshToken(rt string) (*model.RefreshTokenClaims, error) {
	claim := &model.RefreshTokenClaims{}

	_, err := jwt.ParseWithClaims(rt, claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_REFRESH_SECRET")), nil
	})
	if err != nil {
		return claim, err
	}
	return claim, nil
}
