//Models/Chatroom.go
package Models

import (
	"encrypted-chat/Config"
	"log"

	"gorm.io/gorm"
)

type ChatroomRelation struct {
	gorm.Model
	ChatroomID uint   `gorm:"index:,sort:desc"`
	Name       string `gorm:"type:varchar(50)"`
	Bio        string `gorm:"type:text"`
	IconURL    string `gorm:"type:varchar(255)"`
	GuestId    uint   `gorm:"index:,sort:desc"`

	User []*User `gorm:"many2many:user_chatrooms;foreignKey:ChatroomID;References:ID"`
}

func (b *ChatroomRelation) TableName() string {
	return "chatroom"
}

func GetBaseChatroomId() (uint, error) {
	var lastChatroom ChatroomRelation
	if err := Config.DB.Select("chatroom_id").Last(&lastChatroom).Error; err != nil {
		return 0, err
	}
	return lastChatroom.ChatroomID, nil
}
func CreateChatroom(chatrooms *ChatroomRelation) (err error) {
	if err = Config.DB.Create(chatrooms).Error; err != nil {
		return err
	}
	return nil
}
func AssociateChatrooms(user *User, chatrooms *[]ChatroomRelation) error {
	log.Println(chatrooms)
	if err := Config.DB.Model(user).Association("Chatrooms").Append(*chatrooms); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func RemoveChatrooms(users *[]User, chatrooms *[]ChatroomRelation) error {
	var chatroomIDs []uint
	for _, value := range *chatrooms {
		log.Println(value)
		chatroomIDs = append(chatroomIDs, value.ID)
	}
	if err := Config.DB.Model(users).Association("Chatrooms").Delete(chatroomIDs); err != nil {
		return err
	}
	Config.DB.Delete(&ChatroomRelation{}, chatroomIDs)
	return nil
}
