//Controllers/User.go

package Controllers

import (
	"encrypted-chat/Localize"
	"encrypted-chat/Models"
	"encrypted-chat/Validator"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
)

func GetSignupPageTranslation(context *gin.H, localizer *i18n.Localizer) {
	(*context)["GivenNameT"] = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "GivenName",
			Other: "Given Name",
		},
	})
	(*context)["SurnnameT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Surnname",
		Other: "Surnname",
	})
	(*context)["BioT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Bio",
		Other: "Bio",
	})
	(*context)["EmailT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Email",
		Other: "Email",
	})
	(*context)["PasswordT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Password",
		Other: "Password",
	})
	(*context)["ConfirmPasswordT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "ConfirmPassword",
		Other: "Confirm Password",
	})
}

func PasswordValidator(pw string) bool {
	var passwordRegexp = regexp.MustCompile(`(?=[a-z]){8,45}`)
	return passwordRegexp.MatchString(pw)
}

func LoginIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{})
}
func SignupIndex(c *gin.Context) {
	localizer := Localize.GetLocalizer(c)
	var context = make(gin.H)
	GetSignupPageTranslation(&context, localizer)
	fmt.Println(context["confirmPasswordT"])
	c.HTML(http.StatusOK, "signup", context)
}
func Signup(c *gin.Context) {
	binding.Validator = new(Validator.DefaultValidator)
	var user Models.User
	var ve = c.ShouldBind(&user)
	if ve == nil {
		c.ShouldBindJSON(user)
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), DefaultCost)
		if hashErr != nil {
			log.Println(hashErr.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		user.Password = string(hash)

		createError := Models.CreateUser(&user)
		if createError != nil {
			fmt.Println(createError.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			c.JSON(http.StatusOK, user)
		}
	} else {
		var context = make(gin.H)
		c.ShouldBindJSON(user)
		localizer := Localize.GetLocalizer(c)
		GetSignupPageTranslation(&context, localizer)
		context["GivenName"] = user.GivenName
		context["SurnName"] = user.SurnName
		context["IconURL"] = user.IconURL
		context["Bio"] = user.Bio
		context["Email"] = user.Email
		context["Password"] = ""
		for _, fieldErr := range ve.(validator.ValidationErrors) {
			context[fieldErr.StructField()+"Err"] = fieldErr.ActualTag()
		}
		fmt.Println(context)
		fmt.Println(context["GivenNameT"])
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
