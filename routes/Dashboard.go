package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardStruct struct {
	Month  int    `json:"month"`
	Year   string `json:"year"`
	ID     int    `json:"userID"`
	Status string `json:"status"`
}

func Dashboard(c *gin.Context) {
	var data DashboardStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data request",
			"Error":   err.Error(),
		})
		return
	}
	prog, foundProg, err := models.Progress(data.Year, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! ",
			"Error":   err.Error(),
		})
		return
	}

	contact, foundContact, err := models.ContactsCount(data.Year, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! ",
			"Error":   err.Error(),
		})
		return
	}

	company, foundCompany, err := models.CompanyCount(data.Year, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! ",
			"Error":   err.Error(),
		})
		return
	}

	dealCounts, foundDeals, err := models.DealCounts(data.Year, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"Error":   err.Error(),
		})
		return
	}

	response := gin.H{
		"message": "Data fetched successfully",
	}

	quotation, foundQuotation, err := models.QuotationCount(data.Year, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! ",
			"Error":   err.Error(),
		})
		return
	}

	DealStatus, foundStatus, err := models.RecentDeals(data.Status, data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! ",
			"Error":   err.Error(),
		})
		return
	}
	response = gin.H{
		"message": "Data fetched successfully",
	}

	if foundProg {
		response["progress"] = prog
	} else {
		response["progress"] = "No progress found"
	}

	if foundCompany {
		response["companyCount"] = company
	} else {
		response["companyCount"] = "No Contacts Found"
	}

	if foundContact {
		response["contactCount"] = contact
	} else {
		response["contactCount"] = "No Contacts Found"
	}

	if foundQuotation {
		response["quotation"] = quotation
	} else {
		response["quotation"] = "No Contacts Found"
	}

	if foundDeals {
		response["dealCounts"] = dealCounts
	} else {
		response["dealCounts"] = "No Deals Found"
	}

	if foundStatus {
		response["dealStatus"] = DealStatus
	} else {
		response["dealStatus"] = "No Deals Found"
	}

	c.JSON(http.StatusOK, response)
}
