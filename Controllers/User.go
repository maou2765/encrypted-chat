//Controllers/User.go

package Controllers

import (
	"encrypted-chat/Localize"
	"encrypted-chat/Middlewares"
	"encrypted-chat/Models"
	"encrypted-chat/Validator"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

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
	(*context)["SurnameT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Surname",
		Other: "Surname",
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
	localizer := Localize.GetLocalizer(c)
	context := make(gin.H)
	context["AppNameT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "AppName",
		Other: "Encrypted Chat",
	})
	context["EmailT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Email",
		Other: "Email",
	})
	context["PasswordT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "Password",
		Other: "Password",
	})
	context["SignInT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "SignIn",
		Other: "Sign In",
	})
	context["NoACT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "NoAC",
		Other: "If you don't have a account, Please ",
	})
	context["SignUpT"], _ = localizer.LocalizeMessage(&i18n.Message{
		ID:    "SignUp",
		Other: "SIGN UP",
	})
	c.HTML(http.StatusOK, "login", context)
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
		user.Language = Models.EN_US

		createError := Models.CreateUser(&user)
		if createError != nil {
			fmt.Println(createError.Error())
			c.AbortWithStatus(http.StatusInternalServerError)
		} else {
			mw, err := Middlewares.AuthMiddleware()

			// expire := mw.TimeFunc().Add(mw.Timeout)
			expireCookie := mw.TimeFunc().Add(mw.CookieMaxAge)
			maxage := int(expireCookie.Unix() - time.Now().Unix())

			jwtString, cookiesTime, err := mw.TokenGenerator(&user)
			fmt.Println(cookiesTime)
			fmt.Println(err)
			if mw.CookieSameSite != 0 {
				c.SetSameSite(mw.CookieSameSite)
			}

			c.SetCookie(
				mw.CookieName,
				jwtString,
				maxage,
				"/",
				mw.CookieDomain,
				mw.SecureCookie,
				mw.CookieHTTPOnly,
			)
			c.Redirect(http.StatusMovedPermanently, "/friends/add")
		}
	} else {
		var context = make(gin.H)
		c.ShouldBindJSON(user)
		localizer := Localize.GetLocalizer(c)
		GetSignupPageTranslation(&context, localizer)
		context["GivenName"] = user.GivenName
		context["Surname"] = user.Surname
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
	var users []Models.UserDetail
	if err := Models.GetAllUsers(&users); err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		log.Println(users)
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}

func CreateUser(c *gin.Context) {
	var user Models.User
	c.BindJSON(&user)
	log.Println(user.Password)
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), DefaultCost)
	if hashErr != nil {
		log.Println(hashErr.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	}
	user.Password = string(hash)
	user.Status = 0
	log.Println(user.Password)
	err := Models.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
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
		c.AbortWithStatus(http.StatusBadRequest)
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
	err = Models.UpdateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func DeleteUser(c *gin.Context) {
	var user Models.User
	id := c.Params.ByName("id")
	err := Models.DeleteUser(&user, id)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
	} else {
		c.JSON(http.StatusOK, user)
	}
}
