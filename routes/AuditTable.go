package routes

import (
	"bytes"
	"crm/models"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func auditLoggers() {
	logFile, err = os.OpenFile("audit_table.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening contacts log file:", err)
	}
	contactsLogger = log.New(logFile, "AUDIT_TABLE: ", log.LstdFlags|log.Lshortfile)

}
func AuditTable(c *gin.Context) {
	auditLoggers()
	var data models.AuditStructure

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err.Error(),
		})
		return
	}
	//fmt.Println("Request Body: ", string(bodyBytes))

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	auditFunctions := map[string]func() (interface{}, error){
		"Users": func() (interface{}, error) {
			return models.UserAudit()
		},
		"Components": func() (interface{}, error) {
			return models.ComponentAudit()
		},
		"Components list": func() (interface{}, error) {
			return models.ComponentListAudit()
		},
		"Products": func() (interface{}, error) {
			return models.ProductsAudit()
		},
		"Products Main list": func() (interface{}, error) {
			return models.ProductsMainListAudit()
		},
		"Products list": func() (interface{}, error) {
			return models.ProductsListAudit()
		},
		"Companies": func() (interface{}, error) {
			return models.CompaniesAudit()
		},
		"Contacts": func() (interface{}, error) {
			return models.ContactAudit()
		},
		"Deals": func() (interface{}, error) {
			return models.DealsAudit()
		},
		"Margin": func() (interface{}, error) {
			return models.MarginsAudit()
		},
		"Partners": func() (interface{}, error) {
			return models.PartnersAudit()
		},
		"Partner Contacts": func() (interface{}, error) {
			return models.PartnerContctAudit()
		},
		"GeM": func() (interface{}, error) {
			return models.GeMAudit()
		},
	}

	if auditFunc, exists := auditFunctions[data.Audit]; exists {
		res, err := auditFunc()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": err.Error(),
			})
			contactsLogger.Printf("[%s] Error occurred while fetching %s audit table data: %s", time.Now().Format(time.RFC3339), data.Audit, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"Data": res,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid audit type",
			"message": "The provided audit type is not supported.",
		})
	}
}

func AuditData(c *gin.Context) {
	auditLoggers()
	var data models.AuditDetailsData

	bodyBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
		return
	}
	//fmt.Println("Request Body: ", string(bodyBytes))

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	if err := c.ShouldBindJSON(&data); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err,
		})
		contactsLogger.Printf("[%s] Invalid request data while get audit table data ID : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	message := "Neither 'Products_data' nor 'User' found in the description"
	var res interface{}

	if strings.Contains(data.Description, "Products_data") {
		message = "'Products_data' exists in the description"
		res, err = models.GetProductDataQuery(data.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch products audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "User") {
		res, err = models.GetUserDataQuery(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch users audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
		message = "'users' exists in the description"
	} else if strings.Contains(data.Description, "components_list") {
		message = "'components_list' exists in the description"
		res, err = models.AuditComponentsList(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch component list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "components ") {
		message = "'Components' exists in the description"
		res, err = models.AuditComponentsData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch component audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "Components ") {
		message = "'Components' exists in the description"
		res, err = models.AuditComponentsData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch component audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "product_main_list ") {
		message = "'Components' exists in the description"
		res, err = models.AuditProductMainData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product main list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "product list ") {
		message = "'Components' exists in the description"
		res, err = models.AuditProductData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "Company") {
		message = "'Company' exists in the description"
		res, err = models.CompaniesAuditData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "Deal") {
		message = "'Deal' exists in the description"
		res, err = models.DealAuditData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "Deals") {
		message = "'Deals' exists in the description"
		res, err = models.DealAuditData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "contact") {
		message = "'contact' exists in the description"
		res, err = models.ContactAuditData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	} else if strings.Contains(data.Description, "Contact") {
		message = "'Contact' exists in the description"
		res, err = models.ContactAuditData(data.ID)
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"ERR": err,
			})
			contactsLogger.Printf("[%s] Error occurred while fetch product list audit details data : %s", time.Now().Format(time.RFC3339), err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":    res,
		"message": message,
	})
	//contactsLogger.Printf("[%s] Audit table detailed data fetched successfully", time.Now().Format(time.RFC3339))
}
