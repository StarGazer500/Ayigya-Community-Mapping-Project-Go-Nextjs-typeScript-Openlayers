package routers

import (
	"github.com/StarGazer500/ayigya/controllers"
	"github.com/StarGazer500/ayigya/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(route *gin.RouterGroup) {
	route.GET("/profile", middlewares.AuthMiddleware(), controllers.Profile)
	route.POST("/register", controllers.Register)
	route.GET("/register", controllers.Register)
	route.POST("/login", controllers.LoginUser)
	route.GET("/login", controllers.LoginUser)

}


