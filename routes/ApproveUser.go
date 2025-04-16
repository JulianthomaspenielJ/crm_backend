package routes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"crm/models"
)

var db *sql.DB

func init() {
	db = models.GetDB()
}

func ApproveUserList(c *gin.Context) {
	rows, err := db.Query("SELECT ID, Name, Email, Mobile, IsApproved, IsDeleted FROM users WHERE IsApproved = 0 AND IsDeleted = 0")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}

	for rows.Next() {
		var ID int
		var Name, Email, IsApproved, IsDeleted string
		var Mobile sql.NullString

		err := rows.Scan(&ID, &Name, &Email, &Mobile, &IsApproved, &IsDeleted)
		if err != nil {
			fmt.Println("Error scanning user row:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan users"})
			return
		}

		var mobileValue string
		if Mobile.Valid {
			mobileValue = Mobile.String
		} else {
			mobileValue = ""
		}

		users = append(users, map[string]interface{}{
			"id":         ID,
			"name":       Name,
			"email":      Email,
			"mobile":     mobileValue,
			"IsApproved": IsApproved,
			"IsDeleted":  IsDeleted,
		})
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to iterate over users",
			"message": "Something went wrong! Please try again",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":   users,
		"message": "User approved successfully",
	})
}
