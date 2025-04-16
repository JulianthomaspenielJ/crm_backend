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

var (
	activityLog *log.Logger
)

func activityLoggers() {
	logFile, err = os.OpenFile("activity.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening activity log file:", err)
	}
	activityLog = log.New(logFile, "ACTIVITY: ", log.LstdFlags|log.Lshortfile)
}

func InsertCalls(c *gin.Context) {
	activityLoggers()
	var calls models.CallsStruct
	if err := c.ShouldBindJSON(&calls); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while inserting call log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.InsertCall(calls.CallFrom, calls.CallDescription, calls.CallOutCome, calls.CallAt, calls.CallDirection, calls.CreatedBy, calls.CallTo, calls.CreatedUser, calls.IsCompany, calls.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while inserting call log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted successfully",
	})
}

func InsertNotes(c *gin.Context) {
	activityLoggers()
	var notes models.NotesStruct
	if err := c.ShouldBindJSON(&notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while inserting notes %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.InsertNotes(notes.NoteDescription, notes.CreatedBy, notes.IsCompany, notes.CompanyID, notes.CreatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while inserting notes %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted successfully",
	})
	activityLog.Printf("[%s] Notes Inserted notes %s", time.Now().Format(time.RFC3339), notes.CreatedBy)

}

func InsertTask(c *gin.Context) {
	activityLoggers()
	var task models.TaskStruct
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while inserting task %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.InsertTask(task.TaskDescription, task.TaskType, task.DueDate, task.CreatedBy, task.CreatedUser, task.IsCompany, task.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while inserting task %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task inserted successfully",
	})
}

func InsertMeeting(c *gin.Context) {
	activityLoggers()
	var meetings models.MeetingStruct
	if err := c.ShouldBindJSON(&meetings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while inserting meeting log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.InsertMeeting(meetings.Duration, meetings.CreatedBy, meetings.Attendees, meetings.OutCome, meetings.MeetingAt, meetings.MeetingDescription, meetings.CreatedBy, meetings.CreatedUser, meetings.IsCompany, meetings.CompanyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while inserting meeting log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Meeting log inserted successfully",
	})
}

//UPDATE code of activities modules

func UpdateNotes(c *gin.Context) {
	activityLoggers()
	var notes models.NotesStruct
	if err := c.ShouldBindJSON(&notes); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while Updating notes %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.UpdateNotes(notes.NoteDescription, notes.UpdatedBy.String, notes.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while Updating notes %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data Updated successfully",
	})
	activityLog.Printf("[%s] Notes Updated notes %s", time.Now().Format(time.RFC3339), notes.UpdatedBy.String)
}

func UpdateCalls(c *gin.Context) {
	activityLoggers()
	var calls models.CallsStruct
	if err := c.ShouldBindJSON(&calls); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while updating call log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.UpdateCall(calls.CallDescription, calls.CallOutCome, calls.CallAt, calls.CallDirection, calls.UpdatedBy.String, calls.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while updating call log %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
	})
}

func UpdateTask(c *gin.Context) {
	activityLoggers()
	var task models.TaskStruct
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error":   err.Error(),
			"message": "Bad Request! Please provide valid data",
		})
		activityLog.Printf("[%s] Invalid request while inserting task %s", time.Now().Format(time.RFC3339), err)
		return
	}

	err := models.UpdateTask(task.TaskDescription, task.TaskType, task.DueDate, task.UpdatedBy.String, task.Status, task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Something Went wrong! Please try again",
		})
		activityLog.Printf("[%s] Error occurred while inserting task %s", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task inserted successfully",
	})
}

func Attendees(c *gin.Context) {
	res, err := models.GetAttendees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred",
			"error":   err,
		})
		activityLog.Printf("[%s] Error occurred while Getting attendees info %s", time.Now().Format(time.RFC3339), err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched",
		"data":    res,
	})
}
