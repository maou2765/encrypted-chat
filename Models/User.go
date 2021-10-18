//Models/User.go

package Models

import (
	"database/sql/driver"
	"encrypted-chat/Config"
	"errors"
	"fmt"
	"log"
	"reflect"

	"gorm.io/gorm"
)

type Language string

const (
	ZH_HK Language = "zh_hk"
	EN_US Language = "en_us"
)

func (e *Language) Scan(value interface{}) error {
	*e = Language(value.([]byte))
	return nil
}
func (e Language) Value() (driver.Value, error) {
	return string(e), nil
}

type User struct {
	gorm.Model
	GivenName string              `form:"given_name" json:"given_name" binding:"required,alphaunicode" gorm:"type:varchar(72)"`
	Surname   string              `form:"surname"  json:"surname" binding:"omitempty" gorm:"type:varchar(72)"`
	IconURL   string              `form:"icon_url"  json:"icon_url" binding:"omitempty" gorm:"type:varchar(255)"`
	Bio       string              `form:"bio"  json:"bio" binding:"omitempty" gorm:"type:text"`
	Email     string              `form:"email"  json:"email" binding:"required,email" gorm:"type:varchar(72);index:email_unique,uniqueIndex"`
	Password  string              `form:"password"  json:"password" binding:"required" gorm:"type:varchar(72)"`
	Language  Language            `form:"language"  json:"language" gorm:"type:ENUM('zh_hk', 'en_us')" binding:"omitempty"`
	Status    uint                `binding:"omitempty"`
	Chatrooms []*ChatroomRelation `gorm:"many2many:user_chatrooms;foreignKey:ID;References:ChatroomID"`
}

type Chatroom struct {
	ChatroomID      int64   `json:"chatroom_id"`
	ChatroomName    string  `json:"chatroom_name"`
	ChatroomBio     string  `json:"chatroom_bio"`
	ChatroomIconURL string  `json:"chatroom_icon_url"`
	Members         []int64 `json:"members"`
}

type UserDetail struct {
	ID        uint64        `json:"id"`
	GivenName string        `json:"given_name"`
	Surname   string        `json:"surname"`
	IconURL   string        `json:"icon_url"`
	Bio       string        `json:"bio"`
	Email     string        `json:"email"`
	Language  string        `json:"language"`
	Status    int64         `json:"status"`
	Chatrooms [](*Chatroom) `json:"chatrooms"`
}

func (b *User) TableName() string {
	return "user"
}

const UserIdentityKey = "id"

func GetAllUsers(userDetails *[]UserDetail) (err error) {
	var rawUsers []map[string]interface{}
	var chatrooms []Chatroom
	userSubQuery := Config.DB.Table("user").Select("id")
	if err := Config.DB.Table("(?) as chatrooms, user",
		Config.DB.Table("chatroom").Where("chatroom_id IN (?)",
			Config.DB.Table("user_chatrooms").Where("user_id IN (?)", userSubQuery).Select("chatroom_relation_chatroom_id")).Select("*")).Select([]string{
		"user.id",
		"user.given_name",
		"user.surname",
		"user.icon_url",
		"user.bio",
		"user.email",
		"user.language",
		"user.status",
		"chatrooms.chatroom_id AS chatroom_id",
		"chatrooms.name AS chatroom_name",
		"chatrooms.bio AS chatroom_bio",
		"chatrooms.icon_url AS chatroom_icon_url",
		"chatrooms.guest_id AS chatroom_guest_id"}).Order("user.id asc").Find(&rawUsers).Error; err != nil {
		return err
	}
	for _, value := range rawUsers {
		userExists := false
		for userKey, userValue := range *userDetails {
			if userValue.ID == value["id"] {
				userExists = true
				chatroomExists := false
				for chatroomKey, chatroom := range chatrooms {
					if chatroom.ChatroomID == value["chatroom_id"].(int64) {
						memberExist := false
						for _, member := range chatroom.Members {
							if userValue.ID == uint64(value["chatroom_guest_id"].(int64)) {
								for userChatroomKey, userChatroomValue := range userValue.Chatrooms {
									if userChatroomValue.ChatroomID == value["chatroom_id"].(int64) {
										break
									}
									if userChatroomKey == len((*userDetails)[userKey].Chatrooms)-1 {
										log.Println("userChatroomKey == len(userValue.Chatrooms)-1", userChatroomValue.ChatroomID == value["chatroom_id"].(int64))
										(*userDetails)[userKey].Chatrooms = append((*userDetails)[userKey].Chatrooms, &(chatrooms[chatroomKey]))
									}
								}
								if len((*userDetails)[userKey].Chatrooms) == 0 && (*userDetails)[userKey].ID == uint64(value["chatroom_guest_id"].(int64)) {
									log.Println("len(userValue.Chatrooms) == 0", userValue.ID == uint64(value["chatroom_guest_id"].(int64)))
									(*userDetails)[userKey].Chatrooms = append((*userDetails)[userKey].Chatrooms, &(chatrooms[chatroomKey]))
								}
							}
							if member == value["chatroom_guest_id"].(int64) {
								memberExist = true
								break
							}
						}
						if !memberExist {
							chatrooms[chatroomKey].Members = append(chatrooms[chatroomKey].Members, value["chatroom_guest_id"].(int64))
						}
						chatroomExists = true
						break
					}
				}
				if !chatroomExists {
					chatrooms = append(chatrooms, Chatroom{
						ChatroomID:      value["chatroom_id"].(int64),
						ChatroomName:    value["chatroom_name"].(string),
						ChatroomBio:     value["chatroom_bio"].(string),
						ChatroomIconURL: value["chatroom_icon_url"].(string),
						Members:         []int64{value["chatroom_guest_id"].(int64)},
					})
					if len(userValue.Chatrooms) == 0 && userValue.ID == uint64(value["chatroom_guest_id"].(int64)) {
						(*userDetails)[userKey].Chatrooms = append((*userDetails)[userKey].Chatrooms, &(chatrooms[len(chatrooms)-1]))
					}
				}
			}
		}
		if !userExists {
			*userDetails = append(*userDetails, UserDetail{
				ID:        value["id"].(uint64),
				GivenName: value["given_name"].(string),
				Surname:   value["surname"].(string),
				IconURL:   value["icon_url"].(string),
				Bio:       value["bio"].(string),
				Email:     value["email"].(string),
				Language:  value["language"].(string),
				Status:    value["status"].(int64),
				Chatrooms: []*Chatroom{},
			})
		}
	}
	return nil
}
func CreateUser(user *User) (err error) {
	if err = Config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByID(user *User, id string) (err error) {
	if err = Config.DB.Select([]string{"id", "given_name", "surname", "icon_url", "bio", "email", "language", "status"}).Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByEmail(user *User, email string) (err error) {
	if err = Config.DB.Where("email = ? ", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func SearchUser(user *[]User, keyword string, selfId string) (err error) {
	log.Println(keyword, selfId)
	if err = Config.DB.Select([]string{"id", "given_name", "surname", "icon_url", "bio", "email", "language", "status"}).Where("(given_name like ? OR surname like ? OR email like ? OR id = ? )AND email != ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", keyword, selfId).Find(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUsersByIds(user *[]User, ids interface{}) error {
	switch reflect.TypeOf(ids).String() {
	case "string":
		if err := Config.DB.Select([]string{"id", "given_name", "surname", "icon_url", "bio", "email", "language", "status"}).Where("id = ?", ids).Find(user).Error; err != nil {
			return err
		}
		return nil
	case "[]string":
		if err := Config.DB.Select([]string{"id", "given_name", "surname", "icon_url", "bio", "email", "language", "status"}).Where("id in ?", ids).Find(user).Error; err != nil {
			return err
		}
		return nil
	default:
		return errors.New("invalid ids")
	}
	return nil
}
func UpdateUser(user *User) (err error) {
	fmt.Println(user)
	Config.DB.Save(user)
	return nil
}

func DeleteUser(user *User, id string) (err error) {
	Config.DB.Where("id = ?", id).Delete(user)
	return nil
}

func Login(user *User, email string) (err error) {
	if err = Config.DB.Where("email = ?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}
