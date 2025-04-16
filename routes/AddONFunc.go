package routes

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"crm/models"

	"github.com/gin-gonic/gin"
)

func GetComplianceData(c *gin.Context) {
	type ComplianceData struct {
		Link string `json:"link"`
	}
	var data ComplianceData
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}

	payload := map[string]string{"link": data.Link}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error marshaling payload",
			"error":   err.Error(),
		})
		return
	}

	utilsPath := os.Getenv("UTILS_PATH")
	endpoint := "compliance"

	url := fmt.Sprintf("%s%s", utilsPath, endpoint)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to contact Compliance data link service",
			"error":   err.Error(),
		})
		return
	}
	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Println("Error closing response body:", err)
		}
	}()

	if res.StatusCode != http.StatusOK {
		c.JSON(res.StatusCode, gin.H{
			"message": "Compliance web service returned an error",
			"error":   res.Status,
		})
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error reading response body",
			"error":   err.Error(),
		})
		return
	}
	//log.Println("Raw Response Body:", string(body))

	var result [][]string
	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing compliance data service response",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Data":    result,
	})
}

func AddDataSheet(c *gin.Context) {
	var Data models.DataSheetData
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}
	err := models.AddDataSheet(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error adding data",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data added successfully",
	})
}

func GetDataSheetName(c *gin.Context) {
	res, err := models.GetDataSheetName()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Data":    res,
	})
}

func GetDataSheetDatas(c *gin.Context) {
	var Datas models.ProcessorChipsetData
	if err := c.ShouldBindJSON(&Datas); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}

	chip, err := models.GetDataSheetByChipSet(Datas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"error":   err.Error(),
		})
		fmt.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"ChipSet": chip,
		// "Processor": processor,
	})
}

// func AddProcessors(C *gin.Context) {
// 	var Data models.ProcessorData
// 	if err := C.ShouldBindJSON(&Data); err != nil {
// 		C.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Invalid data",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	err := models.AddProcessors(Data)
// 	if err != nil {
// 		C.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Error adding data",
// 			"error":   err.Error(),
// 		})
// 		return
// 	}
// 	C.JSON(http.StatusOK, gin.H{
// 		"message": "Data added successfully",
// 	})
// }

func GetProcessorData(c *gin.Context) {
	var Data models.ProcessorStruct
	err := c.ShouldBindJSON(&Data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"err":     err,
		})
		return
	}
	res, err := models.GetProcessorData(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"err":     err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	})
}

func getImageBase64(imagePath string) (string, error) {
	imageFile, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}

	defer imageFile.Close()

	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return "", err
	}

	base64Image := base64.StdEncoding.EncodeToString(imageData)
	return base64Image, nil
}

func GetDataSheetData(c *gin.Context) {
	var Data models.DataSheetName
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}

	res, err := models.GetDataSheetData(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"error":   err.Error(),
		})
		return
	}

	// baseImagePath := "Assets/Certificates/Images/"
	// var imagePath string

	// switch Data.Product {
	// case "Laptop":
	// 	imagePath = filepath.Join(baseImagePath, "ExecutiveLaptop.png")
	// case "All in one pc":
	// 	imagePath = filepath.Join(baseImagePath, "AIO.png")
	// case "Tower Desktop":
	// 	imagePath = filepath.Join(baseImagePath, "TowerDesktop.png")
	// case "Server":
	// 	imagePath = filepath.Join(baseImagePath, "Workstation.png")
	// default:
	// 	imagePath = filepath.Join(baseImagePath, "TowerDesktop.png")
	// }

	// var base64Image string
	// if imagePath != "" {
	// 	base64Image, err = getImageBase64(imagePath)
	// 	if err != nil {
	// 		c.JSON(http.StatusInternalServerError, gin.H{
	// 			"message": "Error reading image file",
	// 			"error":   err.Error(),
	// 		})
	// 		return
	// 	}
	// } else {
	// 	base64Image = ""
	// }

	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Data":    res,
		// "Image":   base64Image,
	})
}

func GetNotificationsNumber(c *gin.Context) {
	var Data models.Notifications
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}
	res, err := models.GetNotificationsNumber(Data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"error":   err.Error(),
		})
		return
	}
	gem, err := models.GetGemNotificationNumber(Data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error fetching data",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Data":    res,
		"gem":     gem,
	})
}

func CheckAndNotifyDueTasks() {
	fmt.Print("came here the next step process")
	data, err := models.CheckAndNotifyDueTasks()
	if err != nil {
		log.Println("Error checking and notifying due tasks:", err)
		return
	}
	fmt.Print(data)
}

func AddProcessor(c *gin.Context) {
	var data models.ProcessorStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Data! Please Try again",
		})
		return
	}
	res, err := models.AddProcessor(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to add processor",
		})
		return
	}
	if res == "exist" {
		c.JSON(http.StatusConflict, gin.H{
			"message": "Processor Already exist",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Inserted Successfully",
	})
}

func GetProcessor(c *gin.Context) {
	res, err := models.GetProcessor()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get processor",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Data":    res,
	})
}

func UpdateGemNotification(c *gin.Context) {
	var data models.Notifications
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Data! Please Try again",
		})
		return
	}
	err := models.UpdateGemNotification(data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update gem notification",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Updated Successfully",
	})

}

func RespectiveSalesPerson(c *gin.Context) {

	res, err := models.RespectiveSalesPerson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update gem notification",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Updated Successfully",
		"data":    res,
	})
}

func GetProductImage(c *gin.Context) {
	var data models.ProcessorChipsetData
	err := c.ShouldBindJSON(&data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messasge": "Bad Request",
			"err":      err,
		})
		return
	}

	baseImagePath := "Assets/Certificates/Images/"
	var imagePath string

	switch data.Product {
	case "Laptop":
		imagePath = filepath.Join(baseImagePath, "ExecutiveLaptop.png")
	case "All in one pc":
		imagePath = filepath.Join(baseImagePath, "AIO.png")
	case "Tower Desktop":
		imagePath = filepath.Join(baseImagePath, "TowerDesktop.png")
	case "WorkStation":
		imagePath = filepath.Join(baseImagePath, "Workstation.png")
	case "Server":
		imagePath = filepath.Join(baseImagePath, "Workstation.png")
	case "Gaming Laptop":
		imagePath = filepath.Join(baseImagePath, "GamingLaptop.png")
	case "Mini PC":
		imagePath = filepath.Join(baseImagePath, "MiniPC.png")
	case "Tablet":
		imagePath = filepath.Join(baseImagePath, "Tablet.png")
	default:
		imagePath = filepath.Join(baseImagePath, "TowerDesktop.png")
	}

	var base64Image string
	if imagePath != "" {
		base64Image, err = getImageBase64(imagePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error reading image file",
				"error":   err.Error(),
			})
			return
		}
	} else {
		base64Image = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"Image":   base64Image,
	})
}

func MarkAsRead(c *gin.Context) {
	var data models.Notifications
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err,
			"message": "Invalid Data! Please Try again",
		})
		return
	}
	err := models.MarkAsRead(data.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err,
			"message": "Failed to update gem notification",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Updated Successfully",
	})
}
