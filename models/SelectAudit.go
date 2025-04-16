package models

import (
	"encoding/json"
	"fmt"
)

func UserAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%user%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error scanning data %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func ComponentAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Components %' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func ComponentListAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Components_list%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func ProductsAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Products_data%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}
func ProductsListAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Product list%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func CompaniesAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Company%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func ContactAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%contact%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func DealsAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Deal %' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func MarginsAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Margin%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func PartnersAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Partner%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func PartnerContctAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Contact Person%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}
func GeMAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%GEM%' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func ProductsMainListAudit() ([]AuditData, error) {
	query := `SELECT ID, ActionStatus, ActionBy, ActionAt, ActionTableID, ActionDescription FROM audit_table WHERE ActionDescription LIKE '%Product_main_list %' ORDER BY ActionAt DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditData

	for rows.Next() {
		var items AuditData
		err := rows.Scan(&items.ID, &items.ActionStatus, &items.ActionBy, &items.ActionAt, &items.ActionTableID, &items.Description)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

// Fetch individual detailed data's
func GetUserDataQuery(ID int) ([]AuditUser, error) {
	query := `SELECT Name, Email, Mobile, Role, Designation, CreatedBy,CreatedAt,UpdatedBy,UpdatedAt FROM users WHERE ID = ? `
	rows, err := db.Query(query, ID)
	if err != nil {
		//
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []AuditUser

	for rows.Next() {
		var items AuditUser
		err := rows.Scan(&items.UserName, &items.Email, &items.Mobile, &items.UserRole, &items.Designation, &items.CreatedBy, &items.CreatedAt, &items.UpdatedBy, &items.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func GetProductDataQuery(ID int) ([]PricingProduct, error) {
	var products []PricingProduct

	query := `SELECT Products, RefarenceName,CreatedAt,CreatedBy FROM product_data WHERE ID = ? `
	rows, err := db.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()

	for rows.Next() {
		var jsonData []byte
		var referenceName string
		var Time string
		var CreatedBy string
		var product PricingProduct

		if err := rows.Scan(&jsonData, &referenceName, &Time, &CreatedBy); err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}

		if err := json.Unmarshal(jsonData, &product); err != nil {
			return nil, fmt.Errorf("error occurred %s", err)
		}

		product.ReferenceName = referenceName
		product.CreatedAt = Time
		product.CreatedBy = CreatedBy

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred %s", err)
	}

	return products, nil
}

func AuditComponentsList(ID int) ([]Component, error) {

	query := `SELECT Component, Description, ProductType,CreatedBy, CreatedAt FROM components_list WHERE ID = ?`
	rows, err := db.Query(query, ID)
	if err != nil {
		//
		return nil, fmt.Errorf("error occurred %s", err)
	}
	defer rows.Close()
	var data []Component

	for rows.Next() {
		var items Component
		err := rows.Scan(&items.ComponentName, &items.Description, &items.ProductType, &items.CreatedBy, &items.CreatedAt)

		if err != nil {
			return nil, fmt.Errorf("error scanning data %s", err)
		}
		data = append(data, items)
	}
	return data, nil
}

func AuditComponentsData(ID int) ([]VariantItems, error) {

	query := `SELECT VarientName, VarientPrice, VarientDescription,VarientBrand, CreatedBy, CreatedAt, UpdatedBy, UpdatedAt FROM components WHERE ID = ?`
	rows, err := db.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer rows.Close()

	var components []VariantItems

	for rows.Next() {
		var item VariantItems
		err := rows.Scan(&item.VariantName, &item.VariantPrice, &item.VariantDescription, &item.VarientBrand, &item.CreatedBy, &item.CreatedAt, &item.UpdatedBy, &item.UpdatedAt)
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

func AuditProductMainData(ID int) ([]ProductMainList, error) {

	query := `SELECT Product, CreatedBy, CreatedAt FROM product_main_list WHERE ID = ?`
	rows, err := db.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer rows.Close()

	var components []ProductMainList

	for rows.Next() {
		var item ProductMainList
		err := rows.Scan(&item.Product, &item.CreatedBy, &item.CreatedAt)
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

func AuditProductData(ID int) ([]GetProductList, error) {

	query := `SELECT ProductName, Description, CreatedBy, CreatedAt FROM product_list WHERE ID = ?`
	rows, err := db.Query(query, ID)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer rows.Close()

	var components []GetProductList

	for rows.Next() {
		var item GetProductList
		err := rows.Scan(&item.ProductName, &item.Description, &item.CreatedBy, &item.CreatedAt)
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

func CompaniesAuditData(ID int) (Company, error) {
	var company Company
	query := `SELECT ID, CompanyName, CompanyDomain, CompanyEmail, CompanyMobile, IndustryType, City, CreatedBy, CreatedAt FROM company WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(&company.ID, &company.CompanyName, &company.CompanyDomain, &company.Email, &company.Phone, &company.IndustryType, &company.City, &company.CreatedBy, &company.CreatedAt)
	if err != nil {
		return company, fmt.Errorf("error occurred %s", err)
	}
	return company, nil
}

func ContactAuditData(ID int) (Contacts, error) {
	var contacts Contacts
	query := `SELECT ID, Name, Email, CompanyID, Mobile	, LifeCycleStage, CreatedBy, CreatedAt FROM contacts WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(&contacts.ID, &contacts.Name, &contacts.Email,&contacts.CompanyID, &contacts.Mobile, &contacts.LifeCycleStage, &contacts.CreatedBy, &contacts.CreatedAt)
	if err != nil {
		return contacts, fmt.Errorf("error occurred %s", err)
	}
	var Company string
	if contacts.CompanyID != "0" {
		query := `SELECT CompanyName FROM company WHERE ID = ?`
		err := db.QueryRow(query, contacts.CompanyID).Scan(&Company)
		if err != nil {
			return contacts, fmt.Errorf("error occurred %s", err)
		}
		contacts.CompanyName = Company
	} else {
		contacts.CompanyName = "Individual"
	}
	return contacts, nil
}

func DealAuditData(ID int) (DealsStruct, error) {
	var deal DealsStruct
	query := `SELECT ID, ProductType, ProductDetails, CompanyID, IsCompany, DealOwner, DealApprover, QuotationUpdater, CreatedBy, CreatedAt FROM Deals WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(&deal.ID, &deal.ProductType, &deal.ProductDetails, &deal.CompanyID, &deal.IsCompany, &deal.DealOwner, &deal.DealApprover, &deal.QuotationUpdaterName, &deal.CreatedBy, &deal.CreatedAt)
	if err != nil {
		return deal, fmt.Errorf("error occurred %s", err)
	}
	var Company string
	if deal.IsCompany == 1 {
		query := `SELECT CompanyName FROM company WHERE ID = ?`
		err := db.QueryRow(query, deal.CompanyID).Scan(&Company)
		if err != nil {
			return deal, fmt.Errorf("error occurred %s", err)
		}
		deal.CompanyName = Company
	} else if deal.IsCompany == 0 {
		query := `SELECT Name FROM contacts WHERE ID = ?`
		err := db.QueryRow(query, deal.CompanyID).Scan(&Company)
		if err != nil {
			return deal, fmt.Errorf("error occurred %s", err)
		}
		deal.CompanyName = Company
	}
	return deal, nil
}
