package routes

import (
	"context"
	"log"
	"os"
	"strings"

	// "crm/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

var tokenCache = make(map[string]string)

type Attachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

type RequestBody struct {
	Recipient    string       `json:"recipient"`
	Subject      string       `json:"subject"`
	Body         string       `json:"body"`
	Attachments  []Attachment `json:"attachments"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	Email        string       `json:"email"`
	CC           []string     `json:"cc"`
}

func SendEmailHandler(c *gin.Context) {
	var requestBody RequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Refresh access token if needed
	newToken, err := GetCachedToken(requestBody.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to refresh token"})
		return
	}
	requestBody.AccessToken = newToken

	// Enqueue email
	if err := enqueueEmail(requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue email"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Email queued successfully!"})
}

func GetCachedToken(refreshToken string) (string, error) {
	if token, ok := tokenCache[refreshToken]; ok {
		return token, nil
	}

	// Refresh token if not in cache
	newToken, err := RefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	// Store new token in cache with expiration
	tokenCache[refreshToken] = newToken
	go func() {
		time.Sleep(60 * time.Minute) // Expiration period (adjust as needed)
		delete(tokenCache, refreshToken)
	}()
	return newToken, nil
}

func RefreshToken(refreshToken string) (string, error) {
	config, err := ReadCredentials()
	if err != nil {
		return "", fmt.Errorf("unable to load OAuth2 config: %v", err)
	}

	tokenSource := config.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	})

	newToken, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("unable to refresh access token: %v", err)
	}
	return newToken.AccessToken, nil
}

func ReadCredentials() (*oauth2.Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURL := os.Getenv("ORIGIN")
	authURL := os.Getenv("AUTH_URL")
	tokenURL := os.Getenv("TOKEN_URL")
	scopes := os.Getenv("SCOPES")

	if scopes == "" {
		scopes = "https://www.googleapis.com/auth/gmail.send"
	} else {
		scopes = scopes + ",https://www.googleapis.com/auth/gmail.send"
	}

	scopeList := []string{}
	for _, scope := range strings.Split(scopes, ",") {
		scopeList = append(scopeList, scope)
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopeList,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	return config, nil
}

func SendCustomerMail(c *gin.Context) {
	var mailData RequestBody

	if err := c.ShouldBindJSON(&mailData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// fmt.Printf("mailData: %+v\n", mailData)
	err = enqueueEmail(mailData)
	if err != nil {
		fmt.Print("error occured while inserting into queue", err)
		return
	}

}
