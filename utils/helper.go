package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func AbortWithErrJson(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, EJ(err))
}

func AbortWithStrJson(c *gin.Context, code int, str string) {
	c.AbortWithStatusJSON(code, SJ(str))
}

func EJ(e error) gin.H {
	return gin.H{
		"error": e.Error(),
	}
}

func SJ(str string) gin.H {
	return gin.H{
		"error": str,
	}
}

func StringToTime(t string) (*time.Time, error) {
	nt, err := time.Parse(time.RFC3339, t)

	if err != nil {
		return nil, err
	}

	return &nt, nil
}

func GetBearerToken(h string) (*string, error) {
	parts := strings.SplitN(h, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, errors.New("incorrect bearer token type")
	}

	token := parts[1]

	return &token, nil
}
