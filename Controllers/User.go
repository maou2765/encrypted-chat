//Controllers/User.go

package Controllers

import (
	"encrypted-chat/Models"
	"encrypted-chat/Validator"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

func LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}
func SignupIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "signup", gin.H{})
}
func Signup(c *gin.Context) {
	var userValidator struct {
		GivenName string `form:"givenName" json:"given_name" binding:"required,alpha"`
		SurnName  string `form:"surnname" json:"surnname" binding:"alpha"`
		IconURL   string `form:"iconURL" json:"icon_url" binding:"omitempty,uri"`
		Bio       string `form:"bio" json:"bio" binding:"omitempty,alpha"`
		Email     string `form:"email" json:"email" bind:"omitempty,email"`
		Password  string `json:"password"`
	}
	binding.Validator = new(Validator.DefaultValidator)
	var ve = c.ShouldBind(&userValidator)
	if ve == nil {
		var user Models.User
		c.ShouldBindJSON(user)
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), DefaultCost)
		if hashErr != nil {
			log.Println(hashErr.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		user.Password = string(hash)
		log.Println(user.Password)

		createError := Models.CreateUser(&user)
		if createError != nil {
			fmt.Println(createError.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		var context = make(gin.H)
		c.ShouldBindJSON(userValidator)
		context["GivenName"] = userValidator.GivenName
		context["SurnName"] = userValidator.SurnName
		context["IconURL"] = userValidator.IconURL
		context["Bio"] = userValidator.Bio
		context["Email"] = userValidator.Email
		context["Password"] = ""
		for _, fieldErr := range ve.(validator.ValidationErrors) {
			context[fieldErr.StructField()+"Err"] = fieldErr.ActualTag()
		}
		fmt.Println(context)
		c.HTML(http.StatusOK, "signup", context)
	}

}

//GetUsers ... Get all users
func GetUsers(c *gin.Context) {
	var user []Models.User
	err := Models.GetAllUsers(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func CreateUser(c *gin.Context) {
	var user Models.User
	c.BindJSON(&user)
	log.Println(user.Password)
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), DefaultCost)
	if hashErr != nil {
		log.Println(hashErr.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}
	user.Password = string(hash)
	log.Println(user.Password)
	err := Models.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func GetUserByID(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Models.User
	err := Models.GetUserByID(&user, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func UpdateUser(c *gin.Context) {
	id := c.Params.ByName("id")
	var user Models.User
	err := Models.GetUserByID(&user, id)
	if err != nil {
		fmt.Println(err.Error())
	}
	c.BindJSON(&user)
	err = Models.UpdateUser(&user, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(c *gin.Context) {
	var user Models.User
	id := c.Params.ByName("id")
	err := Models.DeleteUser(&user, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
