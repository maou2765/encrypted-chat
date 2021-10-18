package Middlewares

import (
	"encrypted-chat/Models"
	"fmt"
	"log"
	"net/http"
	"time"

	ginJWT "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Origin   string `form:"origin" json:"origin"`
}

var identityKey = Models.UserIdentityKey

func AuthMiddleware() (*ginJWT.GinJWTMiddleware, error) {
	SigningAlgorithm := "HS256"
	Key := []byte("s+h+jHcI7bXb+KUAqXgTEQPdK4dAjgTbk/DhBDSb/bgLyVShY1P0x6uuLzubSlIPaY5flbkBD3+PDsr+D8IzXbBd1igU+9yJG3xu7toZMFreNXnZsbChzCpQjD4hbpTyKaewkMVm8U+P2uU0b0Wky0dpwljpFpNAjNDFwz/Env8jMPa0H1xKSqR7IhOCEAyS6zqb5RSKnIly5AvuHDNnMC89oQwr/QDQVrUbjA==")
	CookieName := "jwt"
	CookieDomain := "localhost"
	SecureCookie := false
	CookieHTTPOnly := false
	var CookieMaxAge time.Duration = 3600 * 24 * 7
	authMiddleware, err := ginJWT.New(&ginJWT.GinJWTMiddleware{
		Realm:       "chatroom",
		Key:         Key,
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) ginJWT.MapClaims {
			if v, ok := data.(*Models.User); ok {
				fmt.Println(v.Email)
				return ginJWT.MapClaims{
					identityKey: v.Email,
				}
			}
			return ginJWT.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := ginJWT.ExtractClaims(c)
			log.Println(claims)
			return &Models.User{
				Email: claims[identityKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				log.Println(err)
				return "", ginJWT.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password
			var user Models.User
			err := Models.Login(&user, email)
			log.Println(err)
			if err != nil {
				return "", ginJWT.ErrFailedAuthentication
			}
			if passwordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); passwordErr != nil {
				log.Println(passwordErr)
				return nil, ginJWT.ErrFailedAuthentication
			}
			return &Models.User{
				Email:     user.Email,
				GivenName: user.GivenName,
				Surname:   user.Surname,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println(data)
			if _, ok := data.(*Models.User); ok {
				cookie, _ := c.Cookie("jwt")
				token, _ := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
					return Key, nil
				})
				claims := token.Claims.(jwt.MapClaims)
				newToken := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
				newClaims := newToken.Claims.(jwt.MapClaims)
				for key := range claims {
					newClaims[key] = claims[key]
				}
				expire := time.Now().Add(time.Hour)
				newClaims["exp"] = expire.Unix()
				newClaims["orig_iat"] = time.Now().Unix()
				tokenString, _ := newToken.SignedString(Key)
				expireCookie := time.Now().Add(CookieMaxAge)
				maxage := int(expireCookie.Unix() - time.Now().Unix())
				log.Println(expireCookie)
				c.SetCookie(
					CookieName,
					tokenString,
					maxage,
					"/",
					CookieDomain,
					SecureCookie,
					CookieHTTPOnly,
				)
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			log.Println("Unauthorized")
			c.Redirect(http.StatusMovedPermanently, "/login")
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			var user Models.User
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				log.Println(err)
				return
			}
			email := loginVals.Email
			err := Models.GetUserByEmail(&user, email)
			if err != nil {
				c.AbortWithStatus(http.StatusBadGateway)
				return
			}
			log.Println("code:%s, token:%s, expire:%d", code, token, expire)
			if loginVals.Origin == "app" {
				c.JSON(http.StatusOK, user)
			}
			if user.Status == 0 {
				c.Redirect(http.StatusMovedPermanently, "/friends/add")
			} else {
				c.Redirect(http.StatusMovedPermanently, "/chats")
			}
		},
		// TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		// TokenHeadName:  "Bearer",
		SigningAlgorithm: SigningAlgorithm,
		SendCookie:       true,
		SecureCookie:     SecureCookie,   //non HTTPS dev environments
		CookieHTTPOnly:   CookieHTTPOnly, // JS can't modify
		CookieDomain:     CookieDomain,
		CookieName:       CookieName, // default jwt
		CookieMaxAge:     CookieMaxAge,
		// TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteDefaultMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, err
}
