package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CallPythonServer(c *gin.Context) {
	var users models.UserDetails

	if c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data!",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		//ds"data": res,
	})
}
