package models

import "database/sql"

type Contacts struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Email            string `json:"email"`
	CompanyID        string `json:"companyID"`
	OwnerID          int    `json:"ownerId"`
	Mobile           string `json:"phoneNumber"`
	LastContacted    sql.NullString
	LifeCycleStage   string `json:"lifeCycleStage"`
	LeadStatus       string `json:"leadStatus"`
	LegalBasis       string `json:"legalBasis"`
	ContactLevelType string `json:"contactLevelType"`
	CreatedBy        string `json:"createdBy"`
	CreatedAt        string `json:"createdAt"`
	UpdatedBy        sql.NullString
	UpdatedAt        sql.NullString
	JobRole          string `json:"jobTitle"`
	EditedBY         string `json:"editedBy"`
	ContactLogo      string `json:"contactLogo"`
	CompanyName      string `json:"companyName"`
	DeletedBy        string `json:"deletedBy"`
	IsCompany        int    `json:"isCompany"`
	Company          int    `json:"company"`
	LeadSource       string `json:"leadSource"`
	Designation      string `json:"designation"`
}
