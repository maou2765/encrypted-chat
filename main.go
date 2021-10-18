package main

import (
	"encrypted-chat/Config"
	"encrypted-chat/Models"
	"encrypted-chat/Routes"

	"fmt"

	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var err error

func SetupServer() *gin.Engine {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	Config.DB, err = gorm.Open(mysql.Open(Config.DbURL(Config.BuildDBConfig())),
		&gorm.Config{
			Logger: newLogger,
		})

	if err != nil {
		fmt.Println("status", err)
	}
	Config.DB.AutoMigrate(&Models.User{})
	migrator := Config.DB.Migrator()
	if emailUniConstraint := migrator.HasConstraint(&Models.User{}, "email_unique"); !emailUniConstraint {
		migrator.CreateConstraint(&Models.User{}, "email_unique")
	}
	r := Routes.SetupRouter()
	// Listen and Server in 0.0.0.0:8080
	return r
}
func main() {
	r := SetupServer()
	r.Run(":8080")
}
