package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CodeRequest struct {
	Code string `json:"code"`
}

func MeetGoogle(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"meeting_url": "Not yet updated",
	})
}
