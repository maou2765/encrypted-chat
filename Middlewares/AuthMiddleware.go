package Middlewares

import (
	"encrypted-chat/Models"
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `from:"email" json:"email" binding:"required"`
	Password string `from:"password" json:"password" binding:"required"`
}

var identityKey = Models.UserIdentityKey

func AuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test",
		Key:         []byte("secret"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*Models.User); ok {
				return jwt.MapClaims{
					identityKey: v.Email,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			log.Println(claims)
			return &Models.User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password
			var user Models.User
			err := Models.Login(&user, email)
			if err != nil {
				return "", jwt.ErrFailedAuthentication
			}

			if passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); passwordErr != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			return &Models.User{
				Email:     user.Email,
				GivenName: user.GivenName,
				SurnName:  user.SurnName,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*Models.User); ok && v.Email == "user1@email.com" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Println("Unauthorized")
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}
