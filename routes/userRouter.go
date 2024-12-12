package routes

import (
	controller "github.com/lucapierini/UserAuthentication/controllers"
	"github.com/gin-gonic/gin"
	"github.com/lucapierini/UserAuthentication/middleware"
)

func UserRouters(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controller.GetUsers())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
}