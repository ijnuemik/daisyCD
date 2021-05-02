package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	models "daisyCD/models"
	controllers "daisyCD/controllers"
	commons "daisyCD/commons"
	middlewares "daisyCD/middlewares"
	routers "daisyCD/routers"
)

func main() {
	r := gin.Default()

	commons.RedisInit()
	commons.ViperInit()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))

	db := models.SetupModels() // new
	// Provide db variable to controllers
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	routers.AuthRouter(r)
	routers.RepositoryRouter(r)
	routers.HelmRouter(r)
	r.GET("/users", middlewares.TokenAuthMiddleware(), controllers.GetUsers)
	r.Run()

}
