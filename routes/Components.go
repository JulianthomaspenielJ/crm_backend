package routes

import (
	"crm/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type DeleteRequest struct {
	ID          int    `json:"id"`
	UpdatedBy   string `json:"deletedBy"`
	RecoveredBy string `json:"recoveredBy"`
}

func componentLoggers() {
	logFile, err = os.OpenFile("pricing.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening contacts log file:", err)
	}
	contactsLogger = log.New(logFile, "PRICING: ", log.LstdFlags|log.Lshortfile)

}
func DeleteComponent(c *gin.Context) {
	componentLoggers()
	var request DeleteRequest
	var err error

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data while delete component data : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err = models.DeleteComponent(request.ID, request.UpdatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error :": err,
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred while delete component data in the db : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data":    "Updated Successfully",
		"message": "Data deleted Successfully",
	})
	contactsLogger.Printf("[%s] Component data deleted successfully", time.Now().Format(time.RFC3339))
}

func FetchComponents(c *gin.Context) {
	componentLoggers()
	components, err := models.GetComponents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch components"})
		contactsLogger.Printf("[%s] Error occurred in query while component data fetching : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"components": components,
		"message":    "data fetched",
	})
}

func PostComponentsData(c *gin.Context) {
	componentLoggers()
	var variantData models.VariantData

	if err := c.BindJSON(&variantData); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data while create a new component data : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.InsertComponents(variantData.VariantName, variantData.VariantType, variantData.VariantPrice, variantData.VariantDescription, variantData.Brand, variantData.VariantID, variantData.CreatedBy)
	if err != nil {
		fmt.Print("Error occurred ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to insert component data into database",
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred while query executing in create a new component data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	if res == "Component Name already exist" {
		c.JSON(http.StatusOK, gin.H{
			"message": "name",
		})
		contactsLogger.Printf("[%s] Error Component name already exist : %s", time.Now().Format(time.RFC3339), variantData.VariantName)
		return
	}
	fmt.Print(res)

	c.JSON(http.StatusOK, gin.H{
		"message": "Data saved successfully",
		"data":    variantData,
	})
	contactsLogger.Printf("[%s] Component created successfully : %s", time.Now().Format(time.RFC3339), variantData.VariantName)
}

func UpdateComponent(c *gin.Context) {
	componentLoggers()
	var component models.UpdateComponent

	if err := c.ShouldBindJSON(&component); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data while update component : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.UpdateComponents(component.ID, component.VarientName, component.VarientDescription, component.VariantPrice, component.Brand, component.UpdatedBy)
	if err != nil {
		fmt.Print("error : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Error: ": err.Error()})
		contactsLogger.Printf("[%s] Error occurred while query executing in update component : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    component,
		"message": "Data Updated Successfully",
	})
	contactsLogger.Printf("[%s] Component updated successfully: %s", time.Now().Format(time.RFC3339), component.VarientName)
}

func GetDeletedData(c *gin.Context) {
	componentLoggers()
	res, err := models.GetDeletedData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err,
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred while query executing in fetching deleted component list : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data":    res,
		"Message": "Data fetched successfully",
	})
}
func RecoverComponent(c *gin.Context) {
	componentLoggers()
	var component DeleteRequest

	if err := c.ShouldBindJSON(&component); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request while recovering component: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.RecoverComponent(component.ID, component.RecoveredBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"Error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while query executing in recover deleted component : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Restored successfully",
		"data":    "Data updated",
	})
	contactsLogger.Printf("[%s] Component recovered successfully, By : %s", time.Now().Format(time.RFC3339), component.RecoveredBy)
}
