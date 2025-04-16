package utils

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
)

var db *sql.DB

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("DSN:", dsn)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func GetDB() *sql.DB {
	return db
}

func createEmailBody(recipient, subject, body string, attachments []Attachment, ccEmails string, team string) ([]byte, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// Add headers
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", recipient))
	if ccEmails != "" {
		buffer.WriteString(fmt.Sprintf("Cc: %s\r\n", ccEmails))
	}
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	buffer.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
	buffer.WriteString("\r\n")

	// Add body
	bodyWriter, _ := writer.CreatePart(textproto.MIMEHeader{
		"Content-Type": {"text/plain; charset=UTF-8"},
	})
	bodyWriter.Write([]byte(body))

	// Add attachments
	for _, att := range attachments {
		attWriter, _ := writer.CreatePart(textproto.MIMEHeader{
			"Content-Disposition":       {fmt.Sprintf(`attachment; filename="%s"`, att.Filename)},
			"Content-Type":              {"application/pdf"},
			"Content-Transfer-Encoding": {"base64"},
		})
		encoded := base64.StdEncoding.EncodeToString(att.Data)
		attWriter.Write([]byte(encoded))
	}

	writer.Close()
	return buffer.Bytes(), nil
}

func SendEmail(service *gmail.Service, recipient, subject, body string, attachments []Attachment, Team string, DealID int, Name, Email string, CC []string) error {
	var message gmail.Message
	fmt.Print("send mail process started")

	// Join the CC emails (if any) into a comma-separated string
	ccEmails := strings.Join(CC, ", ")

	// Create the email body with attachments
	emailBody, err := createEmailBody(recipient, subject, body, attachments, ccEmails, Team) // Pass joined CC here
	if err != nil {
		err = UpdateStatus(DealID, "Mail Sending Failed", Name, Email, recipient, subject, body)
		fmt.Print("Got error while creating email content", err)
		return fmt.Errorf("error creating email body: %v", err)
	}

	message.Raw = base64.URLEncoding.EncodeToString(emailBody)

	response, err := service.Users.Messages.Send("me", &message).Do()
	if err != nil {
		fmt.Printf("Failed to send email via Gmail API: %v\n", err)
		UpdateStatus(DealID, "Mail Sending Failed", Name, Email, recipient, subject, body)
		return fmt.Errorf("unable to send email: %v", err)
	} else {
		fmt.Printf("Email sent successfully: %v\n", response)
	}

	err = UpdateStatus(DealID, "Mail sent Successfully", Name, Email, recipient, subject, body)
	fmt.Print("Mail sent successfully..")
	return nil
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

func UpdateStatus(ID int, Status, Name, Email, recipient, subject, body string) error {
	Query := `UPDATE Deals SET MailStatus = ? WHERE ID = ?`
	_, err := db.Exec(Query, Status, ID)
	if err != nil {
		return err
	}
	var IsCompany, CompanyID int
	CompanyQuery := `SELECT IsCompany, CompanyID FROM Deals WHERE ID = ?`
	err = db.QueryRow(CompanyQuery, ID).Scan(&IsCompany, &CompanyID)
	if err != nil {
		return err
	}

	InsertQuery := `INSERT INTO mail(CompanyID, IsCompany, MailFrom, MailTo, MailSubject, MailContent, Status, CreatedBy, CreatedAt) 
					VALUES(?, ?, ?, ?, ?, ?, ?, ?, Now())`
	_, err = db.Exec(InsertQuery, CompanyID, IsCompany, Email, recipient, subject, body, Status, Name)
	if err != nil {
		return err
	}

	return nil
}
