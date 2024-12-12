package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/lucapierini/UserAuthentication/controllers"
)

func AuthRouters(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
}