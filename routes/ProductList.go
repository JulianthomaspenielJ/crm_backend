package routes

import (
	"crm/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetProductList(c *gin.Context) {
	componentLoggers()
	data, err := models.GetProductLists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error : ": err,
			"message":  "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Failed to get data from product list : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data ":   data,
		"message": "Data fetched successfully",
	})
}

func GetProductMainList(c *gin.Context) {
	componentLoggers()
	data, err := models.GetProductMainList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Failed to get data from product main list : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    data,
		"Message": "Data Get successfully",
	})
}

func ProductID(c *gin.Context) {
	componentLoggers()
	var products models.ProductID
	if err := c.ShouldBindJSON(&products); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While get data from product list : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.ProductIDData(products.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred While get data from product list : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Res":     res,
		"message": "Data fetched successfully",
	})
}

func CreateProductMainList(c *gin.Context) {
	componentLoggers()
	var product models.PricingMainListProduct

	if err := c.ShouldBindJSON(&product); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While inserting new data to product main list : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	insertedData, err := models.InsertProductMainList(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert product"})
		contactsLogger.Printf("[%s] Error occurred While inserting a new data to product main list : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product Stored Successfully",
		"data":    insertedData,
	})
}

func UserList(c *gin.Context) {
	res, err := models.UserList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Failed to get data from product main list : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    res,
		"Message": "Data Get successfully",
	})
}
