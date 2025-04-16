package models

import "database/sql"

type Company struct {
	ID                 int    `json:"id"`
	CompanyName        string `json:"companyName"`
	OwnerID            int    `json:"ownerId"`
	Email              string `json:"email"`
	CompanyDomain      string `json:"domainName"`
	IndustryType       string `json:"industry"`
	Phone              string `json:"mobile"`
	CompanyDescription string `json:"description"`
	City               string `json:"city"`
	State              string `json:"state"`
	Country            string `json:"country"`
	PostalCode         string `json:"postalCode"`
	NumberOfEmployee   string `json:"numberOfEmployees"`
	NumberOfBranches   string `json:"numberOfBranches"`
	AnnualRevenue      string `json:"annualRevenue"`
	CompanyGST         string `json:"gstNumber"`
	CompanyPAN         string `json:"panNumber"`
	MSME               string `json:"msmeNumber"`
	CIN                string `json:"cin"`
	SocialMedia        string `json:"socialUrl"`
	CompanyOwner       string `json:"companyOwner"`
	CompanyLogo        string `json:"companyLogo"`
	TimeZone           string `json:"timeZone"`
	CreatedBy          string `json:"createdBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          sql.NullString
	EditedBY           string `json:"editedBy"`
	CompanyLevel       string `json:"companyLevel"`
	Contacts          []Contacts `json:"contacts"` 
}

type CompanyStruct struct {
	ID   int    `json:"id"`
	Team string `json:"team"`
	Role string `json:"role"`
}

type CompanyName struct {
	ID          int    `json:"id"`
	CompanyName string `json:"companyName"`
}
