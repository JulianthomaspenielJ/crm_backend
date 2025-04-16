package routes

import (
	"crm/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddMargin(c *gin.Context) {
	var margin models.Margin
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input! Please give a proper input",
			"error":   err.Error(),
		})
		return
	}

	res, err := models.AddMargin(margin.MarginType, margin.BPPercentage, margin.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	if res == "Margin Type already Exist" {
		c.JSON(http.StatusOK, gin.H{
			"message": res,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted Successfully",
		"data":    res,
	})
}

func AddMarginValues(c *gin.Context) {
	var margin models.MarginValues
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Input! Please give a proper input",
			"error":   err.Error(),
		})

		return
	}

	res, err := models.AddMarginValues(margin.MarginID, margin.Tier, margin.TierPercentage, margin.CreatedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	if res == "Tier already Exist" {
		c.JSON(http.StatusOK, gin.H{
			"message": res,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data inserted Successfully",
		"data":    res,
	})
}

func GetMargin(c *gin.Context) {
	res, err := models.GetMargin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched Successfully",
		"data":    res,
	})
}

func GetMarginTableData(c *gin.Context) {
	res, err := models.GetMarginTableData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched Successfully",
		"data":    res,
	})
}

func GetMarginValuesTableData(c *gin.Context) {
	res, err := models.GetMarginValuesTableData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data fetched Successfully",
		"data":    res,
	})
}

func UpdateMargin(c *gin.Context) {
	var margin models.Margin
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please give a valid input",
			"error":   err.Error(),
		})
		return
	}
	err := models.UpdateMargin(margin.MarginType, margin.CreatedBy, margin.BPPercentage, margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully",
	})
}

func UpdateMarginValues(c *gin.Context) {
	var margin models.MarginValues
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please provide a valid input.",
			"error":   err.Error(),
		})
		return
	}

	res, err := models.UpdateMarginValues(margin.Tier, margin.UpdatedBy.String, margin.TierPercentage, margin.ID)
	if err != nil {
		if res == "exist" {
			c.JSON(http.StatusOK, gin.H{
				"message": "Tier name already exists.",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong.",
				"error":   err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data updated successfully.",
	})
}
func GetMarginValuesData(c *gin.Context) {
	var margin models.MarginValues
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please give a valid input",
			"error":   err.Error(),
		})

		return
	}
	res, err := models.GetMarginValuesData(margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Data FETCHED successfully",
		"data":    res,
	})
}

func DeleteMargin(c *gin.Context) {
	var margin struct {
		UpdatedBy string `json:"updatedBy"`
		ID        int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please provide valid input",
			"error":   err.Error(),
		})

		return
	}
	err := models.DeleteMargin(margin.UpdatedBy, margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Margin Deleted successfully",
	})
}

func DeleteMarginValues(c *gin.Context) {
	var margin struct {
		UpdatedBy string `json:"updatedBy"`
		ID        int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please provide valid input",
			"error":   err.Error(),
		})

		return
	}
	err := models.DeleteMarginValues(margin.UpdatedBy, margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Margin deleted successfully",
	})
}

func RecoverMargin(c *gin.Context) {
	var margin struct {
		UpdatedBy string `json:"updatedBy"`
		ID        int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please provide valid input",
			"error":   err.Error(),
		})

		return
	}
	err := models.RecoverMargin(margin.UpdatedBy, margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Margin recovered successfully",
	})
}

func RecoverMarginValues(c *gin.Context) {
	var margin struct {
		UpdatedBy string `json:"updatedBy"`
		ID        int    `json:"id"`
	}
	if err := c.ShouldBindJSON(&margin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input! Please provide valid input",
			"error":   err.Error(),
		})

		return
	}
	err := models.RecoverMarginValues(margin.UpdatedBy, margin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong",
			"error":   err.Error(),
		})

		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Margin Recovered successfully",
	})
}
