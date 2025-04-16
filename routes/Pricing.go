package routes

import (
	"crm/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Pricing(c *gin.Context) {
	componentLoggers()

	var requestData struct {
		ID int `json:"id"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While get component data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	fmt.Printf("Received ID: %d\n", requestData.ID)

	components, err := models.GetComponentsData(requestData.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch component ID",
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred While fetching component data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Data fetched successfully",
		"component": components,
	})
}
