package routers

import (
	controllers "daisyCD/controllers"
	middlewares "daisyCD/middlewares"
	"github.com/gin-gonic/gin"
)

func HelmRouter(e *gin.Engine) {
	e.GET("/helm/applist", middlewares.TokenAuthMiddleware(), controllers.GetApplicationList)
	e.POST("/helm/appdelete", middlewares.TokenAuthMiddleware(), controllers.DeleteApplication)
	e.GET("/helm/chartlist", middlewares.TokenAuthMiddleware(), controllers.ListHelmChart)
	e.POST("/helm/installchart", middlewares.TokenAuthMiddleware(), controllers.InstallHelmChart)
}