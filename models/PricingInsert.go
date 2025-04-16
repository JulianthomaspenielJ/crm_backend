package models

import (
	"database/sql"
	"encoding/json"
	"fmt"

	//"testing/quick"
	"time"

	_ "github.com/go-sql-driver/mysql"
	//"github.com/pelletier/go-toml/query"
)

var currentTime = time.Now().Format("2006/01/02 15:04:05")

func InsertComponent(componentName, description, productType, createdBy string) error {
	var ExistComponentName string

	if db == nil {
		return fmt.Errorf("database connection is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %v", err)
	}

	stmt, err := tx.Prepare("INSERT INTO components_list (Component, Description, ProductType, CreatedBy, CreatedAt, IsDeleted) VALUES (?, ?, ?, ?, Now(), ?)")
	if err != nil {
		fmt.Print("ERR : ", err)
		tx.Rollback()
		return fmt.Errorf("error preparing statement for componentlist insert: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(componentName, description, productType, createdBy, 0)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing componentlist insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting last insert ID: %v", err)
	}

	res, err := InsertAudit("Create", createdBy, id, "127.1.1", "Data inserted into the components_list table")
	if err != nil {
		return nil
	}
	fmt.Print(res)

	productCheckQuery := `SELECT ProductName FROM product_list WHERE ProductName = ? AND ProductID = ? AND IsDeleted = ?`
	err = db.QueryRow(productCheckQuery, componentName, 0, 0).Scan(&ExistComponentName)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error occurred While checking select query %s", err)
	}

	if ExistComponentName == "" {
		resPro, err := tx.Exec(`INSERT INTO product_list(ProductName, Description, IsDeleted, CreatedBy, CreatedAt, ProductID) VALUES(?, ?, ?, ?, Now(), ?)`, componentName, "Components", 0, createdBy, 0)
		if err != nil {
			tx.Rollback()
			return err
		}
		mainProductID, err := resPro.LastInsertId()
		if err != nil {
			tx.Rollback()
			return err
		}

		auditRes, err := InsertAudit("Create", createdBy, mainProductID, "127.0.0", "Data inserted into product list table")
		if err != nil {
			return err
		}
		fmt.Print(auditRes)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

func InsertComponents(componentName string, componentType string, price float64, description, brand string, variantID string, createdBy string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("database connection is nil")
	}

	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("error beginning transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var existingID string

	checkQuery := `SELECT ComponentID FROM components WHERE VarientName = ? AND ComponentID = ?`
	err = tx.QueryRow(checkQuery, componentName, variantID).Scan(&existingID)

	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error checking existing component: %v", err)
	}

	if existingID != "" {
		res := "Component Name already exist"
		return res, nil
	}

	query := `INSERT INTO components (VarientName, VarientType, VarientPrice, VarientDescription, VarientBrand, ComponentID, IsDeleted, CreatedBy, CreatedAt)
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := tx.Exec(query, componentName, componentType, price, description, brand, variantID, 0, createdBy, currentTime)
	if err != nil {
		return "", fmt.Errorf("error executing components insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting last insert ID: %v", err)
	}

	res, err := InsertAudit("Create", createdBy, id, "127.1.0", "Data inserted to the components table")
	fmt.Print(res)

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("error committing transaction: %v", err)
	}

	return "", nil
}

func StoreProductData(data PricingProduct, AddressTo string, CreatedBy string, RefarenceName string, IsDeal, DealID int) (interface{}, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error occurred while marshaling data: %s", err)
	}

	// checkQuery := `SELECT RefarenceName FROM product_data WHERE RefarenceName = ? AND IsDeleted =?`
	// var existingRefName string
	// err = db.QueryRow(checkQuery, RefarenceName, 0).Scan(&existingRefName)
	// if err != nil && err != sql.ErrNoRows {
	// 	fmt.Println("Error occurred during query execution:", err)
	// 	return nil, fmt.Errorf("error ocurred : %s", err)
	// }

	// if existingRefName != "" {
	// 	fmt.Println("Reference name already exists.")
	// 	data := "Reference Name already exist"
	// 	return data, nil
	// }

	query := `INSERT INTO product_data(Products, CreatedBy, CreatedAt, IsDeleted, RefarenceName, TOAddress, IsDeal, DealID)
			  VALUES (?, ?, Now(), ?, ?, ?, ?, ?)`

	result, err := db.Exec(query, jsonData, CreatedBy, 0, RefarenceName, AddressTo, IsDeal, DealID)

	if err != nil {
		fmt.Print("error came frm here")
		return nil, fmt.Errorf("error occurred while data inserting: %s", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Print("error data affecting frm here")
		return nil, fmt.Errorf("no data affected: %s", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("error getting last insert id: %v", err)
	}

	auditquery, err := InsertAudit("Create", CreatedBy, id, "127.0.0.1", "Data Inserted to the Products_data table")
	if err != nil {
		return nil, fmt.Errorf("error occurred while inserting data to the Audit Table %s", err)
	}
	fmt.Print(rows, auditquery)

	return id, nil
}

func InsertProductMainList(product PricingMainListProduct) (interface{}, error) {
	fmt.Printf("Inserting Product: %+v\n", product)

	checkQuery := `SELECT Product FROM product_main_list WHERE Product = ?`
	var ExitProductName string
	err := db.QueryRow(checkQuery, product.ProductName).Scan(&ExitProductName)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("error occurred %s", err)
	}

	if ExitProductName != "" {
		fmt.Println("Reference name already exists.")
		data := "Reference Name already exist"
		return data, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
		}
	}()

	res, err := tx.Exec("INSERT INTO product_main_list(Product, CreatedBy, CreatedAt,IsDeleted) VALUES(?, ?, Now(),?)", product.ProductName, product.CreatedBy, 0)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	mainProductID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	auditRes, err := InsertAudit("Create", product.CreatedBy, mainProductID, "127.0.0", "Data inserted into product_main_list table")
	if err != nil {
		return nil, err
	}
	fmt.Print(auditRes)
	for _, component := range product.SelectedComponents {
		stmt, err := tx.Prepare("INSERT INTO product_list(ProductName, Description, IsDeleted, CreatedBy, CreatedAt, ProductID) VALUES(?, ?, ?, ?, Now(), ?)")
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		productRes, err := stmt.Exec(component.ProductName, component.Description, 0, product.CreatedBy, mainProductID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		productID, err := productRes.LastInsertId()
		if err != nil {
			return nil, err
		}
		auditResult, err := InsertAudit("Create", product.CreatedBy, productID, "127.0.0", "Data inserted into product list table")
		if err != nil {
			return nil, err
		}
		fmt.Print(auditResult)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	insertedData := struct {
		ProductName        string                 `json:"productName"`
		SelectedComponents []PricingMainComponent `json:"selectedComponents"`
		MainProductID      int64                  `json:"mainProductID"`
	}{
		ProductName:        product.ProductName,
		SelectedComponents: product.SelectedComponents,
		MainProductID:      mainProductID,
	}

	return insertedData, nil
}

func InsertAudit(Action string, ActionBy string, ID int64, IP string, Description string) (int64, error) {
	query := `
		INSERT INTO audit_table (ActionStatus, ActionBy, ActionAt, ActionTableID, IpAddress, ActionDescription)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := db.Exec(query, Action, ActionBy, currentTime, ID, IP, Description)
	if err != nil {
		return 0, fmt.Errorf("failed to insert audit data: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected in audit: %v", err)
	}

	return rowsAffected, nil
}
