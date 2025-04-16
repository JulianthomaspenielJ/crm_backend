package models

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

var genericDomains = []string{
	"gmail.com",
	"yahoo.com",
	"outlook.com",
	"hotmail.com",
	"icloud.com",
}

func isGenericDomain(domain string) bool {
	domain = strings.ToLower(domain)
	for _, generic := range genericDomains {
		if domain == generic {
			return true
		}
	}
	return false
}

func InsertCompany(CompanyName, Email, companyDomain, Phone, IndustryType, CompanyDescription, City, State, Country, PostalCode, NumberOfEmployee, AnnualRevenue, CompanyGST, CompanyPAN, MSME, NumberOfBranches, CIN, SocialMedia, CreatedBy, CompanyOwner, CompanyLogo, TimeZone string, OwnerId int) (string, int, error) {
	var ExistingID, PAN, GST, CEmail, PhoneNumber, CompanyLevel string

	if !isGenericDomain(companyDomain) {
		checkQuery := `SELECT ID FROM company WHERE CompanyDomain = ?`
		err := db.QueryRow(checkQuery, companyDomain).Scan(&ExistingID)
		if err != nil && err != sql.ErrNoRows {
			return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
		}

		if ExistingID != "" {
			return "domain", 0, nil
		}
	}

	checkMailQuery := `SELECT ID FROM company WHERE CompanyEmail = ?`
	err := db.QueryRow(checkMailQuery, Email).Scan(&CEmail)
	if err != nil && err != sql.ErrNoRows {
		return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
	}
	if CEmail != "" {
		res := "Email"
		return res, 0, nil
	}

	checkPhoneQuery := `SELECT ID FROM company WHERE CompanyMobile = ?`
	err = db.QueryRow(checkPhoneQuery, Phone).Scan(&PhoneNumber)
	if err != nil && err != sql.ErrNoRows {
		return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
	}
	if PhoneNumber != "" {
		res := "Phone"
		return res, 0, nil
	}

	if CompanyPAN != "N/A" {
		checkPanQuery := `SELECT ID FROM company WHERE CompanyPAN = ?`
		err = db.QueryRow(checkPanQuery, CompanyPAN).Scan(&PAN)
		if err != nil && err != sql.ErrNoRows {
			return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
		}

		if PAN != "" {
			return "PAN", 0, nil
		}
	}
	if CompanyGST != "N/A" {
		checkGSTQuery := `SELECT ID FROM company WHERE CompanyGST = ?`
		err = db.QueryRow(checkGSTQuery, CompanyGST).Scan(&GST)
		if err != nil && err != sql.ErrNoRows {

			return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
		}
		if GST != "" {
			res := "GST"
			return res, 0, nil
		}
	}

	if NumberOfEmployee == "1 to 10" {
		CompanyLevel = "Bronze"
	} else if NumberOfEmployee == "10 to 50" {
		CompanyLevel = "Bronze"
	} else if NumberOfEmployee == "50 to 100" {
		CompanyLevel = "Silver"
	} else {
		CompanyLevel = "Gold"
	}
	TimeZoneIn := "Kolkata,India"
	query := `INSERT INTO company(CompanyName, CompanyDomain, CompanyEmail, CompanyMobile, IndustryType, CompanyDescription, City, State, Country, PostelCode, NumberOfEmployee, AnnualRevenue, CompanyGST, CompanyPAN, MSME, NumberOfBranches, CIN, SocialMedia, CreatedBy, ComapnyOwner, CompanyLogo, TimeZone, CreatedAt, IsDeleted, OwnerID, CompanyLevel) 
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, Now(), ?, ?, ?)`
	result, err := db.Exec(query, CompanyName, companyDomain, Email, Phone, IndustryType, CompanyDescription, City, State, Country, PostalCode, NumberOfEmployee, AnnualRevenue, CompanyGST, CompanyPAN, MSME, NumberOfBranches, CIN, SocialMedia, CreatedBy, CompanyOwner, CompanyLogo, TimeZoneIn, 0, OwnerId, CompanyLevel)
	if err != nil {
		return "", 0, fmt.Errorf("error occurred while executing query: %s", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return "", 0, fmt.Errorf("error occurred while fetching last insert id: %s", err)
	}
	err = AuditFunction("Create", CreatedBy, id, "127.127.0.0", "New Company created")
	if err != nil {
		return "", 0, fmt.Errorf("error occurred while executing audit query: %s", err)
	}

	return "Data inserted successfully", int(id), nil
}

func GetCompanyTableData(ID int, Role, Team string) ([]Company, error) {
	var query string
	var rows *sql.Rows
	var err error

	if Team == "Admin" {
		query = `
		SELECT 
			ID, CompanyName, CompanyEmail, CreatedBy, CreatedAt, CompanyMobile, 
			UpdatedAt, Country, City, IndustryType, CompanyLevel , CompanyDomain
		FROM 
			company 
		WHERE 
			IsDeleted = 0
		ORDER BY CreatedAt DESC`

		rows, err = db.Query(query)
	} else {
		query = `
		SELECT 
			ID, CompanyName, CompanyEmail, CreatedBy, CreatedAt, CompanyMobile, 
			UpdatedAt, Country, City, IndustryType, CompanyLevel, CompanyDomain 
		FROM 
			company 
		WHERE 
			IsDeleted = 0 AND (OwnerID = ? OR OwnerID IN (SELECT ID FROM users WHERE Team = ? AND Role = 'Manager'))
		ORDER BY CreatedAt DESC`
		rows, err = db.Query(query, ID, Team)
	}

	if err != nil {
		return nil, fmt.Errorf("error occurred while query executing: %s", err)
	}
	defer rows.Close()

	var companies []Company

	for rows.Next() {
		var comp Company
		var createdAt []byte
		if err := rows.Scan(
			&comp.ID, &comp.CompanyName, &comp.Email, &comp.CreatedBy,
			&createdAt, &comp.Phone, &comp.UpdatedAt, &comp.Country,
			&comp.City, &comp.IndustryType, &comp.CompanyLevel, &comp.CompanyDomain,
		); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data: %s", err)
		}

		if len(createdAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
			if err != nil {
				return nil, fmt.Errorf("error parsing CreatedAt time: %s", err)
			}
			comp.CreatedAt = parsedTime.Format("02 Jan 2006 ")
		} else {
			comp.CreatedAt = ""
		}
		companies = append(companies, comp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return companies, nil
}

func GetCompanyData(ID int) (Company, error) {
	var company Company

	query := `SELECT ID, CompanyName, CompanyDomain, CompanyEmail, CompanyMobile, IndustryType, CompanyDescription, City, State, Country, PostelCode, NumberOfEmployee, AnnualRevenue, CompanyGST, CompanyPAN, MSME, NumberOfBranches, CIN, SocialMedia, CreatedBy, ComapnyOwner, CompanyLogo, TimeZone, CreatedAt FROM company WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(
		&company.ID,
		&company.CompanyName,
		&company.CompanyDomain,
		&company.Email,
		&company.Phone,
		&company.IndustryType,
		&company.CompanyDescription,
		&company.City,
		&company.State,
		&company.Country,
		&company.PostalCode,
		&company.NumberOfEmployee,
		&company.AnnualRevenue,
		&company.CompanyGST,
		&company.CompanyPAN,
		&company.MSME,
		&company.NumberOfBranches,
		&company.CIN,
		&company.SocialMedia,
		&company.CreatedBy,
		&company.CompanyOwner,
		&company.CompanyLogo,
		&company.TimeZone,
		&company.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {

			return company, fmt.Errorf("no company found with ID %d", ID)
		}

		return company, fmt.Errorf("error occurred while executing query: %v", err)
	}

	return company, nil
}

func EditCompanyData(ID int64, IndustryType, Email, City, State, Country, PostelCode, NumberOfEmployees, NumberOfBranches, AnnualRevenue, TimeZone, Description, SocialMedia, Mobile, UpdatedBy string) (Company, error) {
	var company Company
	query := `UPDATE company SET IndustryType = ?, CompanyEmail = ?, City = ? , State = ?, Country = ?, PostelCode = ?, NumberOfEmployee = ?, AnnualRevenue = ?, NumberOfBranches = ?, SocialMedia =? , UpdatedBy = ?, UpdatedAt = Now(), TimeZone = ?, CompanyMobile = ?, CompanyDescription = ? WHERE ID = ? `
	_, err := db.Exec(query, IndustryType, Email, City, State, Country, PostelCode, NumberOfEmployees, AnnualRevenue, NumberOfBranches, SocialMedia, UpdatedBy, TimeZone, Mobile, Description, ID)
	if err != nil {
		return company, fmt.Errorf("error occurred while query executing %s", err)
	}

	err = AuditFunction("Update", UpdatedBy, ID, "127.0.0.1", "Company  data Updated successfully")
	if err != nil {
		return company, err
	}
	result, err := GetCompanyData(int(ID))
	if err != nil {
		return company, err
	}
	return result, nil
}

func AuditFunction(Action string, ActionBy string, ID int64, IP string, Description string) error {
	query := `
		INSERT INTO audit_table (ActionStatus, ActionBy, ActionAt, ActionTableID, IpAddress, ActionDescription)
		VALUES (?, ?, Now(), ?, ?, ?)
	`

	result, err := db.Exec(query, Action, ActionBy, ID, IP, Description)
	if err != nil {
		return fmt.Errorf("failed to insert audit data: %v", err)
	}
	fmt.Print(result)
	return nil
}

func GetContactList(ID int) ([]Contacts, bool, error) {
	var ContactsData []Contacts

	fmt.Print("it was here", ID)
	query := `SELECT ID, Name, Email, Mobile, LeadStatus, ContactLevelType, JobRole, CreatedBy FROM contacts WHERE CompanyID = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, ID)
	if err != nil {
		return ContactsData, false, fmt.Errorf("error occurred while query executing %s", err)
	}
	defer rows.Close()
	for rows.Next() {
		var contact Contacts
		if err := rows.Scan(&contact.ID, &contact.Name, &contact.Email, &contact.Mobile, &contact.LeadStatus, &contact.ContactLevelType, &contact.Designation, &contact.CreatedBy); err != nil {
			return ContactsData, false, fmt.Errorf("error occurred while scanning data %s", err)
		}
		ContactsData = append(ContactsData, contact)
	}
	fmt.Print("it was done")
	return ContactsData, true, nil
}

func GetAllCompany() ([]CompanyName, error) {
	var companies []CompanyName
	query := `SELECT ID, CompanyName FROM company WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)

	if err != nil {
		return companies, fmt.Errorf("error occurred while query executing %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		var comp CompanyName
		if err := rows.Scan(
			&comp.ID, &comp.CompanyName); err != nil {
			return companies, fmt.Errorf("error occurred while scanning data %s", err)
		}
		companies = append(companies, comp)
	}
	return companies, nil
}
