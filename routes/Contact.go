package routes

import (
	"crm/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AddContact(c *gin.Context) {
	initLoggers()
	var data models.Contacts
	if err := c.ShouldBindJSON(&data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid Request while inserting a new contact data %s", time.Now().Format(time.RFC3339), err)
		return
	}
	fmt.Print(data.Company)
	res, err := models.InsertContact(data.Name, data.Email, data.Mobile, data.LifeCycleStage, data.LeadStatus, data.LegalBasis, data.ContactLevelType, data.CreatedBy, data.JobRole, data.IsCompany, data.Company, data.OwnerID, data.LeadSource)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err,
			"message": "Something went wrong! Please try again",
		})
		contactsLogger.Printf("[%s] Error occurred while inserting a new contact data %s", time.Now().Format(time.RFC3339), err)
		return
	}
	if res == "Email" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "User Email Already exist",
			"Response": res,
		})
		contactsLogger.Printf("[%s]Error occurred User email already exist %s", time.Now().Format(time.RFC3339), data.Email)
		return
	}
	if res == "Mobile" {
		c.JSON(http.StatusOK, gin.H{
			"Message":  "Mobile Number Already exist",
			"Response": res,
		})
		contactsLogger.Printf("[%s]Error occurred User Mobile already exist %s", time.Now().Format(time.RFC3339), data.Mobile)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message":  "Data Inserted Successfully",
		"Response": res,
	})
	contactsLogger.Printf("[%s]New Contact inserted successfully %s", time.Now().Format(time.RFC3339), data.Name)
}

func GetContactTableData(c *gin.Context) {
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
	res, err := models.GetContactTableData(data.ID, data.Team, data.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching contact data from DB : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	})
	// contactsLogger.Printf("[%s] Contact data fetched successfully ", time.Now().Format(time.RFC3339))
}

func GetContactData(c *gin.Context) {
	initLoggers()
	var Data models.Contacts

	if err := c.ShouldBindJSON(&Data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		contactsLogger.Printf("[%s] Invalid request data while get contacts detailed data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	res, err := models.GetContactData(Data.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while data fetching from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}
	mail, found, err := models.FetchMailData(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching mailing data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	note, foundNote, err := models.FetchNotesData(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching notes data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	meet, foundMeet, err := models.FetchMeetingData(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching meeting logs data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}
	call, foundCall, err := models.FetchCallsData(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	task, foundTask, err := models.FetchTaskData(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	deals, foundDeals, err := models.FetchDeals(int(Data.ID), 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		contactsLogger.Printf("[%s] Error occurred while fetching calls data from DB: %s", time.Now().Format(time.RFC3339), err)
		return
	}

	response := gin.H{
		"message": "Data fetched successfully",
		"data":    res,
	}

	if foundDeals {
		response["deals"] = deals
	} else {
		response["deals"] = "No Deals Found till now!"
	}

	if foundCall {
		response["calls"] = call
	} else {
		response["calls"] = "No calls Found till now!"
	}

	if foundMeet {
		response["meetings"] = meet
	} else {
		response["meetings"] = "No meeting logs Found till now!"
	}

	if foundTask {
		response["tasks"] = task
	} else {
		response["tasks"] = "No tasks Found till now!"
	}

	if foundNote {
		response["notes"] = note
	} else {
		response["notes"] = "No Notes Found till now!"
	}

	if found {
		response["mail"] = mail
	} else {
		response["mail"] = "No Mail Data Found"
	}
	c.JSON(http.StatusOK, response)
}

func EditContactData(c *gin.Context) {
	initLoggers()
	var Data models.Contacts
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Invalid request data while get contact details editing data : %s", time.Now().Format(time.RFC3339), err)
		return
	}

	res, err := models.EditContactData(Data.ID, Data.Name, Data.Email, Data.Mobile, Data.LifeCycleStage, Data.LeadStatus, Data.ContactLevelType, Data.JobRole, Data.EditedBY)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] error occurred while executing contact details editing data : %s : The contact person name %s", time.Now().Format(time.RFC3339), err, Data.Name)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Data":    res,
		"message": "Data updated successfully",
	})
	contactsLogger.Printf("[%s]Contact Updated successfully %s", time.Now().Format(time.RFC3339), Data.Name)
}

func DeleteContact(c *gin.Context) {
	initLoggers()
	var Data models.Contacts
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Invalid request data while get contact details deleting data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	err = models.DeleteContact(Data.ID, Data.IsCompany, Data.DeletedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		contactsLogger.Printf("[%s] error occurred while executing contact details deleting data : %s : The contact person name %s", time.Now().Format(time.RFC3339), err, Data.Name)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data deleted successfully",
	})
}

func InsertCompanyContact(c * gin.Context) {
	initLoggers()
	var Data models.Contacts
	if err := c.ShouldBindJSON(&Data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid data format!",
			"Error":   err.Error(),
		})
		contactsLogger.Printf("[%s] Invalid request data while get contact details inserting data : %s", time.Now().Format(time.RFC3339), err)
		return
	}
	err = models.InsertCompanyContact(Data.CompanyID, Data.Name, Data.Email, Data.Mobile, Data.Designation, Data.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted successfully",
	})
}