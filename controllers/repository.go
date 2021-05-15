package controllers

import (
	models "daisyCD/models"
	handlers "daisyCD/handlers"
	"strings"
	"fmt"
	"os"
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	git "gopkg.in/src-d/go-git.v4"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func SetRepository(c *gin.Context) {
	metadata, err := handlers.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// gorm
	db := c.MustGet("db").(*gorm.DB)

	var repository models.Repository
	if err := c.ShouldBindJSON(&repository); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, "Invalid json provided 1")
	   return
	}

	repository.UserId = metadata.UserId

	if res := db.Create(&repository); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided 2")
	   return
	} 	
			
	cloneErr := cloneRepository(repository.Url, repository.AccessUser, repository.AccessToken)
	if cloneErr != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
	}

	c.JSON(http.StatusOK, gin.H{"data": repository})
}

func DeleteRepository(c *gin.Context){
	userId, err := handlers.ExtractTokenUserID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	// gorm
	db := c.MustGet("db").(*gorm.DB)

	var repository models.Repository
	if err := c.ShouldBindJSON(&repository); err != nil {
	   c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   return
	}

	repository.UserId = userId

	if res := db.Where("url = ? AND user_id = ?", repository.Url, repository.UserId).Delete(&repository); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   	return
	} 	
	
	c.JSON(http.StatusOK, gin.H{"data": repository})
}

func GetRepositoryList(c *gin.Context) {
	userId, err := handlers.ExtractTokenUserID(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	var repositories []models.Repository
	if res := db.Where("user_id = ?", userId).Find(&repositories); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   	return
	}

	c.JSON(http.StatusOK, gin.H{"data": repositories})
}

func cloneRepository(url string, username string, token string) error {
	directory := "git/" + strings.Split(url,"/")[3] + "/" + strings.Split(url, "/")[4]
	
	options := &git.CloneOptions{
		Auth:	  &githttp.BasicAuth{
						Username: username,
						Password: token,
					},
		URL:	  url,
		Progress: os.Stdout,
	}

	r, err := git.PlainClone(directory, false, options)
	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}
	fmt.Println(commit)
	return nil
}

func PullRepository(c *gin.Context){
	chartPath := c.PostForm("chartPath")	
	url := "https://github.com/" + strings.Split(chartPath, "/")[1] + "/" + strings.Split(chartPath, "/")[2]
	var repoInfo models.Repository
	db := c.MustGet("db").(*gorm.DB)
	if res := db.Where("url = ?", url).Find(&repoInfo); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid URL")
		return
	}

	repository, err := git.PlainOpen(chartPath)
	if err != nil {
		fmt.Println(err)
	}
	
	worktree, err := repository.Worktree()
	if err != nil {
		panic(err)
	}
	
	options := &git.PullOptions{
		Auth:	  &githttp.BasicAuth{
						Username: repoInfo.AccessUser,
						Password: repoInfo.AccessToken,
					},
		RemoteName: "origin",
		Progress: os.Stdout,
	}

	err = worktree.Pull(options)
	if err != nil {
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{"data": repository})

}