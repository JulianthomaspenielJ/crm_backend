package models

import (
	"database/sql"
	"fmt"
	"strings"
)

func InsertContact(Name, Email, Mobile, LifeCycleStage, LeadStatus, LegalBasis, ContactLevelType, CreatedBy, JobRole string, IsCompany, CompanyID, OwnerID int, Designation string) (string, error) {

	fmt.Print(CompanyID)
	InsertQuery := `INSERT INTO contacts(Name, Email, CompanyID, Mobile, LifeCycleStage, LeadStatus, ContactLevelType, CreatedBy, CreatedAt, IsActive, IsDeleted, OwnerID, JobRole ) 
					VALUES (?, ?, ?, ?, ?, ?, ?, ?,Now(), ?, ?, ?, ?)`
	res, err := db.Exec(InsertQuery, Name, Email, CompanyID, Mobile, LifeCycleStage, LeadStatus, ContactLevelType, CreatedBy, 1, 0, OwnerID, Designation)
	if err != nil {
		return "", fmt.Errorf("error occurred while inserting query: %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error occurred while fetching last insert id: %s", err)
	}

	err = AuditFunction("Create", CreatedBy, ID, "127.127.0.0", "New contact data created")
	if err != nil {
		return "", fmt.Errorf("error occurred while executing audit query: %s", err)
	}
	return "Data inserted successfully", nil
}

func GetContactTableData(ID int, Team, Role string) ([]Contacts, error) {
	var query string
	var rows *sql.Rows
	var err error

	if Team == "Admin" {
		query = `
			SELECT 
				ID, Name, Email, Mobile, CreatedBy, CompanyID, UpdatedAt, 
				ContactLevelType, CreatedAt 
			FROM 
				contacts 
			WHERE 
				IsDeleted = ? AND IsActive = ? AND CompanyID = ?
			ORDER BY 
				CreatedAt DESC`
		rows, err = db.Query(query, 0, 1, 0)
	} else {
		query = `
			SELECT 
				ID, Name, Email, Mobile, CreatedBy, CompanyID, UpdatedAt, 
				ContactLevelType, CreatedAt 
			FROM 
				contacts 
			WHERE 
				IsDeleted = ? AND IsActive = ? AND CompanyID = ?
				AND (OwnerID = ? OR OwnerID IN (SELECT ID FROM users WHERE Team = ? AND Role = ?))
			ORDER BY 
				CreatedAt DESC`
		rows, err = db.Query(query, 0, 1, 0, ID, Team, "Manager")
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while executing query: %s", err)
	}
	defer rows.Close()

	var contacts []Contacts

	for rows.Next() {
		var con Contacts
		if err := rows.Scan(&con.ID, &con.Name, &con.Email, &con.Mobile, &con.CreatedBy, &con.CompanyID, &con.UpdatedAt, &con.ContactLevelType, &con.CreatedAt); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %s", err)
		}
		contacts = append(contacts, con)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred while iterating rows: %s", err)
	}

	companyNames := make(map[string]string)
	companyIDs := make([]string, 0, len(contacts))

	for _, con := range contacts {
		if con.CompanyID != "0" && !contains(companyIDs, con.CompanyID) {
			companyIDs = append(companyIDs, con.CompanyID)
		}
	}

	if len(companyIDs) > 0 {
		placeholders := make([]string, len(companyIDs))
		for i := range placeholders {
			placeholders[i] = "?"
		}
		companyQuery := fmt.Sprintf("SELECT ID, CompanyName FROM company WHERE ID IN (%s)", strings.Join(placeholders, ","))

		rows, err := db.Query(companyQuery, interfaceSlice(companyIDs)...)
		if err != nil {
			return nil, fmt.Errorf("error occurred while fetching company names: %s", err)
		}
		defer rows.Close()

		for rows.Next() {
			var id, name string
			if err := rows.Scan(&id, &name); err != nil {
				return nil, fmt.Errorf("error occurred while scanning company data: %s", err)
			}
			companyNames[id] = name
		}
	}

	for i := range contacts {
		con := &contacts[i]
		if con.CompanyID != "0" {
			con.CompanyName = companyNames[con.CompanyID]
		} else {
			con.CompanyName = "Individual"
		}
	}

	return contacts, nil
}

func interfaceSlice(slice []string) []interface{} {
	intfSlice := make([]interface{}, len(slice))
	for i, v := range slice {
		intfSlice[i] = v
	}
	return intfSlice
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func GetContactData(ID int64) (Contacts, error) {
	var contact Contacts

	query := `SELECT ID, Name, Email, CompanyID, Mobile, LastContacted, LifeCycleStage, LeadStatus, ContactLevelType, CreatedBy, CreatedAt, UpdatedBy, JobRole FROM contacts WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(
		&contact.ID,
		&contact.Name,
		&contact.Email,
		&contact.CompanyID,
		&contact.Mobile,
		&contact.LastContacted,
		&contact.LifeCycleStage,
		&contact.LeadStatus,
		&contact.ContactLevelType,
		&contact.CreatedBy,
		&contact.CreatedAt,
		&contact.UpdatedBy,
		&contact.JobRole,
	)
	if contact.CompanyID != "0" {
		GetQuery := `SELECT CompanyName FROM company WHERE ID = ?`
		err := db.QueryRow(GetQuery, contact.CompanyID).Scan(&contact.CompanyName)
		if err != nil {
			return contact, fmt.Errorf("error occurred while query executing %s", err)
		}
		if err == sql.ErrNoRows {
			return contact, fmt.Errorf("no contact found with ID %s", err)
		}
	} else {
		contact.CompanyName = "Individual"
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return contact, fmt.Errorf("no contact found with ID %d", err)
		}
		return contact, fmt.Errorf("error occurred while executing query: %v", err)
	}

	return contact, nil
}

func EditContactData(ID int64, Name, Email, Phone, LifeCycleStage, LeadStatus, ContactLevelType, JobRole, UpdatedBy string) (Contacts, error) {
	var contacts Contacts
	query := `UPDATE contacts SET Name = ? , Email = ?, Mobile = ?, LifeCycleStage = ?, LeadStatus = ?, ContactLevelType = ?, JobRole = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, Name, Email, Phone, LifeCycleStage, LeadStatus, ContactLevelType, JobRole, UpdatedBy, ID)
	if err != nil {
		return contacts, fmt.Errorf("error occurred while query executing %s", err)
	}

	err = AuditFunction("Update", UpdatedBy, ID, "127.0.0.1", "contact data Updated successfully")
	if err != nil {
		return contacts, fmt.Errorf("error occurred while query executing %s", err)
	}

	res, err := GetContactData(ID)
	if err != nil {
		return contacts, fmt.Errorf("error occurred while getting data %s", err)
	}
	return res, nil
}

func DeleteContact(ID int64, IsCompany int, UpdatedBy string) error {
	var query string
	if IsCompany == 1 {
		query = `UPDATE company SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
		err := AuditFunction("Delete", UpdatedBy, ID, "127.0.0.1", "Company data Deleted successfully")
		if err != nil {
			return err
		}
	} else {
		query = `UPDATE contacts SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
		err := AuditFunction("Delete", UpdatedBy, ID, "127.0.0.1", "Contact data Deleted successfully")
		if err != nil {
			return err
		}
	}
	_, err := db.Exec(query, 1, UpdatedBy, ID)
	if err != nil {
		return err
	}

	return nil
}

func InsertCompanyContact(CompanyID, Name, Email, Mobile, Designation, CreatedBy string) error {
	InsertQuery := `INSERT INTO contacts(Name, Email, Mobile, JobRole, CompanyID, LeadStatus, ContactLevelType, CreatedBy, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, Now())`
	_, err := db.Exec(InsertQuery, Name, Email, Mobile, Designation, CompanyID,  "N/A", "N/A", CreatedBy)
	if err != nil {
		return err
	}
	return nil
}