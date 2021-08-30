//Models/UserModel.go

package Models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	GivenName string `json:"given_name"`
	SurnName  string `json:"surnname"`
	IconURL   string `json:"icon_url"`
	Bio       string `json:"bio"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (b *User) TableName() string {
	return "user"
}

const UserIdentityKey = "id"
