package queue

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"

	"gopkg.in/gomail.v2"
)

var redisClient *redis.Client
var logFile *os.File

type MailStruct struct {
	From        string   `json:"fromEmail"`
	To          []string `json:"toMail"`
	CC          []string `json:"cc"`
	BCC         []string `json:"bcc"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	Attachments []string `json:"attachments"`
}

func init() {
	var err error
	err = godotenv.Load()
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

	// var err error
	logFile, err = os.OpenFile("mail_service.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}
	log.SetOutput(logFile)
}

func SendGoMail(From string, To, CC, BCC []string, Subject, Body string, base64Attachments []string) (string, error) {
	fmt.Printf("From mail is %s, To %s, CC %s", From, To, CC)
	m := gomail.NewMessage()
	m.SetHeader("From", From)
	m.SetHeader("Reply-To", From)
	m.SetHeader("Sender", From)
	m.SetHeader("To", To...)
	if len(CC) > 0 {
		m.SetHeader("Cc", CC...)
	}
	if len(BCC) > 0 {
		m.SetHeader("Bcc", BCC...)
	}
	m.SetHeader("Subject", Subject)
	m.SetBody("text/html", Body)

	for _, base64Attachment := range base64Attachments {
		if base64Attachment != "" {
			const prefix = "data:"
			if strings.HasPrefix(base64Attachment, prefix) {
				parts := strings.SplitN(base64Attachment, ",", 2)
				if len(parts) == 2 {
					mimeTypeData := parts[0]
					base64Data := parts[1]

					mimeTypeParts := strings.SplitN(mimeTypeData, ";", 2)
					mimeType := mimeTypeParts[0][len(prefix):]
					fileExtension := extractFileExtension(mimeType)

					data, err := base64.StdEncoding.DecodeString(base64Data)
					if err != nil {
						fmt.Printf("Error decoding base64: %s\n", err)
						return "error", fmt.Errorf("failed to decode base64 attachment: %v", err)
					}

					attachmentName := "attachment" + fileExtension
					m.Attach(attachmentName, gomail.SetCopyFunc(func(w io.Writer) error {
						_, err := w.Write(data)
						if err != nil {
							fmt.Printf("Error writing attachment: %s\n", err)
						}
						return err
					}))
				}
			}
		}
	}

	d := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("ACCOUNT_MAIL"), os.Getenv("MAIL_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return "error", err
	}

	return "Mail sent successfully", nil
}

func extractFileExtension(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "application/pdf":
		return ".pdf"
	default:
		return ".bin"
	}
}

func ConsumeMailQueue() {
	subChannel := "ticket_queue"
	fmt.Print("REDIS BEGIN")
	sub1 := redisClient.Subscribe(subChannel)
	subCh := sub1.Channel()

	for msg := range subCh {
		var mailData MailStruct
		if err := json.Unmarshal([]byte(msg.Payload), &mailData); err != nil {
			log.Printf("[%s] Unmarshal error: %v\n", time.Now().Format(time.RFC3339), err)
			continue
		}
		result, err := SendGoMail(mailData.From, mailData.To, mailData.CC, mailData.BCC, mailData.Subject, mailData.Body, mailData.Attachments)
		fmt.Printf("Result: FROM MAIL %s, To %s, cc %s, \n", mailData.From, mailData.To, mailData.CC)
		if err != nil {
			log.Printf("[%s] Error sending email: %v\n", time.Now().Format(time.RFC3339), err)
		} else {
			log.Printf("[%s] %s\n", time.Now().Format(time.RFC3339), result)
		}
	}
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

type Attachment struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`
}

type Attachments struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

func StartEmailWorker() {
	fmt.Println("Worker started, waiting for emails...")
	channel := "gmail_queue"
	pubsub := redisClient.Subscribe(channel)
	defer pubsub.Close()

	subCh := pubsub.Channel()

	fmt.Println("CHANNEL CONNECTED")
	for msg := range subCh {
		var emailData EmailStruct

		if err := json.Unmarshal([]byte(msg.Payload), &emailData); err != nil {
			log.Printf("[%s] Unmarshal error: %v\n", time.Now().Format(time.RFC3339), err)
			continue
		}

		var attachments []Attachments
		for _, att := range emailData.Attachments {
			data, err := base64.StdEncoding.DecodeString(att.Data)
			if err != nil {
				log.Println("Invalid base64 data for attachment:", err)
				continue
			}
			attachments = append(attachments, Attachments{
				Filename: att.Filename,
				Data:     data,
			})
		}

		err := SendMailing(emailData.Email, emailData.Recipient, emailData.Subject, emailData.Body, attachments, emailData.CC)
		if err != nil {
			log.Println("Failed to send email:", err)
		} else {
			fmt.Println("Email sent successfully")
		}
	}
}

func createEmailBody(Email, Recipient, Subject, Body string, Attachments []Attachments, CC []string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fmt.Print(Email, "Email id ")
	// Add headers
	buffer.WriteString(fmt.Sprintf("From: %s\r\n", Email))
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", Recipient))

	buffer.WriteString(fmt.Sprintf("Reply-To: %s\r\n", Email)) // Always set Reply-To to Email

	if len(CC) > 0 {
		ccEmails := strings.Join(CC, ", ")
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", ccEmails))
	}

	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", Subject))

	// Set boundary for the email content
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	buffer.WriteString("\r\n")

	// Add the body of the email (plain text)
	bodyWriter, _ := writer.CreatePart(map[string][]string{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	bodyWriter.Write([]byte(Body))

	fmt.Printf("Attachments: %s\n", Attachments)
	for _, att := range Attachments {
		// fmt.Printf("Attachment: %s\n", att.Filename)
		attWriter, _ := writer.CreatePart(map[string][]string{
			"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, att.Filename)},
			"Content-Type":              {"application/pdf"},
			"Content-Transfer-Encoding": {"base64"},
		})

		encoded := base64.StdEncoding.EncodeToString(att.Data) // Encoding the attachment data to base64
		attWriter.Write([]byte(encoded))                       // Write the base64 encoded data to the part
	}

	// Close the multipart writer
	writer.Close()

	return buffer.Bytes(), nil
}

func SendMailing(Email, Recipient, Subject, Body string, Attachments []Attachments, CC []string) error {
	// fmt.Printf("Attachments from send mail: %s\n", Attachments)
	emailBody, err := createEmailBody(Email, Recipient, Subject, Body, Attachments, CC)
	if err != nil {
		return fmt.Errorf("error creating email body: %v", err)
	}

	// SMTP server configuration
	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	smtpUser := os.Getenv("ACCOUNT_MAIL")
	smtpPass := os.Getenv("MAIL_PASSWORD") // or use an App Password for Gmail if 2FA is enabled

	// Setup authentication
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpServer)

	msg := emailBody
	to := []string{Recipient}
	if len(CC) > 0 {
		to = append(to, CC...)
	}

	// Send the email
	err = smtp.SendMail(smtpServer+":"+smtpPort, auth, smtpUser, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}
