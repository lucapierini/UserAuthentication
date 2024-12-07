package routes

import (
	
	"github.com/gin-gonic/gin"
	controller "github.com/lucapierini/UserAutentication/controllers"
)

func AuthRouters(incomingRoutes *gin.Engine) {
	auth := incomingRoutes.Group("/users")
	{
		auth.POST("/singup", controller.Singup())
		auth.POST("/login", controller.Login())
	}
}