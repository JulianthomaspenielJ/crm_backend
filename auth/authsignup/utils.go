package signup

import (
	"context"
	"crm/models"
	"crm/routes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/api/option"
)

func StoreUserInfo(userInfo map[string]interface{}) error {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./oAuth.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		return fmt.Errorf("failed to create Firebase app %v", nil)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("failed to create Firebase Auth client %v", nil)
	}

	email := userInfo["email"].(string)
	_, err = client.GetUserByEmail(ctx, email)
	if err == nil {
		return fmt.Errorf("user already registered")
	} else if !auth.IsUserNotFound(err) {
		return fmt.Errorf("error fetching user %v", err)
	}

	params := (&auth.UserToCreate{}).
		Email(email).
		DisplayName(userInfo["name"].(string))
	userRecord, err := client.CreateUser(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to create user %v", err)
	}
	userInfo["uid"] = userRecord.UID

	err = StoreUserInfoInDB(userInfo)
	if err != nil {
		return fmt.Errorf("failed to store user info in DB %v", err)
	}

	return nil
}

func StoreUserInfoInDB(userInfo map[string]interface{}) error {
	db := models.GetDB()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction %v", err)
	}

	query := `INSERT INTO users (Name, Email, AuthUID, IsApproved, IsDeleted)
              VALUES (?, ?, ?, ?, ?)
              ON DUPLICATE KEY UPDATE Name=VALUES(Name), Email=VALUES(Email), AuthUID=VALUES(AuthUID), IsApproved=VALUES(IsApproved), IsDeleted=VALUES(IsDeleted)`

	_, err = tx.Exec(query,
		userInfo["name"], userInfo["email"], userInfo["uid"], 0, 0)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute query %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction %v", err)
	}
	var Mail models.DealsStruct
	var Email string
	AdminQuery := `SELECT Email FROM users WHERE Role = ? AND IsDeleted = ?`
	err = db.QueryRow(AdminQuery, "Super Admin", 0).Scan(&Email)
	if err != nil {
		fmt.Print("Error while fetching admin email", err)
	}
	// Mail.From = "crm@holoware.co"
	Mail.To = []string{Email}
	fmt.Print("TO PERSON MAIL", Mail.To)
	Mail.From = userInfo["email"].(string)
	Mail.DealOwner = userInfo["name"].(string)
	Mail.CC = []string{Mail.From}
	Mail.Subject = "New User Request for CRM"
	Body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello Sir,</p>
				<p>A new user <strong>%s (%s)</strong> has been requested to Access to the CRM. And we need your approval to proceed with setting up their login credentials.</p> <br/>
				<p>Looking forward to your confirmation.<p>
				<p>Thanks & Regards,<br>Team CRM</p>
			</body>
		</html>
		`,
		Mail.DealOwner, Mail.From,
	)
	Mail.Body = Body
	// 	fmt.Printf("The user mail %s To mail %s", Mail.Email, Mail.To)
	err = routes.DealsMail(Mail)
	return nil
}

func GetUserInfo(accessToken string) (map[string]interface{}, error) {
	userInfoEndPoint := "https://www.googleapis.com/oauth2/v2/userinfo"
	response, err := http.Get(fmt.Sprintf("%s?access_token=%s", userInfoEndPoint, accessToken))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
		return nil, err
	}
	return userInfo, nil
}
func SignJWT(userInfo map[string]interface{}) (string, error) {
	clientSecret := os.Getenv("CLIENT_SECRET")

	claims := jwt.MapClaims{
		"sub":   userInfo["id"],
		"name":  userInfo["name"],
		"email": userInfo["email"],
		"iss":   "oauth-app-golang",
		"exp":   time.Now().Add(time.Minute * 60).Unix(), // Set expiration to 60 minutes
		"iat":   time.Now().Unix(),                       // Issued at time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(clientSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}
