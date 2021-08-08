package Routes

import (
	"io"
	"log"
	"maou2765/encrypted-chat/Controllers"
	"maou2765/encrypted-chat/Middlewares"
	"maou2765/encrypted-chat/Models"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func welcomeHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(Models.UserIdentityKey)
	c.JSON(200, gin.H{
		"userID":   claims[Models.UserIdentityKey],
		"userName": user.(*Models.User).Email,
		"text":     "Welcome",
	})
}

func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authMiddleware, err := Middlewares.AuthMiddleware()

	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", welcomeHandler)
	}

	grp1 := r.Group("/user-api")
	{
		grp1.GET("user", Controllers.GetUsers)
		grp1.POST("user", Controllers.CreateUser)
		grp1.GET("user/:id", Controllers.GetUserByID)
		grp1.PUT("user/:id", Controllers.UpdateUser)
		grp1.DELETE("user/:id", Controllers.DeleteUser)
	}

	return r
}
