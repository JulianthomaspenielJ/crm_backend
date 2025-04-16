package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// fetch Component list data
func GetComponents() ([]Component, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query("SELECT ID, Component, Description, ProductType FROM components_list WHERE IsDeleted = 0 ORDER BY CreatedAt DESC")
	if err != nil {
		return nil, fmt.Errorf("error querying componentlist table: %v", err)
	}
	defer rows.Close()

	var components []Component

	for rows.Next() {
		var comp Component
		if err := rows.Scan(&comp.ID, &comp.ComponentName, &comp.Description, &comp.ProductType); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		components = append(components, comp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return components, nil
}

func GetUserByEmail(email string) (*User, error) {
	query := `SELECT ID, IsApproved, IsDeleted, Name, Role, Mobile, Designation, AvatarURL, Team, SignatureURL FROM users WHERE Email = ?`
	row := db.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.IsApproved, &user.IsDeleted, &user.UserName, &user.UserRole, &user.Mobile, &user.Designation, &user.AvatarURL, &user.Team, &user.SignatureURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found in the database")
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if user.AvatarURL.Valid && user.AvatarURL.String != "" {
		base64Avatar, err := imageToBase64("./Assets/AvatarImages/" + user.AvatarURL.String)
		if err != nil {
			return nil, fmt.Errorf("failed to convert avatar image to base64: %v", err)
		}
		user.AvatarURL.String = base64Avatar
	}

	if user.SignatureURL.Valid && user.SignatureURL.String != "" {
		base64Signature, err := imageToBase64("./Assets/SignatureImages/" + user.SignatureURL.String)
		if err != nil {
			return nil, fmt.Errorf("failed to convert signature image to base64: %v", err)
		}
		user.SignatureURL.String = base64Signature
	}

	return &user, nil
}

func GetComponentsData(ID int) ([]VariantItems, error) {

	query := `SELECT ID, VarientName, VarientPrice, VarientDescription FROM components WHERE ComponentID = ?  AND IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, ID, 0)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer rows.Close()

	var components []VariantItems

	for rows.Next() {
		var item VariantItems
		err := rows.Scan(&item.ID, &item.VariantName, &item.VariantPrice, &item.VariantDescription)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		components = append(components, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return components, nil
}

func GetComponentslist() ([]VariantItems, error) {
	query := `SELECT ID,VarientName, VarientDescription, VarientPrice,VarientBrand FROM components WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return nil, fmt.Errorf("error occurred while executing query : %v", err)
	}

	defer rows.Close()

	var data []VariantItems

	for rows.Next() {
		var items VariantItems
		err := rows.Scan(&items.ID, &items.VariantName, &items.VariantDescription, &items.VariantPrice, &items.VarientBrand)
		if err != nil {
			return nil, fmt.Errorf("error scanning data %s ", err)
		}
		data = append(data, items)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating data row : %v", err)
	}
	return data, nil
}

func GetProductLists() ([]GetProductList, error) {
	query := `SELECT ID, ProductName, Description FROM product_list WHERE IsDeleted = ? AND ProductID = 0 ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return nil, fmt.Errorf("error occurred while query executing %s", err)
	}
	defer rows.Close()

	var data []GetProductList

	for rows.Next() {
		var items GetProductList
		err := rows.Scan(&items.ID, &items.ProductName, &items.Description)
		if err != nil {
			return nil, fmt.Errorf("error occurred while executing %s", err)
		}
		data = append(data, items)
	}

	return data, nil
}

func GetProductData() ([]PricingProduct, error) {
	var products []PricingProduct
	query := `
		SELECT 
			pd.ID, 
			pd.CreatedBy, 
			pd.CreatedAt, 
			pd.RefarenceName, 
			pd.IsDeal, 
			d.IsCompany, 
			COALESCE(
				CASE 
					WHEN pd.IsDeal = 2 THEN g.DealerName
					ELSE COALESCE(c.CompanyName, ct.Name)
				END, 'Unknown') AS CustomerName,
			d.TicketNo AS DealTicketNo,       -- Add TicketNo from Deals
			g.TicketNo AS GeMTicketNo         -- Add TicketNo from GeM
		FROM 
			product_data pd
		LEFT JOIN Deals d ON pd.DealID = d.ID AND pd.IsDeal = 1
		LEFT JOIN company c ON d.IsCompany = 1 AND d.CompanyID = c.ID
		LEFT JOIN contacts ct ON d.IsCompany = 0 AND d.CompanyID = ct.ID
		LEFT JOIN GeM g ON pd.DealID = g.ID AND pd.IsDeal = 2
		WHERE 
			pd.IsDeleted = 0 
		ORDER BY 
			pd.CreatedAt DESC;

	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving data: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var createdAtStr string
		var ID int
		var createdBy, referenceName, customerName, dealTicketNo, gemTicketNo sql.NullString
		var IsDeal, IsCompany sql.NullInt64
		var product PricingProduct

		// Include dealTicketNo and gemTicketNo in the scan
		if err := rows.Scan(&ID, &createdBy, &createdAtStr, &referenceName, &IsDeal, &IsCompany, &customerName, &dealTicketNo, &gemTicketNo); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %s", err)
		}

		createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("error occurred while parsing CreatedAt: %s", err)
		}

		product.CreatedBy = createdBy.String
		product.CreatedAt = createdAt.Format("02 Jan 2006 ")
		product.ReferenceName = referenceName.String
		product.ID = ID

		if IsDeal.Valid {
			product.IsDeal = IsDeal
		} else {
			product.IsDeal = sql.NullInt64{Int64: 0, Valid: false}
		}

		if customerName.Valid {
			product.CustomerName = customerName.String
		} else {
			product.CustomerName = ""
		}

		if dealTicketNo.Valid {
			product.TicketNo = dealTicketNo.String
		} else if gemTicketNo.Valid {
			product.TicketNo = gemTicketNo.String
		} else {
			product.TicketNo = ""
		}

		// Append product to the list
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %s", err)
	}

	return products, nil
}

func GetProductListData(id int) ([]PricingProduct, error) {
	var products []PricingProduct

	query := `SELECT Products, RefarenceName, CreatedAt FROM product_data WHERE IsDeleted = 0 AND ID = ?`
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("error occurred while retrieving data: %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var jsonData []byte
		var referenceName string
		var product PricingProduct
		var createdAt string

		if err := rows.Scan(&jsonData, &referenceName, &createdAt); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %v", err)
		}

		if err := json.Unmarshal(jsonData, &product); err != nil {
			return nil, fmt.Errorf("error occurred while unmarshaling data: %s", err)
		}

		product.ReferenceName = referenceName
		product.CreatedAt = createdAt

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %s", err)
	}

	return products, nil
}

func GetProductMainList() ([]ProductMainList, error) {
	query := `SELECT ID, Product FROM product_main_list WHERE IsDeleted = 0`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred while querying: %v", err)
	}
	defer rows.Close()

	var data []ProductMainList
	for rows.Next() {
		var items ProductMainList
		err := rows.Scan(&items.ID, &items.Product)
		if err != nil {
			return nil, fmt.Errorf("error occurred while Scanning: %v", err)
		}
		data = append(data, items)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during row iteration: %w", err)
	}

	return data, nil
}
func ProductIDData(ID int) ([]GetProductList, error) {
	query := `SELECT ID, ProductName, Description, ProductID FROM product_list WHERE ProductID = ? AND IsDeleted = ?`

	rows, err := db.Query(query, ID, 0)
	if err != nil {
		return nil, fmt.Errorf("error occurred while executing query: %w", err)
	}
	defer rows.Close()

	var data []GetProductList

	for rows.Next() {
		var items GetProductList
		err := rows.Scan(&items.ID, &items.ProductName, &items.Description, &items.ProductID)
		if err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %w", err)
		}
		data = append(data, items)
	}

	return data, nil
}

func GetDeletedData() ([]VariantItems, error) {

	query := `SELECT ID, VarientName, VarientPrice, VarientDescription FROM components WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 1)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer rows.Close()

	var components []VariantItems

	for rows.Next() {
		var item VariantItems
		err := rows.Scan(&item.ID, &item.VariantName, &item.VariantPrice, &item.VariantDescription)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		components = append(components, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return components, nil
}
func UserList() ([]UserDetails, error) {

	query := `SELECT ID, Name FROM users WHERE (Team = ? OR Team = ?) AND Role NOT LIKE '%Admin%' AND IsDeleted = ?`

	// Execute the query with the correct parameters
	rows, err := db.Query(query, "Marketing", "Sales", 0)
	if err != nil {
		return nil, fmt.Errorf("error occurred while executing query: %w", err)
	}
	defer rows.Close()

	var data []UserDetails

	for rows.Next() {
		var items UserDetails
		err := rows.Scan(&items.ID, &items.Name)
		if err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %w", err)
		}
		data = append(data, items)
	}

	return data, nil
}
