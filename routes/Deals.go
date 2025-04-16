package routes

import (
	"crm/models"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AddDeal(c *gin.Context) {
	var deal models.DealsStruct

	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}

	deals, err := models.AddDeal(deal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err,
		})
		return
	}

	if err := DealsMail(deals); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enqueue mail",
			"error":   err.Error(),
		})
		return
	}
	// fmt.Print(deals)

	c.JSON(http.StatusOK, gin.H{
		"message": "Deal created successfully",
	})
}

func GetApprovalDeal(c *gin.Context) {
	res, err := models.GetApproveDeal()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal data fetched successfully",
		"data":    res,
	})
}

func DealDetailedData(c *gin.Context) {
	var deal models.DealsStruct

	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}
	fmt.Print(deal.ID, deal.NotificationID)
	res, sum, err := models.GetDetailedDeal(deal.ID, deal.NotificationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		fmt.Print(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal data fetched successfully",
		"data":    res,
		"sum":     sum,
	})
}

func ApproveDeal(c *gin.Context) {
	var deal models.DealsStruct

	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})

		return
	}
	status := "Quotation Required"
	mail, err := models.ApproveDeal(deal.ID, deal.DealApproverID, deal.DealApprover, deal.QuotationUpdaterName, deal.QuotationUpdaterID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	if err := DealsMail(mail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enqueue mail",
			"error":   err.Error(),
		})
		return
	}
	fmt.Print(mail)
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal Approved successfully",
	})
}

func GetQuotationDeal(c *gin.Context) {
	var data models.DealsStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})

		return
	}
	res, err := models.GetQuotationDeal(data.ID, data.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal data fetched successfully",
		"data":    res,
	})
}

// UpdateQuotation is a handler function that updates the quotation for a deal.
// It takes in a JSON payload with the following fields:
//   - Quotation: the base64-encoded PDF data for the quotation
//   - ID: the ID of the deal
//   - QuotationUpdaterID: the ID of the user who is updating the quotation
//   - QuotationUpdaterName: the name of the user who is updating the quotation
//
// The function first decodes the base64-encoded PDF data, then creates a file with the PDF data and saves it to the "Assets/Quotation" directory.
// Finally, it updates the quotation information in the database using the provided data.
// If any errors occur during the process, the function returns appropriate error messages in the response.
func UpdateQuotation(c *gin.Context) {
	var data struct {
		Quotations           []string `json:"quotations"`
		ID                   int      `json:"id"`
		QuotationUpdaterID   int      `json:"quotationUpdaterID"`
		QuotationUpdaterName string   `json:"quotationUpdaterName"`
		PricingID            int      `json:"pricingID"`
		TotalPrice           string   `json:"pricingAmount"`
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})
		return
	}

	currentTime := time.Now()
	timestamp := currentTime.Format("Jan02-2006_150405")
	baseDir := filepath.Join("Assets", "Deals", "Quotation", fmt.Sprintf("%d", data.ID), timestamp)

	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create directory",
			"error":   err.Error(),
		})
		fmt.Print("While creating directory", err)
		return
	}

	var storedFiles []string

	for i, quotation := range data.Quotations {
		cleanBase64 := strings.TrimPrefix(quotation, "data:application/pdf;base64,")
		cleanBase64 = strings.ReplaceAll(cleanBase64, " ", "+")

		pdfData, err := base64.StdEncoding.DecodeString(cleanBase64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("Failed to decode Base64 data for file %d", i+1),
				"error":   err.Error(),
			})
			return
		}

		fileName := fmt.Sprintf("%s_%d_%d.pdf", timestamp, data.ID, i+1) // File name with S.No
		filePath := filepath.Join(baseDir, fileName)

		file, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to create file %d", i+1),
				"error":   err.Error(),
			})
			fmt.Print("While creating file", err)
			return
		}
		defer file.Close()

		if _, err := file.Write(pdfData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Failed to write file %d", i+1),
				"error":   err.Error(),
			})
			fmt.Print("While writing file", err)
			return
		}

		storedFiles = append(storedFiles, fileName)
	}
	filename := strings.Join(storedFiles, ",")
	fmt.Print(filename)
	deal, err := models.UpdateQuotation(data.ID, data.PricingID, filename, data.QuotationUpdaterID, data.QuotationUpdaterName, data.TotalPrice)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update quotation in database",
			"error":   err.Error(),
		})
		fmt.Print("While updating quotation in database", err)
		return
	}

	if err := DealsMail(deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to enqueue mail",
			"error":   err.Error(),
		})
		fmt.Print("While sending mail", err)
		return
	}

	fmt.Print(deal)

	c.JSON(http.StatusOK, gin.H{
		"message": "PDF files uploaded successfully",
		"files":   storedFiles,
	})
}

func GetAllDealsData(c *gin.Context) {
	var data models.DealsStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data",
			"error":   err.Error(),
		})

		return
	}
	res, err := models.GetAllDealsData(data.ID, data.Team, data.Role, data.UserID, data.MyDeals)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal data fetched successfully",
		"data":    res,
	})
}

func UpdateDealStatus(c *gin.Context) {
	var deal models.DealsStruct

	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})

		return
	}
	if deal.Reason == "" {
		deal.Reason = "N/A"
	}
	res, err := models.UpdateDealStatus(deal.ID, deal.Status, deal.EditedBy, deal.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	fmt.Print(res)
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal Status Updated Successfully",
	})
}

func DenyDeal(c *gin.Context) {
	var deal models.DenyDeals
	if err := c.ShouldBindJSON(&deal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input data! Please try again",
			"error":   err.Error(),
		})
		return
	}
	err := models.DenyDeal(deal.ID, deal.Reason, deal.DeletedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal Denied Successfully",
	})
}

func GetQuotationUploadedData(c *gin.Context) {
	res, err := models.GetQuotationUploadedData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Deal data fetched successfully",
		"data":    res,
	})
}
