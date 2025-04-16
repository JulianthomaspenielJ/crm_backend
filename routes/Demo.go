package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddDemo(c *gin.Context) {
	var Demo models.Demo

	if err := c.ShouldBindJSON(&Demo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}

	err := models.AddDemo(Demo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Demo created successfully",
	})
}

func GetDemo(c *gin.Context) {
	var Demo models.Demo
	if err := c.ShouldBindJSON(&Demo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}

	res, err := models.GetDemo(Demo.Team, Demo.Role, Demo.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Demo data fetched successfully",
		"data":    res,
	})
}

func GetDetailedDemo(c *gin.Context) {
	var Demo models.Demo

	if err := c.ShouldBindJSON(&Demo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}

	data, err := models.GetDemoData(Demo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}

	// Fetch the summary data
	summary, err := models.GetDemoSummary(Demo.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong while fetching summary!",
			"error":   err.Error(),
		})
		return
	}

	// Send both demo and summary data together
	c.JSON(http.StatusOK, gin.H{
		"message": "Demo data fetched successfully",
		"data":    data,
		"summary": summary,
	})
}


func UpdateDemo(c * gin.Context){
	var Demo models.Demo
	if err := c.ShouldBindJSON(&Demo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}
	err := models.UpdateDemo(Demo.ID, Demo.Status, Demo.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Demo updated successfully",
	})
}