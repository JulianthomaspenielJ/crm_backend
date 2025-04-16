package routes

import (
	"crm/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)


func GetProductsData(c *gin.Context) {
	var productData models.GetProductsDetails
	// fmt.Println(productData)

	if err := c.BindJSON(&productData); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": productData,
		"message": "Data fetched successfully",
	})
}
