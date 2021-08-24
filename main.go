package main

import (
	"encrypted-chat/Config"
	"encrypted-chat/Models"
	"encrypted-chat/Routes"
	"fmt"

	"github.com/jinzhu/gorm"
)

var err error

func main() {
	Config.DB, err = gorm.Open("mysql",
		Config.DbURL(Config.BuildDBConfig()))

	if err != nil {
		fmt.Println("status", err)
	}
	defer Config.DB.Close()
	Config.DB.AutoMigrate(&Models.User{})
	r := Routes.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
