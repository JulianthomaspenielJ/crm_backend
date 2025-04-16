package routes

import (
	"crm/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Products(c *gin.Context) {
	componentLoggers()
	var product models.PricingProduct

	if err := c.ShouldBindJSON(&product); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While storing product data data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.StoreProductData(product, product.AddressTo, product.CreatedBy, product.ReferenceName, int(product.IsDeal.Int64), product.DealID)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Invalid request data While storing product data data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product stored successfully",
		"data":    res,
	})
}

func GetProduct(c *gin.Context) {
	componentLoggers()
	product, err := models.GetProductData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred While fetching product data data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product retrieved Successfully",
		"data":    product,
	})
}

func PostProduct(c *gin.Context) {
	componentLoggers()
	var Data models.GetProductsDetails

	if err := c.ShouldBindJSON(&Data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While storing product list data data %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.GetProductListData(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err,
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred While fetching product list data data %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data":    res,
		"message": "Data retried successfully",
	})
}

func UpdateProductData(c *gin.Context) {
	componentLoggers()
	var data models.PricingProduct
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		return
	}
	res, err := models.UpdatePricing(data.ReferenceID, data, data.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Error occurred while executing query ",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data":    res,
		"message": "Data updated successfully",
	})
}

func RemoveProduct(c *gin.Context) {
	componentLoggers()
	var data models.GetProductsDetails
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		return
	}
	err := models.RemoveProduct(data.ID, data.DeletedBy, data.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Error occurred while executing query ",
		})
		fmt.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data removed successfully",
	})
}

func ProductReference(c *gin.Context) {
    componentLoggers()

    var data models.ProductReferenceStruct

    // Bind JSON data
    if err := c.ShouldBindJSON(&data); err != nil {
        fmt.Println("Error binding JSON:", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   err.Error(),
            "message": "Invalid request data. Please check your input and try again.",
        })
        return
    }

    // Fetch product details by ReferenceName
    res, err := models.ProductReference(data.ReferenceName)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   err.Error(),
            "message": "Error occurred while executing query.",
        })
        return
    }

    // Unmarshal product data
    var productsData models.ProductReferenceStruct
    if err := json.Unmarshal([]byte(res.Products), &productsData); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   err.Error(),
            "message": "Failed to parse product data.",
        })
        return
    }

    // Collect product IDs
    var productIDs []int
    for _, product := range productsData.Products {
        productIDs = append(productIDs, product.ID)
    }

    // Fetch component details for the collected product IDs
    componentDetails, err := models.FetchComponentDetails(productIDs)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   err.Error(),
            "message": "Failed to retrieve component details.",
        })
        return
    }

    // Send the response to the client
    c.JSON(http.StatusOK, gin.H{
        "message":    "Data retrieved successfully",
        "components": componentDetails,
        "info":       productsData,
    })
}
