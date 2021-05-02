package controllers

import (
	models "daisyCD/models"
	"log"
    "os"
	"fmt"
    "helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
    // "helm.sh/helm/v3/pkg/cli"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	// "k8s.io/client-go/rest"
	"path/filepath"
	"strings"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

func ListHelmChart(c *gin.Context){
	var files []string

    root := "/Users/isabel/Documents/ijnuemik-github/daisyCD/git-repos/"
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.Name() == "Chart.yaml" { 
			path = strings.Replace(path, root, "", 1)
			path = strings.Replace(path, "Chart.yaml", "", 1)
			files = append(files, path)
		}
        return nil
    })
    if err != nil {
        panic(err)
    }
    for _, file := range files {
        fmt.Println(file)
    }
	
	c.JSON(http.StatusOK, gin.H{"data": files})
}

func getActionConfig(namespace string) (*action.Configuration, error){
	actionConfig := new(action.Configuration)
	var kubeConfig *genericclioptions.ConfigFlags
	// Create the rest config instance with ServiceAccount values loaded in them
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	return nil, err
	// }
	kubeconfigPath := "/Users/isabel/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load kubeconfig from %s: %v", kubeconfigPath, err)
	}
	
	// Create the ConfigFlags struct instance with initialized values from ServiceAccount
	kubeConfig = genericclioptions.NewConfigFlags(false)
	kubeConfig.APIServer = &config.Host
	kubeConfig.BearerToken = &config.BearerToken
	kubeConfig.CAFile = &config.CAFile
	kubeConfig.Namespace = &namespace
	if err := actionConfig.Init(kubeConfig, namespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		return nil, err
	}
		
	return actionConfig, nil
}

func loadHelmchart(chartPath string) (*chart.Chart, error){
	chart, err := loader.Load(chartPath)
	if err != nil{
		return nil, err
	}
	return chart, nil
} 

func InstallHelmChart(c *gin.Context){
	var metadata models.Application
	if err := c.ShouldBindJSON(&metadata); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chart, err := loadHelmchart("/Users/isabel/Documents/ijnuemik-github/daisyCD/git-repos/" + metadata.ChartPath)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Cannot Load Helm Chart")
		fmt.Println(err)
		return
	}

	actionConfig, err := getActionConfig(metadata.ReleaseNamespace)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Cannot Connect Kubernetes")
		fmt.Println(err)
		return
	}
	client := action.NewInstall(actionConfig)
	client.Namespace = metadata.ReleaseNamespace
	client.ReleaseName = metadata.ReleaseName
	release, err := client.Run(chart, nil)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, "Cannot Install Chart")
		return
	}
	fmt.Println("Successfully installed release: ", release.Name)

	db := c.MustGet("db").(*gorm.DB)
	if res := db.Create(&metadata); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
	   return
	} 	

	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

func GetApplicationList(c *gin.Context){
	db := c.MustGet("db").(*gorm.DB)
	var applications []models.Application
	if res := db.Find(&applications); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid")
	   	return
	}

	c.JSON(http.StatusOK, gin.H{"data": applications})
}

func DeleteApplication(c *gin.Context){
	db := c.MustGet("db").(*gorm.DB)
	var application models.Application

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided 1")
		return
	}
	
	actionConfig, err := getActionConfig(application.ReleaseNamespace)
	if err != nil{
		c.JSON(http.StatusUnprocessableEntity, "Cannot Connect Kubernetes")
		fmt.Println(err)
		return
	}
	
	client := action.NewUninstall(actionConfig)
	release, err := client.Run(application.ReleaseName)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnprocessableEntity, "Cannot Uninstall Chart")
		return
	}
	fmt.Println("Successfully uninstalled release: ", release.Release.Name)	

	if res := db.Where("release_name = ?", application.ReleaseName).Delete(&application); res.Error != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided 2")
	   	return
	} 

	c.JSON(http.StatusOK, gin.H{"data": application})
}
	// func updateHelm(c *gin.Context){

// }