package main

import (
	"os"
	"github.com/gin-gonic/gin"
	"github.com/lucapierini/UserAuthentication/routes"
)

func main(){
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		}
	
	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRouters(router)
	routes.UserRouters(router)
	
	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to User Authentication API",
		})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to User Authentication API",
		})
	})

	router.Run(":" + port)


}