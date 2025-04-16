package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func exchangeCodeForToken(code string) (string, string, error) {
	tokenEndpoint := googleOauthConfig.Endpoint.TokenURL

	reqBody := map[string]string{
		"code":          code,
		"client_id":     googleOauthConfig.ClientID,
		"client_secret": googleOauthConfig.ClientSecret,
		"redirect_uri":  googleOauthConfig.RedirectURL,
		"grant_type":    "authorization_code",
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	resp, err := http.Post(tokenEndpoint, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", "", fmt.Errorf("failed to make request to token endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status code from token endpoint: %d", resp.StatusCode)
	}

	var tokenResp GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", "", fmt.Errorf("failed to decode token response: %v", err)
	}

	if tokenResp.IDToken == "" {
		return "", "", fmt.Errorf("no ID token in token response")
	}

	return tokenResp.IDToken, tokenResp.RefreshToken, nil
}

func SignJWT(userInfo map[string]interface{}) (string, error) {
	now := time.Now().Unix() // Current time in seconds
	expiration := now + 3600
	clientSecret := os.Getenv("CLIENT_SECRET")
	claims := jwt.MapClaims{
		"sub":   userInfo["id"],
		"name":  userInfo["name"],
		"email": userInfo["email"],
		"iss":   "oauth-app-golang",
		"iat":   now,
		"exp":   expiration,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	log.Printf("JWT generated with iat: %v and exp: %v", now, expiration)
	return signedToken, nil
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
}

func refreshAccessToken(refreshToken string) (string, error) {
	tokenEndpoint := googleOauthConfig.Endpoint.TokenURL

	reqBody := map[string]string{
		"client_id":     googleOauthConfig.ClientID,
		"client_secret": googleOauthConfig.ClientSecret,
		"refresh_token": refreshToken,
		"grant_type":    "refresh_token",
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

	if tokenResp.AccessToken == "" {
		return "", fmt.Errorf("no access token in token response")
	}

	return tokenResp.AccessToken, nil
}
