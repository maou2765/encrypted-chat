//Models/UserModel.go

package Models

import (
	"database/sql/driver"

	"gorm.io/gorm"
)

type language string

const (
	ZH_HK language = "zh_hk"
	EN_US language = "en_us"
)

func (e *language) Scan(value interface{}) error {
	*e = language(value.([]byte))
	return nil
}
func (e language) Value() (driver.Value, error) {
	return string(e), nil
}

type User struct {
	gorm.Model
	GivenName string   `form:"given_name" json:"given_name" binding:"required,alphaunicode"`
	SurnName  string   `form:"surnname"  json:"surnname" binding:"omitempty"`
	IconURL   string   `form:"icon_url"  json:"icon_url" binding:"omitempty"`
	Bio       string   `form:"bio"  json:"bio" binding:"omitempty"`
	Email     string   `form:"email"  json:"email" binding:"required,email"`
	Password  string   `form:"password"  json:"password" binding:"required"`
	Language  language `form:"language"  json:"language" gorm:"type:language" binding:"omitempty"`
}

func (b *User) TableName() string {
	return "user"
}

const UserIdentityKey = "id"
