package controllers

import (
	models "daisyCD/models"
	handlers "daisyCD/handlers"
	"net/http"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GET /books
// Get all books
func GetUsers(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	
	var users []models.User
	db.Find(&users)

	tokenAuth, err := handlers.ExtractTokenMetadata(c.Request)
	if err != nil {
	   c.JSON(http.StatusUnauthorized, "unauthorized")
	   return
	}

    userId, err := handlers.FetchAuth(tokenAuth)
	if err != nil {
	   c.JSON(http.StatusUnauthorized, "unauthorized")
	   return
	}

	fmt.Println("userId is\t", userId)

	c.JSON(http.StatusOK, gin.H{"data": users})
}