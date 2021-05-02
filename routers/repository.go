package routers

import (
	controllers "daisyCD/controllers"
	middlewares "daisyCD/middlewares"
	"github.com/gin-gonic/gin"
)

func RepositoryRouter(e *gin.Engine) {
	e.POST("/repository/set", middlewares.TokenAuthMiddleware(), controllers.SetRepository)

	e.POST("/repository/delete", middlewares.TokenAuthMiddleware(), controllers.DeleteRepository)

	e.GET("/repository/list", middlewares.TokenAuthMiddleware(), controllers.GetRepositoryList)
}