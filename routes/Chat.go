package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InsertChat(c *gin.Context) {
	var chat models.ChatStruct
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}
	err := models.InsertChat(chat.DealID, 0, chat.Message, "No media", chat.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Chat inserted successfully",
	})
}

func GetChat(c *gin.Context) {
	var data models.ChatStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}
	res, err := models.GetChat(data.DealID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Chat fetched successfully",
		"data":    res,
	})
}