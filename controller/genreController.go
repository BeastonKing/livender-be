package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateGenre(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
