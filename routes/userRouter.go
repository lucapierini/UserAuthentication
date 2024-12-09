package routes

import (
	controller "github.com/lucapierini/UserAutentication/controllers"
	"github.com/gin-gonic/gin"
	"github.com/lucapierini/UserAutentication/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}