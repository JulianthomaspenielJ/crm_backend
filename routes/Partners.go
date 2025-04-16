package routes

import (
	"crm/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddPartner(c *gin.Context) {
	var data models.PartnerStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid Data! Please Try again",
		})
		return
	}
	// fmt.Printf("Partner Name %s, AuthCode %s, AreaCovered %s, CreatedBy %s, respectiveSalesPerson %s, dueDate %s, registeredDate %s And Contacts are %s",
	// 	data.PartnerName, data.AuthCode, data.AreaCovered, data.CreatedBy, data.RespectiveSalesPerson, data.DueDate, data.RegisteredDate, data.Contacts)
	res, err := models.AddPartner(data.PartnerName, data.AuthCode, data.AreaCovered, data.CreatedBy, data.Address, data.RespectiveSalesPerson, data.DueDate, data.RegisteredDate, data.Email, data.Phone, data.Contacts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to add partner",
		})
		return
	}
	if res == "exist" {
		c.JSON(http.StatusOK, gin.H{
			"message": "Partner Name or Authorization code Already exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data Inserted Successfully",
	})
}

func GetPartners(c *gin.Context) {
	res, err := models.GetPartnersData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get partners",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"partners": res,
	})
}

func GetPartnersList(c *gin.Context) {
	res, err := models.GetPartnerListData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get partners",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"partners": res,
	})
}

func AddContactPerson(c *gin.Context) {
	var data models.ContactPerson
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	res, err := models.AddContactPerson(data.PartnerID, data.ContactName, data.ContactEmail, data.ContactPhone, data.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to add contact person",
		})
		fmt.Print(err)
		return
	}
	if res == "exist" {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Contact person already exist",
			"message": "Contact person already exist",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":        "Contact added successfully",
		"contact_person": res,
	})

}

func GetDataContactPerson(c *gin.Context) {
	res, err := models.GetDataContactPerson()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get partners",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"contacts": res,
	})
}

func AddGem(c *gin.Context) {
	var data models.GemStructure
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print("bind", err)
		return
	}

	today := time.Now().Format("Jan02-2006_150405")
	dealerIDStr := strconv.Itoa(data.DealerID)
	dirPath := fmt.Sprintf("Assets/Gem/GemCertificates/%s/%s", dealerIDStr, today)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create directory",
		})
		fmt.Print("mkdir", err)
		return
	}

	var fileNames []string

	for i, fileStr := range data.GemCertificate {
		pdfData, err := base64.StdEncoding.DecodeString(fileStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid base64 data",
			})
			fmt.Print("decode", err)
			return
		}

		fileName := fmt.Sprintf("%s_%s_%d.pdf", today, dealerIDStr, i+1)
		filePath := fmt.Sprintf("%s/%s", dirPath, fileName)

		if err := os.WriteFile(filePath, pdfData, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to save PDF file",
			})
			fmt.Print("write", err)
			return
		}

		fileNames = append(fileNames, fileName)
	}
	Status := "Gem Requested"
	gemCertificateStr := strings.Join(fileNames, ",")
	Link := "not now"
	res, err := models.AddGem(data.DealerID, data.OwnerID, data.GemNumber, data.DueDate, data.Quantity, data.DealerName, data.Customer, data.Description, data.GemRaisedBy, Status, gemCertificateStr, Link, data.Team, data.GemUrl, data.GemDate, data.ProductType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to add gem",
		})
		fmt.Print("error occurred here", err)
		return
	}

	dealsData := models.DealsStruct{
		To:      res.To,
		From:    res.From,
		Subject: res.Subject,
		Body:    res.Body,
	}

	if err := DealsMail(dealsData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enqueue mail",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gem added successfully",
	})
}

func GetGemWaitingForApproval(c *gin.Context) {
	data, err := models.GetGemWaitingForApproval()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get gem waiting for approval",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Gem waiting for approval retrieved successfully",
		"data":    data,
	})
}

func GetGem(c *gin.Context) {
	var request struct {
		ID int `json:"id"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	data, err := models.GetGem(request.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to retrieve gem data",
		})
		fmt.Print(err)
		return
	}
	if len(data.GemCertificate) > 0 {
		if len(data.GemCertificate) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "No gem certificates found",
				"message": "Gem certificates data is empty",
			})
			return
		}
		firstFileName := data.GemCertificate[0]
		parts := strings.Split(firstFileName, "_")
		if len(parts) < 3 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invalid file name format",
				"message": "The file name does not follow the expected format",
			})
			return
		}
		timestamp := parts[0] + "_" + parts[1]
		dealerID := parts[2]
		dirPath := fmt.Sprintf("Assets/Gem/GemCertificates/%s/%s", dealerID, timestamp)
		var base64Certificates []string
		for _, fileName := range data.GemCertificate {
			filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
			pdfData, err := os.ReadFile(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": fmt.Sprintf("Failed to read file: %s", fileName),
				})
				fmt.Print(err)
				return
			}
			base64Str := base64.StdEncoding.EncodeToString(pdfData)
			base64Certificates = append(base64Certificates, base64Str)
		}
		data.GemCertificate = base64Certificates
	}
	if len(data.Quotation) > 0 {
		firstFileName := data.Quotation[0]
		parts := strings.Split(firstFileName, "_")
		if len(parts) < 3 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Invalid file name format",
				"message": "The file name does not follow the expected format",
			})
			return
		}
		timestamp := parts[0] + "_" + parts[1]
		dealerID := parts[2]
		dirPath := fmt.Sprintf("Assets/Gem/Pricing/%s/%s", dealerID, timestamp)
		var base64Quotations []string
		for _, fileName := range data.Quotation {
			filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
			pdfData, err := os.ReadFile(filePath)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   err.Error(),
					"message": fmt.Sprintf("Failed to read file: %s", fileName),
				})
				fmt.Print(err)
				return
			}
			base64Str := base64.StdEncoding.EncodeToString(pdfData)
			base64Quotations = append(base64Quotations, base64Str)
		}
		data.Quotation = base64Quotations
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Gem data retrieved successfully",
		"data":    data,
	})
}

func GetAwaitingPricingGEMList(c *gin.Context) {
	res, err := models.GetAwaitingPricingGEMList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get GEM list",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "GEM list retrieved successfully",
		"data":    res,
	})
}

func GetUpdatedPricingGEMList(c *gin.Context) {
	res, err := models.GetUpdatedPricingGEMList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get GEM list",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "GEM list retrieved successfully",
		"data":    res,
	})
}

func UploadQuotationGem(c *gin.Context) {
	fmt.Print("Gem function started")
	var data models.GemStructure
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	dealerIDStr := strconv.Itoa(data.ID)
	fmt.Print("The gem ID ", dealerIDStr, " ")
	today := time.Now().Format("Jan02-2006_150405")
	dirPath := fmt.Sprintf("Assets/Gem/Pricing/%s/%s", dealerIDStr, today)

	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to create directory",
		})
		fmt.Print("error while creating dir ", err)
		return
	}
	fmt.Print("DIR PATH , ", dirPath, " ")
	var fileNames []string
	for i, fileStr := range data.Quotation {

		if strings.HasPrefix(fileStr, "data:application/pdf;base64,") {
			fileStr = strings.TrimPrefix(fileStr, "data:application/pdf;base64,")
		}

		pdfData, err := base64.StdEncoding.DecodeString(fileStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   err.Error(),
				"message": "Invalid data! Please try again",
			})
			fmt.Print("invalid data", err)
			return
		}
		fileName := fmt.Sprintf("%s_%s_%d.pdf", today, dealerIDStr, i+1)
		filePath := fmt.Sprintf("%s/%s", dirPath, fileName)
		if err := os.WriteFile(filePath, pdfData, os.ModePerm); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to save PDF file",
			})
			fmt.Print("while storing file", err)
			return
		}
		fileNames = append(fileNames, fileName)
		quotationStr := strings.Join(fileNames, ",")
		fmt.Print("File names, ", quotationStr, " ")
		Status := "Quotation Updated"
		err = models.UploadQuotationGem(data.ID, quotationStr, data.UpdatedBy, Status, data.GEMFinalURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   err.Error(),
				"message": "Failed to add gem",
			})
			fmt.Print(err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Quotation uploaded successfully",
	})
}

func ApproveGeM(c *gin.Context) {
	var Data models.GeMApprove
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	err := models.ApproveGeM(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to approve GEM",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "GEM approved successfully",
	})
}

func DenyGeM(c *gin.Context) {
	var Data models.GeMDeny
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	err := models.DenyGeM(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to approve GEM",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "GEM approved successfully",
	})
}

func GetMyGem(c *gin.Context) {
	var Data models.AllGeM
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		fmt.Print(err)
		return
	}
	fmt.Print(Data)
	res, err := models.GetMyGem(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to approve GEM",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "GEM approved successfully",
		"data":    res,
	})
}

func UpdateStatus(c *gin.Context) {
	var Data models.UpdateStatus
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		return
	}
	err := models.UpdateGemStatus(Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to update status",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Status updated successfully",
	})

}

func GetPartner(c * gin.Context){
	var Data models.PartnerStruct
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request",
		})
		return
	}
	res, err := models.GetPartner(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Failed to get partner",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Partner fetched successfully",
		"data":    res,
	})
}	