package routes

import (
	"crm/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsersData(c *gin.Context) {
	res, err := models.GetUsersData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data fetched successfully",
		"data":    res,
	})
}

func UpdateUserManagementData(c *gin.Context) {
	var users models.UserDetails

	if c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data!",
			"error":   err.Error(),
		})
		return
	}
	err := models.UpdateUserManagementData(users.ID, users.Name, users.Email, users.Role, users.Designation, users.Team, users.Mobile, users.EditedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
	})
}

func GetDeletedUsersData(c *gin.Context) {
	res, err := models.GetDeletedUsersData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data fetched successfully",
		"data":    res,
	})
}

func RestoreUserManagementData(c *gin.Context) {
	var users models.UserDetails
	if c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data!",
			"error":   err.Error(),
		})
		return
	}
	err := models.RestoreUserManagementData(users.ID, users.EditedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data restored successfully",
	})
}

func RemoveUser(c *gin.Context) {
	var users models.UserDetails
	if c.ShouldBindJSON(&users); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data!",
			"error":   err.Error(),
		})
		return
	}
	err := models.RemoveUser(users.ID, users.EditedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data removed successfully",
	})
}

func UserNameList(c *gin.Context) {
	users, err := models.UserNameList()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    users,
	})
}

func UserNameListDashboards(c *gin.Context) {
	users, err := models.UserNameListDashboards()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   err.Error(),
			"message": "Something went wrong! Please try again.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched successfully",
		"data":    users,
	})
}

type FCMStruct struct {
	FCMToken string `json:"token"`
	UserID   int    `json:"userID"`
}

func StoreFCMToken(c *gin.Context) {
	var user FCMStruct
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input data!",
			"error":   err.Error(),
		})
		return
	}

	err := models.StoreFCMToken(user.UserID, user.FCMToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong! Please try again.",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
	})
}

func GetTeam(c *gin.Context) {
	res, err := models.GetTeam()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong!",
			"error":   err.Error(),
		})
		fmt.Print("error", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "data fetched successfully",
		"data":    res,
	})
}
