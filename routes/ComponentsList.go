package routes

import (
	"crm/models"
	"fmt"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Component struct {
	ID            int    `json:"id"`
	ComponentName string `json:"componentName"`
	Description   string `json:"description"`
	ProductType   string `json:"productType"`
	CreatedBy     string `json:"createdBy"`
}

func ComponentList(c *gin.Context) {
	componentLoggers()
	var components []Component

	if err := c.ShouldBindJSON(&components); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While inserting new component list : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	for _, comp := range components {
		err := models.InsertComponent(comp.ComponentName, comp.Description, comp.ProductType, comp.CreatedBy)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to insert component",
				"message": "Something went wrong! Please try again",
			})
			contactsLogger.Printf("[%s] Failed to inserting a new component list : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"ComponentList": components,
		"message":       "Data stored successfully.",
	})
	contactsLogger.Printf("[%s] New component list data created successfully", time.Now().Format(time.RFC3339))
}

func GetComponents(c *gin.Context) {
	componentLoggers()
	data, err := models.GetComponentslist()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to fetch data", "Errors are :": err})
		contactsLogger.Printf("[%s] Error occurred while get component data list %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"DATA : ": data,
		"message": "Data fetched successfully",
	})
}

func UpdateComponentList(c *gin.Context) {
	var Data Component
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input format! Please try again.",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while updating component data list %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.UpdateComponentList(Data.ID, Data.ComponentName, Data.Description, Data.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while updating component data list %s", time.Now().Format(time.RFC3339), err)
		return
	}
	fmt.Print(Data)
	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
	})
}

func DeleteComponentList(c *gin.Context) {
	var Data Component
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input format! Please try again.",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while deleting component data list %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.DeleteComponentList(Data.ID, Data.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while deleting component data list %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Deleted successfully",
	})
}
