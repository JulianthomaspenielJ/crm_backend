// routes/update_user.go

package routes

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crm/models"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func InitializeDB(database *sql.DB) {
	db = database
}
func userLoggers() {
	logFile, err = os.OpenFile("users.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening contacts log file:", err)
	}
	contactsLogger = log.New(logFile, "USER: ", log.LstdFlags|log.Lshortfile)

}

func UpdateUser(c *gin.Context) {
	userLoggers()
	var updateData models.UpdateUserData
	if err := c.BindJSON(&updateData); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While activating user : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	rowsAffected, Mail, err := models.UpdateUser(db, updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred while executing query in activating user : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err = models.UpdateAudit("Update", updateData.CreatedBy, updateData.ID, "127.0.0", "User Activated")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		return
	}
	err = DealsMail(Mail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again",
		})
		fmt.Println("Error sending email:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    fmt.Sprintf("User activated successfully. Rows affected: %d", rowsAffected),
		"message": "User data activated Successfully",
	})
	contactsLogger.Printf("[%s] User activated successfully : %s", time.Now().Format(time.RFC3339), updateData.Name)
}

func UpdateUserDetails(c *gin.Context) {
	userLoggers()
	var userDetails models.UserDetails

	if err := c.BindJSON(&userDetails); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data While updating user : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	avatarName, err := saveBase64Image(userDetails.AvatarURL, "avatar", "./Assets/AvatarImages/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Error storing data! Please try again",
		})
		contactsLogger.Printf("[%s] Error storing avatar image While updating user : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	// signatureName, err := saveBase64Image(userDetails.SignatureURL, "signature", "./Assets/SignatureImages/")
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	// 	fmt.Printf("Error saving signature: %v\n", err)
	// 	return
	// }

	//fmt.Println("Avatar saved to:", avatarName)
	//fmt.Println("Signature saved to:", signatureName)

	res, err := models.UpdateUserDetails(userDetails.Name, userDetails.Email, userDetails.Mobile, avatarName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error ":  err,
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred While Updating user : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "User details updated successfully",
		"data":       userDetails,
		"avatarPath": avatarName,
		//"signaturePath": signatureName,
		"Response": res,
	})
	contactsLogger.Printf("[%s] User updated successfully : %s", time.Now().Format(time.RFC3339), userDetails.Name)
}

func saveBase64Image(dataURL string, prefix string, dirPath string) (string, error) {
	splitData := strings.Split(dataURL, ",")
	if len(splitData) != 2 {
		return "", fmt.Errorf("invalid data URL")
	}
	base64Data := splitData[1]

	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %v", err)
	}

	timestamp := time.Now().Unix()
	filename := fmt.Sprintf("%d_%s.png", timestamp, prefix)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		fmt.Printf("Directory does not exist, creating directory: %s\n", dirPath)
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create directory: %v", err)
		}
	}

	filePath := filepath.Join(dirPath, filename)
	err = ioutil.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not saved correctly: %s", filePath)
	}

	return filename, nil
}
