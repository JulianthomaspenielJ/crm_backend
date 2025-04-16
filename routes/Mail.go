package routes

import (
	"crm/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

var redisClient *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
	})
	fmt.Printf("Connecting to Redis: Host=%s Port=%s Password=%s\n", redisHost, redisPort, redisPassword)

}

type EmailStruct struct {
	Recipient    string       `json:"recipient"`
	Subject      string       `json:"subject"`
	Body         string       `json:"body"`
	Attachments  []Attachment `json:"attachments"`
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	Team         string       `json:"team"`
	DealID       int          `json:"dealID"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	CC           []string     `json:"cc"`
}

type MailStruct struct {
	From        string   `json:"fromEmail"`
	To          []string `json:"toMail"`
	CC          []string `json:"cc"`
	BCC         []string `json:"bcc"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments"`
	CreatedBy   string   `json:"createdBy"`
	IsCompany   int      `json:"isCompany"`
	CompanyID   int      `json:"companyID"`
	CreatedUser int      `json:"CreatedUser"`
}

func Mail(c *gin.Context) {
	var Data MailStruct
	var err error

	logFile, err := os.OpenFile("mail_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	log.SetOutput(logFile)

	if err := c.ShouldBindJSON(&Data); err != nil {
		fmt.Println("Error binding JSON:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   err.Error(),
			"message": "Invalid request data. Please check your input and try again.",
		})
		return
	}

	if err := enqueueMail(Data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Error occurred",
			"err":     err,
		})
		log.Printf("[%s] Error occurred while enqueue the mail: %v\n", time.Now().Format(time.RFC3339), err)
		return
	}

	err = models.StoreMailContent(Data.From, strings.Join(Data.To, ","), Data.Subject, Data.Body, strings.Join(Data.CC, ","), Data.CreatedBy, Data.IsCompany, Data.CompanyID, Data.CreatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   err.Error(),
			"message": "Error occurred while storing data",
		})
		log.Printf("[%s] Error occurred while storing mail data to the db: %v\n", time.Now().Format(time.RFC3339), err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Mail enqueued successfully",
	})
	log.Printf("[%s] Mail enqueued successfully, From: %s, To: %s\n",
		time.Now().Format(time.RFC3339), Data.From, strings.Join(Data.To, ", "))
}

func enqueueMail(mail MailStruct) error {
	mailData, err := json.Marshal(mail)
	if err != nil {
		return err
	}

	err = redisClient.Publish("ticket_queue", string(mailData)).Err()
	if err != nil {
		return err
	}

	return nil
}

func DealsMail(deal models.DealsStruct) error {

	mailData, err := json.Marshal(deal)
	if err != nil {
		return fmt.Errorf("failed to marshal mail data: %v", err)
	}

	err = redisClient.Publish("ticket_queue", string(mailData)).Err()
	if err != nil {
		return fmt.Errorf("failed	 to publish mail to Redis: %v", err)
	}

	return nil
}

func enqueueEmail(emailData RequestBody) error {
	mailData, err := json.Marshal(emailData)
	if err != nil {
		return err
	}
	err = redisClient.Publish("gmail_queue", string(mailData)).Err()
	if err != nil {
		return err
	}
	return nil
}

func RefreshTokenHandler(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refreshToken"`
	}

	// Bind the incoming JSON request to the requestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call your existing RefreshToken function
	newToken, err := RefreshToken(requestBody.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unable to refresh access token"})
		return
	}

	// Return the new access token in the response
	c.JSON(http.StatusOK, gin.H{
		"accessToken": newToken,
	})
}
