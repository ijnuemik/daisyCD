
package models

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
	"github.com/spf13/viper"
)

func SetupModels() *gorm.DB {

	// viper package read .env
	viper_user := viper.Get("POSTGRES_USER")
	viper_password := viper.Get("POSTGRES_PASSWORD")
	viper_db := viper.Get("POSTGRES_DB")
	viper_host := viper.Get("POSTGRES_HOST")
	viper_port := viper.Get("POSTGRES_PORT")

	// https://gobyexample.com/string-formatting
	prosgret_conname := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=disable", viper_host, viper_port, viper_user, viper_db, viper_password)

	fmt.Println("conname is\t\t", prosgret_conname)

	db, err := gorm.Open(postgres.Open(prosgret_conname), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Repository{})
	db.AutoMigrate(&Application{})

	return db
}