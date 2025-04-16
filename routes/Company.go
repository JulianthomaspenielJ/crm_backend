package routes

import (
	"bytes"
	"crm/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	logFile        *os.File
	err            error
	contactsLogger *log.Logger
)

func initLoggers() {
	logFile, err = os.OpenFile("contacts.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening contacts log file:", err)
	}
	contactsLogger = log.New(logFile, "CONTACTS: ", log.LstdFlags|log.Lshortfile)

}

func AddCompany(c *gin.Context) {
	initLoggers()
	var data models.Company

	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request while adding a new company %s", time.Now().Format(time.RFC3339), err)
		return
	}

	Logo := "Company.png"
	res, id, err := models.InsertCompany(data.CompanyName, data.Email, data.CompanyDomain, data.Phone, data.IndustryType, data.CompanyDescription, data.City, data.State, data.Country, data.PostalCode, data.NumberOfEmployee, data.AnnualRevenue, data.CompanyGST, data.CompanyPAN, data.MSME, data.NumberOfBranches, data.CIN, data.SocialMedia, data.CreatedBy, data.CompanyOwner, Logo, data.TimeZone, data.OwnerID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error": err,
		})
		contactsLogger.Printf("[%s] Error occurred while executing insert query : %s", time.Now().Format(time.RFC3339), err)
		return
	} else if res == "domain" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "Company Domain Already exists",
			"Response": res,
		})
		contactsLogger.Printf("[%s] Company domain already exists : %s", time.Now().Format(time.RFC3339), res)
		return
	} else if res == "Email" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "Email Already exists",
			"Response": res,
		})
		contactsLogger.Printf("[%s] Company email already exists : %s", time.Now().Format(time.RFC3339), res)
		return
	} else if res == "Phone" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "Mobile number already exists",
			"Response": res,
		})
		return
	} else if res == "PAN" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "PAN number already exists",
			"Response": res,
		})
		contactsLogger.Printf("[%s] Company PAN already exists : %s", time.Now().Format(time.RFC3339), res)
		return
	} else if res == "GST" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "GST number already exists",
			"Response": res,
		})
		contactsLogger.Printf("[%s] Company GST already exists : %s", time.Now().Format(time.RFC3339), res)
		return
	} else {
		contactsLogger.Printf("[%s] Company data stored successfully: %s", time.Now().Format(time.RFC3339), data.CompanyName)
		fmt.Print("It's came here")
		fmt.Print("Contacts data's", data.Contacts)
		for _, contact := range data.Contacts {
			fmt.Printf("Name %s, Email %s, Mobile %s, Lead Status %s", contact.Name, contact.Email, contact.Mobile, contact.LeadStatus)
			_, err := models.InsertContact(contact.Name, contact.Email, contact.Mobile, contact.LifeCycleStage, contact.LeadStatus, "N/A", "N/A", data.CreatedBy, "N/A", 1, id, data.OwnerID, contact.Designation)
			if err != nil {
				fmt.Print(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"Error":   err,
					"Message": "Failed to insert contact",
				})
				contactsLogger.Printf("[%s] Error while inserting contact %s: %s", time.Now().Format(time.RFC3339), contact.Email, err)
				return
			}
			contactsLogger.Printf("[%s] Contact data stored successfully: %s", time.Now().Format(time.RFC3339), contact.Email)
		}

		c.JSON(http.StatusOK, gin.H{
			"Message":  "Data Inserted Successfully",
			"Response": res,
		})
	}
}

func GetCompanyTableData(c *gin.Context) {
	initLoggers()
	var data models.CompanyStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Invalid request data while fetching company data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	res, err := models.GetCompanyTableData(data.ID, data.Role, data.Team)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching company data from DB : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	})
	//contactsLogger.Printf("[%s] Company data fetched successfully ", time.Now().Format(time.RFC3339))
}
func GetCompanyData(c *gin.Context) {
	initLoggers()
	var Data models.Company

	if err := c.ShouldBindJSON(&Data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data while getting company detailed data: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.GetCompanyData(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	mail, foundMail, err := models.FetchMailData(Data.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching mailing data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	note, foundNote, err := models.FetchNotesData(Data.ID, 1)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching notes data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	call, foundCall, err := models.FetchCallsData(Data.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	task, foundTask, err := models.FetchTaskData(Data.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	meet, foundMeet, err := models.FetchMeetingData(Data.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching meeting log data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	deals, foundDeals, err := models.FetchDeals(Data.ID, 1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	contacts, foundContact, err := models.GetContactList(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching contact data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	response := gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	}

	if foundContact {
		response["contacts"] = contacts
	} else {
		response["contacts"] = "No Contacts Found till now!"
	}
	if foundDeals {
		response["deals"] = deals
	} else {
		response["deals"] = "No Deals Found till now!"
	}

	if foundMeet {
		response["meetings"] = meet
	} else {
		response["meetings"] = "No meeting logs Found till now!"
	}

	if foundNote {
		response["notes"] = note
	} else {
		response["notes"] = "No Notes Found till now!"
	}

	if foundTask {
		response["tasks"] = task
	} else {
		response["tasks"] = "No Tasks Found till now!"
	}

	if foundCall {
		response["calls"] = call
	} else {
		response["calls"] = "No calls Found till now!"
	}

	if foundMail {
		response["mail"] = mail
	} else {
		response["mail"] = "No Mail Data Found"
	}

	c.JSON(http.StatusOK, response)
}

func EditCompanyData(c *gin.Context) {
	initLoggers()
	var Data models.Company
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Invalid request data while edit company data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	res, err := models.EditCompanyData(int64(Data.ID), Data.IndustryType, Data.Email, Data.City, Data.State, Data.Country, Data.PostalCode, Data.NumberOfEmployee, Data.NumberOfBranches, Data.AnnualRevenue, Data.TimeZone, Data.CompanyDescription, Data.SocialMedia, Data.Phone, Data.EditedBY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while edit company data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated Successfully",
		"Data":    res,
	})
	contactsLogger.Printf("[%s] Company data updated successfully: %s", time.Now().Format(time.RFC3339), err)
}

type GstNum struct {
	GstNumber string `json:"gstNumber"`
}

func GST(c *gin.Context) {
	var data GstNum
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		return
	}
	fmt.Print(data.GstNumber)

	payload := map[string]string{"gstNumber": data.GstNumber}
	jsonPayload, _ := json.Marshal(payload)

	utilsPath := os.Getenv("UTILS_PATH")
	endpoint := "gst"

	url := fmt.Sprintf("%s%s", utilsPath, endpoint)
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to contact GST service",
			"Error":   err.Error(),
		})

		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.JSON(response.StatusCode, gin.H{
			"message": "GST service returned an error",
		})
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error parsing GST service response",
			"Error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched Successfully",
		"Data":    result,
	})
}

func GetAllCompany(c *gin.Context) {
	initLoggers()
	res, err := models.GetAllCompany()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Error occurred while fetching company data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched Successfully",
		"Data":    res,
	})
	contactsLogger.Printf("[%s] Company data fetched successfully: %s", time.Now().Format(time.RFC3339), err)
}

// func GetCompaniesContactData(c *gin.Context) {
// 	initLoggers()
// 	var Data models.Company
// 	if err := c.ShouldBindJSON(&Data); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"message": "Invalid data format!",
// 			"Error":   err.Error(),
// 		})
// 		contactsLogger.Printf("[%s] Invalid request data while fetching company data : %s", time.Now().Format(time.RFC3339), err)
// 		return
// 	}
// 	res, err := models.GetCompanyData(Data.ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"message": "Something went wrong",
// 			"error":   err.Error(),
// 		})
// 		contactsLogger.Printf("[%s] Error occurred while fetching company data : %s", time.Now().Format(time.RFC3339), err)
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Data fetched Successfully",
// 		"Data":    res,
// 	})
// }
