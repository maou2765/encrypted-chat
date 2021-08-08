//Config/Database.go
package Config

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var DB *gorm.DB // DBConfig represents db configuration
type DBConfig struct {
	Host     string
	Port     int
	User     string
	DBName   string
	Password string
}

func BuildDBConfig() *DBConfig {
	dbConfig := DBConfig{
		Host:     "127.0.0.10",
		Port:     6603,
		User:     "root",
		Password: "my-secret-pw",
		DBName:   "encrypted_chat_db",
	}
	return &dbConfig
}
func DbURL(dbConfig *DBConfig) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)
}
