package signup

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"crm/models"
)

var googleOauthConfig *oauth2.Config

func init() {
	googleOauthConfig = models.GetGoogleOauthConfig()
}

func GoogleCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("Error exchanging code: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error exchanging code! Please try again"})
		return
	}

	userInfo, err := GetUserInfo(token.AccessToken)
	if err != nil {
		fmt.Println("Error getting user info: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error getting user info! Please try again"})
		return
	}

	email := userInfo["email"].(string)
	if !strings.HasSuffix(email, "@holoware.co") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized domain! Please SignUp with holoware.co account"})
		return
	}

	err = StoreUserInfo(userInfo)
	if err != nil {
		if strings.Contains(err.Error(), "already registered") {
			c.JSON(http.StatusConflict, gin.H{"error": "The user email is already registered!"})
		} else {
			fmt.Println("Error storing user info: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error storing user info! Please try again"})
		}
		return
	}

	signedToken, err := SignJWT(userInfo)
	if err != nil {
		fmt.Println("Error signing token: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong! Please try again"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": signedToken})
}
