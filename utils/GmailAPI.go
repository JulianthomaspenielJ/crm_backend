package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type Attachment struct {
	Filename string `json:"filename"`
	Data     []byte `json:"data"`
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	// Exchange the authorization code for a token.
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// func ReadCredentials() (*oauth2.Config, error) {
// 	config := &oauth2.Config{
// 		ClientID:     "558462773191-hhl6ts0fh26k6a1e137a9r456nhneckk.apps.googleusercontent.com",
// 		ClientSecret: "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",
// 		RedirectURL:  "http://localhost:3000",
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/gmail.send",
// 		},
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
// 			TokenURL: "https://oauth2.googleapis.com/token",
// 		},
// 	}

// 	return config, nil
// }

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
	scopes := os.Getenv("SCOPES") // This could be a comma-separated list of scopes

	// Add the Gmail send scope explicitly
	if scopes == "" {
		scopes = "https://www.googleapis.com/auth/gmail.send"
	} else {
		scopes = scopes + ",https://www.googleapis.com/auth/gmail.send" // Append if not already present
	}

	// Convert the scopes from a comma-separated string to a slice
	scopeList := []string{}
	for _, scope := range strings.Split(scopes, ",") {
		scopeList = append(scopeList, scope)
	}

	// Create and return the OAuth2 config with the environment variables
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

func GetGmailServiceWithToken(accessToken string) (*gmail.Service, error) {
	ctx := context.Background()
	config, err := ReadCredentials()
	if err != nil {
		return nil, fmt.Errorf("unable to load OAuth2 config: %v", err)
	}

	tokenSource := config.TokenSource(ctx, &oauth2.Token{AccessToken: accessToken})
	return gmail.NewService(ctx, option.WithTokenSource(tokenSource))
}
