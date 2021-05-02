package routers

import (
	"github.com/gin-gonic/gin"
	handlers "daisyCD/handlers"
	middlewares "daisyCD/middlewares"
)

func AuthRouter(e *gin.Engine) {
	e.POST("/login", handlers.Login)

	e.POST("/logout", middlewares.TokenAuthMiddleware(), handlers.Logout)

	e.POST("/token/refresh", handlers.Refresh)
}