package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	handlers "daisyCD/handlers"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
	   err := handlers.TokenValid(c.Request)
	   if err != nil {
		  c.JSON(http.StatusUnauthorized, err.Error())
		  c.Abort()
		  return
	   }
	   c.Next()
	}
}