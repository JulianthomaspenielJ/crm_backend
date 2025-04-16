package login

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	//"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

var db *sql.DB
var (
	googleOauthConfig = &oauth2.Config{}
)

func init() {

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		RedirectURL:  os.Getenv("ORIGIN"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
}

func LoginHandler(c *gin.Context) {
	var request struct {
		Code string `json:"code"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	idToken, err := exchangeCodeForToken(request.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	var isApproved, isDeleted int
	query := `SELECT IsApproved, IsDeleted FROM users WHERE Email = ?`
	err = db.QueryRow(query, email).Scan(&isApproved, &isDeleted)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in the database"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	if isApproved == 1 && isDeleted == 0 {
		signedToken, err := SignJWT(map[string]interface{}{
			"id":    userRecord.UID,
			"name":  userRecord.DisplayName,
			"email": email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Something went wrong! Please try again"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": signedToken, "message": "Login successful", "Data": email})
	} else if isApproved == 0 && isDeleted == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User waiting for admin approval", "data": email, "uid": userRecord.UID})
		return
	} else if isApproved == 0 && isDeleted == 1 {
		c.JSON(http.StatusForbidden, gin.H{"error": "User no longer active!"})
		return
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Waiting for admin approval"})
		return
	}
}

func exchangeCodeForToken(code string) (string, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	redirectURI := os.Getenv("ORIGIN")

	tokenEndpoint := "https://oauth2.googleapis.com/token"

	reqBody := map[string]string{
		"code":          code,
		"client_id":     clientID,
		"client_secret": clientSecret,
		"redirect_uri":  redirectURI,
		"grant_type":    "authorization_code",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(tokenEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to make request to token endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code from token endpoint: %d", resp.StatusCode)
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %v", err)
	}

	if tokenResp.IDToken == "" {
		return "", fmt.Errorf("no ID token in token response")
	}

	return tokenResp.IDToken, nil
}

func SignJWT(userInfo map[string]interface{}) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userInfo["id"],
		"name":  userInfo["name"],
		"email": userInfo["email"],
		"iss":   "oauth-app-golang",
		"exp":   time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("JWT token signing error: %v", err)
	}
	return signedToken, nil
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"` // Add this line
	IDToken      string `json:"id_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}
