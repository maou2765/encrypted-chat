package Controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChatIndex(c *gin.Context) {
	context := gin.H{}
	c.HTML(http.StatusOK, "chat", context)
}
