package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gym_friend_auth_server/controllers"
	"gym_friend_auth_server/initializers"
	"net/http"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DBConnection()
	initializers.InMemoryConnection()
	initializers.SetLogger()

}

func main() {
	r := gin.Default()

	r.Use(cors.New(
		cors.Config{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	controllers.UseRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
