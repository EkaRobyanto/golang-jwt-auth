package main

import (
	"golang-auth/controllers"
	"golang-auth/initializers"
	"golang-auth/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	r := gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "ping",
	// 	})
	// })
	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/profile", middleware.RequireAuth, controllers.GetProfile)
	r.Run()
}
