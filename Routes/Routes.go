package Routes

import (
	"encrypted-chat/Controllers"
	"encrypted-chat/Middlewares"
	"encrypted-chat/Models"
	"io"
	"log"
	"net/http"
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-contrib/static"
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
func createHTMLRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("login", "templates/layouts/base.html", "templates/login/index.html")
	r.AddFromFiles("signup", "templates/layouts/base.html", "templates/signup/index.html")
	r.AddFromFiles("add_friend", "templates/layouts/base.html", "templates/add_friend/index.html")
	r.AddFromFiles("chat", "templates/layouts/base.html", "templates/chat/index.html")
	return r
}
func SetupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.Use(gin.Recovery())
	// r.StaticFS("/static", http.Dir("/static"))
	r.Use(static.Serve("/static", static.LocalFile("./static", false)))
	r.HTMLRender = createHTMLRender()

	gin.DefaultWriter = io.MultiWriter(os.Stdout)

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	authMiddleware, err := Middlewares.AuthMiddleware()

	if err != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}

	r.GET("/login", Controllers.LoginIndex)
	r.POST("/login", authMiddleware.LoginHandler)
	r.GET("/signup", Controllers.SignupIndex)
	r.POST("/signup", Controllers.Signup)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := r.Group("/")
	auth.GET("/refresh-token", authMiddleware.RefreshHandler)
	// Refresh time can be longer than token timeout
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", welcomeHandler)
		userGp := auth.Group("/user")
		{
			userGp.GET("", Controllers.GetUsers)
			userGp.POST("", Controllers.CreateUser)
			userGp.GET("/:id", Controllers.GetUserByID)
			userGp.PUT("/:id", Controllers.UpdateUser)
			userGp.DELETE("/:id", Controllers.DeleteUser)
		}
		fdGp := auth.Group("/friends")
		{
			fdGp.GET("", Controllers.SearchFriends)
			fdGp.POST("", Controllers.AddFriends)
			fdGp.GET("/add", Controllers.AddFriendIndex)
		}
		chatGp := auth.Group("/chats")
		{
			chatGp.GET("", Controllers.ChatIndex)
		}
	}

	return r
}
