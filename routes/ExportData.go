package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ExportData(c *gin.Context) {
	var Data models.ExportStruct
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"err":     err.Error(),
		})
		return
	}

	res, err := models.ExportData(Data.FromDate, Data.ToDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	})
}
