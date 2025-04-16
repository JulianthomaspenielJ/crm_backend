package routes

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"crm/models"

	firebase "firebase.google.com/go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

var googleOauthConfig *oauth2.Config

func init() {
	db = models.GetDB()
	googleOauthConfig = models.GetGoogleOauthConfig()
}

func LoginHandler(c *gin.Context) {
	fmt.Println("It was reach there")
	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idToken, refreshToken, err := exchangeCodeForToken(request.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Print("refresh token", refreshToken)
	payload, err := idtoken.Validate(context.Background(), idToken, os.Getenv("CLIENT_ID"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	email := payload.Claims["email"].(string)

	ctx := context.Background()
	sa := option.WithCredentialsFile("./oAuth.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Firebase"})
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Firebase Auth client"})
		return
	}

	userRecord, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in Firebase"})
		return
	}

	user, err := models.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "user not found in the database" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else {
			fmt.Print("error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error", "err": err})
		}
		return
	}

	if user.IsApproved == 1 && user.IsDeleted == 0 {
		signedToken, err := SignJWT(map[string]interface{}{
			"id":    userRecord.UID,
			"name":  userRecord.DisplayName,
			"email": email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong! Please try again"})
			return
		}

		userData := map[string]interface{}{
			"uid":          userRecord.UID,
			"name":         userRecord.DisplayName,
			"email":        email,
			"userName":     user.UserName.String,
			"userRole":     user.UserRole.String,
			"mobile":       user.Mobile.String,
			"designation":  user.Designation.String,
			"avatarURL":    user.AvatarURL.String,
			"signatureURL": user.SignatureURL.String,
			"team":         user.Team.String,
			"id":           user.ID,
			"refreshToken": refreshToken,
		}

		c.JSON(http.StatusOK, gin.H{"token": signedToken, "message": "Login successful", "userData": userData})
	} else if user.IsApproved == 0 && user.IsDeleted == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User waiting for admin approval", "data": email, "uid": userRecord.UID})
		return
	} else if user.IsApproved == 0 && user.IsDeleted == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User no longer active!"})
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Waiting for admin approval"})
		return
	}
}
