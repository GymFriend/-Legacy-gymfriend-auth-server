package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/shyunku-libraries/go-logger"
	"gorm.io/gorm"
	"gym_friend_auth_server/crypto"
	"gym_friend_auth_server/dto"
	"gym_friend_auth_server/entity"
	"gym_friend_auth_server/initializers"
	"gym_friend_auth_server/model"
	"gym_friend_auth_server/utils"
	"io"
	"net/http"
	"os"
	"strings"
)

// 인증 관련 라우터입니다

const (
	kakaoTokenURL = "https://kauth.kakao.com/oauth/token"
	kakaoAPIURL   = "https://kapi.kakao.com/v2/user/me"
)

// 카카오 로그인 관련 코드입니다
func KakaoLogin(c *gin.Context) {
	// 클라이언트로부터 받은 카카오 code를 파싱하는 코드입니다
	authHeader := c.Request.Header.Get("Authorization")
	code, err := utils.GetBearerToken(authHeader)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	client := &http.Client{}

	// code를 이용해 kakao access token을 요청하는 코드입니다
	postData := "grant_type=authorization_code" + "&client_id=" + os.Getenv("KAKAO_ID") +
		"&redirect_uri=" + os.Getenv("KAKAO_REDIRECT_URI") +
		"&code=" + *code

	kakaoTokenReq, err := http.NewRequest("POST", kakaoTokenURL, strings.NewReader(postData))
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	kakaoTokenReq.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	kakaoTokenRes, err := client.Do(kakaoTokenReq)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer kakaoTokenRes.Body.Close()

	kakaoTokenBody, err := io.ReadAll(kakaoTokenRes.Body)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	kakaoTokenResponse := model.KakaoTokenResponse{}

	if err = json.Unmarshal(kakaoTokenBody, &kakaoTokenResponse); err != nil {
		if err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// 카카오 access token을 이용해 유저 정보를 요청하는 코드입니다
	kakaoUserReq, err := http.NewRequest("GET", kakaoAPIURL, nil)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	kakaoUserReq.Header.Add("Authorization", "Bearer "+kakaoTokenResponse.AccessToken)

	kakaoUserRes, err := client.Do(kakaoUserReq)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer kakaoUserRes.Body.Close()

	kakaoUserBody, err := io.ReadAll(kakaoUserRes.Body)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var (
		user      *entity.UserEntity
		userDto   *dto.UserDto
		kakaoResp model.KakaoResponse
	)

	if err = json.Unmarshal(kakaoUserBody, &kakaoResp); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if kakaoResp.ID == 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	props := kakaoResp.Properties

	tx := initializers.DB.Begin()

	// 유저 정보가 DB에 존재하면 해당 유저 정보를 응답, 존재하지 않다면 DB에 유저정보 저장 후 응답하는 코드입니다
	if err = tx.Where("kakao_id = ?", kakaoResp.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			*user = entity.UserEntity{
				Uuid:                 uuid.New().String(),
				KakaoId:              &kakaoResp.ID,
				KakaoEmail:           &kakaoResp.KakaoAccount.Email,
				KakaoNickname:        &props.Nickname,
				KakaoProfileImgUrl:   &props.ProfileImage,
				KakaoThumbnailImgUrl: &props.ThumbnailImage,
			}

			if err = tx.Create(&user).Error; err != nil {
				log.Error(err)
				tx.Rollback()
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			log.Info("New Kakao user created:", user.Uuid)
		} else {
			log.Error(err)
			tx.Rollback()
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	// jwt를 생성하는 코드입니다. 로그인 이외의 서비스 내에서 발생하는 모든 요청의 header에는 access token이 포함되어 있어야 합니다.
	tokenDto, err := crypto.GenerateTokens(user)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userDto = dto.UserDtoFromEntity(*user, *tokenDto)

	if err = tx.Commit().Error; err != nil {
		log.Error(err)
		tx.Rollback()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, userDto)
	log.Info("Kakao user logged in / [uuid]:", userDto.Uuid)
}

// 자동 로그인 코드입니다 (현재 유저 정보를 브라우저에 저장하고 있으므로 사용 여부 파악 불가)
func AutoLogin(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	refreshToken := c.Request.Header.Get("X-Refresh-Token")
	accessToken, err := utils.GetBearerToken(authHeader)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var user *entity.UserEntity

	atClaims, atErr := crypto.ValidateAccessToken(*accessToken)
	rtClaims, rtErr := crypto.ValidateRefreshToken(refreshToken)
	if atClaims.Uuid != rtClaims.Uuid {
		log.Errorf("Payload mismatch / [access payload]: $s [refresh payload]: $s", atClaims.Uuid, rtClaims.Uuid)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if atErr == nil {
		initializers.DB.Where("uuid = ?", atClaims.Uuid).First(&user)
		log.Info("Valid access token / [uuid]:", atClaims.Uuid)
	} else if atErr != nil && rtErr == nil {
		initializers.DB.Where("uuid = ?", rtClaims.Uuid).First(&user)
		log.Info("Valid refresh token / [uuid]:", rtClaims.Uuid)
	} else {
		if err = crypto.DeleteTokens(rtClaims.Uuid); err != nil {
			log.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		log.Error(atErr, rtErr)
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	tokenDto, err := crypto.GenerateTokens(user)
	if err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if err = crypto.SaveTokens(user.Uuid, tokenDto.RefreshToken); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userDto := dto.UserDtoFromEntity(*user, *tokenDto)

	c.JSON(201, userDto)
	log.Info("User auto logged in / [uuid]:", userDto.Uuid)
}

// 로그아웃을 요청받는 코드입니다
func Logout(c *gin.Context) {
	var body dto.LogoutDto
	if err := c.Bind(&body); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := crypto.DeleteTokens(body.Uuid); err != nil {
		log.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(201, gin.H{})
	log.Info("User logged out / [uuid]:", body.Uuid)
}

func UseAuthRouter(g *gin.Engine) {
	sg := g.Group("/auth")

	sg.GET("/kakao", KakaoLogin)
	sg.POST("/auto-login", AutoLogin)
	sg.POST("/logout", Logout)
}
