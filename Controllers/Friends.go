package Controllers

import (
	"encrypted-chat/Config"
	"encrypted-chat/Localize"
	"encrypted-chat/Models"
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func AddFriendIndex(c *gin.Context) {
	var context = make(gin.H)
	localizer := Localize.GetLocalizer(c)
	context["KeywordT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:  "Keyword",
		One: "Keyword",
	})
	context["FriendT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Friend",
		One:   "Friend",
		Other: "Friends",
	})
	context["AddedFriendsT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "AddedFriends",
		Other: "Added Friends",
	})
	context["SearchFdInstructionT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "SearchFdInstruction",
		Other: "Please search your friends by name, email or id",
	})
	context["SkipT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Skip",
		Other: "Skip",
	})
	context["AddT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Add",
		Other: "Add",
	})
	context["Theme"] = Config.HackerTheme
	c.HTML(http.StatusOK, "add_friend", context)
}

func SearchFriends(c *gin.Context) {
	searchKeyword, _ := c.GetQuery("search")
	claims := jwt.ExtractClaims(c)
	var users []Models.User
	if err := Models.SearchUser(&users, searchKeyword, claims["id"].(string)); err == nil {
		c.JSON(http.StatusOK, gin.H{"friends": users})
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}
}

func AddFriends(c *gin.Context) {
	fdIds := c.PostFormArray("fd[]")
	claims := jwt.ExtractClaims(c)
	var selfUser Models.User
	if err := Models.GetUserByEmail(&selfUser, claims["id"].(string)); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	var fds []Models.User
	log.Println(fdIds)
	if err := Models.GetUsersByIds(&fds, fdIds); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}
	log.Println(selfUser, fds)
	var failedChatrooms []Models.ChatroomRelation
	var failedFds []Models.User
	baseId, _ := Models.GetBaseChatroomId()

	for _, value := range fds {
		log.Println(value, baseId)
		if value.Email != selfUser.Email {
			var selfChatroom Models.ChatroomRelation
			selfChatroom.ChatroomID = baseId
			selfChatroom.Name = value.GivenName
			selfChatroom.Bio = ""
			selfChatroom.IconURL = value.IconURL
			selfChatroom.GuestId = value.ID
			if err := Models.CreateChatroom(&selfChatroom); err != nil {
				log.Println(err)
				failedChatrooms = append(failedChatrooms, selfChatroom)
				break
			}
			log.Println("selfChatroom created")
			var fdChatroom Models.ChatroomRelation
			fdChatroom.ChatroomID = baseId
			fdChatroom.Name = selfUser.GivenName
			fdChatroom.Bio = ""
			fdChatroom.IconURL = selfUser.IconURL
			fdChatroom.GuestId = selfUser.ID
			if err := Models.CreateChatroom(&fdChatroom); err != nil {
				log.Println(err)
				failedChatrooms = append(failedChatrooms, fdChatroom)
				break
			}
			log.Println("fdChatroom created")

			var chatrooms = []Models.ChatroomRelation{fdChatroom, fdChatroom}

			if err := Models.AssociateChatrooms(&selfUser, &chatrooms); err != nil {
				log.Println(err)
				failedChatrooms = append(failedChatrooms, fdChatroom)
				failedFds = append(failedFds, value)
				break
			}
			baseId += 1
		}
	}
	if len(failedChatrooms) > 0 {
		log.Println("failedChatrooms")
		failedFds = append(failedFds, selfUser)
		if err := Models.RemoveChatrooms(&failedFds, &failedChatrooms); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}
	selfUser.Status = 1
	if err := Models.UpdateUser(&selfUser); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	c.Redirect(http.StatusMovedPermanently, "/chats")
}
