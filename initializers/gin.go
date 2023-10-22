package initializers

import (
	"github.com/gin-gonic/gin"
	log "github.com/shyunku-libraries/go-logger"
)

func SetLogger() {
	gin.DefaultWriter = log.GetLogger()
	gin.DefaultErrorWriter = log.GetLogger()
}
