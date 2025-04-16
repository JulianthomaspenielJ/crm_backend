package models

import (
	// "crm/routes"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func UpdateComponents(id int, name string, description string, price int64, brand string, updatedBy string) error {
	query := `UPDATE components SET VarientName=?, VarientDescription=?, VarientPrice=?, VarientBrand=?, UpdatedBy=?, UpdatedAt=NOW() WHERE ID=?`
	result, err := db.Exec(query, name, description, price, brand, updatedBy, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows are affected %v", err)
	}

	err = UpdateAudit("Update", updatedBy, id, "127.0.0.1", "Components table data updated")
	if err != nil {
		return err
	}

	return nil
}

func DeleteComponent(id int, UpdatedBy string) error {
	query := `UPDATE components SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`

	result, err := db.Exec(query, 1, UpdatedBy, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows are affected %v", err)
	}

	err = UpdateAudit("Delete", UpdatedBy, id, "127.0.0.1", "Components table data Deleted")
	if err != nil {
		return err
	}

	return nil
}

func UpdateUser(db *sql.DB, data UpdateUserData) (int64, DealsStruct, error) {
	var Mail DealsStruct
	query := `
        UPDATE users 
        SET Name =?, Mobile = ?, Role = ?, Designation = ?, IsApproved = ?, CreatedBy = ?, CreatedAt = ?, Team = ?
        WHERE ID = ?
    `
	currentTime := time.Now().Format("2006/01/02 15:04:05")

	result, err := db.Exec(query, data.Name, data.Mobile, data.Role, data.Designation, 1, data.CreatedBy, currentTime, data.Team, data.ID)
	if err != nil {
		return 0, Mail, fmt.Errorf("failed to update user: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, Mail, fmt.Errorf("failed to get rows affected: %v", err)
	}

	var Email string
	GetQuery := `SELECT Email FROM users WHERE ID = ?`
	err = db.QueryRow(GetQuery, data.ID).Scan(&Email)
	if err != nil {
		fmt.Printf("Error while fetching email: %v", err)
		return 0, Mail, nil
	}

	AdminQuery := `SELECT Email FROM users WHERE Role = ? AND IsDeleted = ?`
	err = db.QueryRow(AdminQuery, "Super Admin", 0).Scan(&Email)
	if err != nil {
		fmt.Print("Error while fetching admin email", err)
	}
	// Mail.From = "crm@holoware.co"
	Mail.To = []string{Email}
	fmt.Print("TO PERSON MAIL", Mail.To)
	Mail.From = "crm@holoware.co"
	Mail.DealOwner = data.CreatedBy
	Mail.CC = []string{Mail.From}
	Mail.Subject = "Access Approved for CRM"
	Body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello %s,</p>
				<p>As per your request you have been approved in to the CRM Application by <strong>%s</strong> And you are role will be %s</p> <br/>
				<p>Wish you all the very best!!<p>
				<p>Thanks & Regards,<br>Team CRM</p>
			</body>
		</html>

		`,
		data.Name, data.CreatedBy, data.Role,
	)
	Mail.Body = Body
	// 	fmt.Printf("The user mail %s To mail %s", Mail.Email, Mail.To)
	// err = routes.DealsMail(Mail)
	return rowsAffected, Mail, nil
}

func UpdateUserDetails(Name string, Email string, mobile string, avatarURL string) (int64, error) {
	query := `UPDATE users SET Mobile=?, AvatarURL=?, UpdatedBy=?, UpdatedAt= Now() WHERE Email = ?`
	result, err := db.Exec(query, mobile, avatarURL, Name, Email)
	if err != nil {
		return 0, fmt.Errorf("error occurred while executing query: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("no rows are affected: %v", err)
	}

	var userID int
	GetIDQuery := `SELECT ID FROM users WHERE Email = ?`
	err = db.QueryRow(GetIDQuery, Email).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("error occurred while fetching ID: %v", err)
	}

	err = UpdateAudit("Update", Name, userID, "127.0.0.1", "User data updated")
	if err != nil {
		return 0, fmt.Errorf("error occurred while updating audit: %v", err)
	}
	return rowsAffected, nil
}

func RecoverComponent(ID int, UpdatedBy string) error {
	query := `UPDATE components SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	rows, err := db.Exec(query, 0, UpdatedBy, ID)
	if err != nil {
		return fmt.Errorf("error occurred while query executing %s", err)
	}
	fmt.Print(rows)
	err = UpdateAudit("Restored", UpdatedBy, ID, "127.0.0.1", "Components table data Restored")
	if err != nil {
		return err
	}
	return nil
}

func UpdatePricing(ID int, data PricingProduct, UpdatedBy string) (interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error occurred while marshaling data: %s", err)
	}
	query := `UPDATE product_data SET Products = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	res, err := db.Exec(query, jsonData, UpdatedBy, ID)
	if err != nil {
		return nil, fmt.Errorf("error occurred while data update: %s", err)
	}
	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0", "Products_data Updated successfully")
	if err != nil {
		return nil, fmt.Errorf("error occurred while data update: %s", err)
	}

	return res, nil
}

func UpdateAudit(Action string, ActionBy string, ID int, Ip string, Description string) error {
	query := `INSERT INTO audit_table (ActionStatus, ActionBy, ActionAt, ActionTableID, IpAddress, ActionDescription) 
			VALUES(?,?,NOW(),?,?,?)`

	_, err := db.Query(query, Action, ActionBy, ID, Ip, Description)
	if err != nil {
		return err
	}
	return nil
}

func UpdateComponentList(ID int, Component, Description, UpdatedBy string) error {
	query := `UPDATE components_list SET Component = ?, Description = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Query(query, Component, Description, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0.1", "components_list data updated successfully")
	if err != nil {
		return err
	}
	return nil
}

func DeleteComponentList(ID int, UpdatedBy string) error {
	query := `UPDATE components_list SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Query(query, 1, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0.1", "components_list data updated successfully")
	if err != nil {
		return err
	}
	return nil
}

func RemoveProduct(ID int, DeletedBy, Reason string) error {
	query := `UPDATE product_data SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now(), DeletedReason = ? WHERE ID =?`
	_, err := db.Query(query, 1, DeletedBy, Reason, ID)
	if err != nil {
		return err
	}

	err = UpdateAudit("Delete", DeletedBy, ID, "127.0.0.1", "Products_data Deleted successfully")
	if err != nil {
		return err
	}
	return nil
}

func ProductReference(ReferenceName string) (ProductDetailedData, error) {
	var data ProductDetailedData
	query := `SELECT ID, RefarenceName, Products FROM product_data WHERE RefarenceName = ?`
	err := db.QueryRow(query, ReferenceName).Scan(&data.ID, &data.ReferenceName, &data.Products)
	if err != nil {
		return data, err
	}
	// fmt.Println(data.Products)
	return data, nil
}

func FetchComponentDetails(productIDs []int) ([]ComponentDetail, error) {
	var componentDetails []ComponentDetail

	placeholders := make([]string, len(productIDs))
	for i := range productIDs {
		placeholders[i] = "?"
	}

	query := `SELECT ID, VarientName, VarientDescription, VarientPrice, VarientType 
              FROM components 
              WHERE ID IN (` + strings.Join(placeholders, ",") + `)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return componentDetails, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(intToInterfaceSlice(productIDs)...)
	if err != nil {
		return componentDetails, err
	}
	defer rows.Close()

	re := regexp.MustCompile(`\((.*?)\)`)

	for rows.Next() {
		var component ComponentDetail
		if err := rows.Scan(&component.ID, &component.VarientName, &component.VarientDescription, &component.VarientPrice, &component.VarientType); err != nil {
			return componentDetails, err
		}

		fmt.Printf("Original VarientType: %s\n", component.VarientType)

		matches := re.FindStringSubmatch(component.VarientType)
		if len(matches) > 1 {
			component.VarientType = matches[1] // Update VarientType to just the text inside parentheses
		} else {
			component.VarientType = "" // Set to empty if no match found
		}

		//fmt.Printf("Modified VarientType: %s\n", component.VarientType)

		componentDetails = append(componentDetails, component)
	}

	return componentDetails, nil
}

func intToInterfaceSlice(intSlice []int) []interface{} {
	interfaceSlice := make([]interface{}, len(intSlice))
	for i, v := range intSlice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}
