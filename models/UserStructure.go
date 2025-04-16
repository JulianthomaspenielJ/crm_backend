package models

import "database/sql"

type User struct {
	ID           int
	IsApproved   int
	IsDeleted    int
	UserName     sql.NullString
	UserRole     sql.NullString
	Mobile       sql.NullString
	Designation  sql.NullString
	Team         sql.NullString
	AvatarURL    sql.NullString
	SignatureURL sql.NullString
}

type UpdateUserData struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobileNo"`
	Role        string `json:"role"`
	Team        string `json:"team"`
	Designation string `json:"designation"`
	CreatedBy   string `json:"approvedBy"`
}

type UserDetails struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	Role        string `json:"role"`
	Designation string `json:"designation"`
	AvatarURL   string `json:"AvatarURL"`
	Team        string `json:"team"`
	CreatedBy   string `json:"approvedBy"`
	CreatedAt   string `json:"createdAt"`
	UpdatedBy   string `json:"updatedBy"`
	EditedBy    string `json:"editedBy"`
}

type AuditStructure struct {
	Audit string `json:"newValue"`
}

type AuditData struct {
	ID            int    `json:"id"`
	ActionStatus  string `json:"actionStatus"`
	ActionBy      string `json:"actionBy"`
	ActionAt      string `json:"actionAt"`
	ActionTableID int    `json:"actionTableID"`
	Description   string `json:"description"`
}

type AuditDetailsData struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type AuditUser struct {
	UserName    string `json:"name"`
	UserRole    string `json:"role"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	Designation string `json:"designation"`
	CreatedAt   string `json:"createdAt"`
	CreatedBy   string `json:"createdBy"`
	UpdatedBy   sql.NullString
	UpdatedAt   sql.NullString
}

type TeamStructure struct {
	ID       int    `json:"id"`
	TeamName string `json:"teamName"`
}

type AdminTable struct {
	ID            int    `json:"id"`
	UserID        int    `json:"userID"`
	UserRole      string `json:"userRole"`
	AllocatedTeam string `json:"allocatedTeam"`
	UnderBy       int    `json:"underBy"`
	CreatedBy     string `json:"createdBy"`
	CreatedAt     string `json:"createdAt"`
	UpdatedBy     sql.NullString
	UpdatedAt     sql.NullString
	IsActive      int `json:"isActive"`
}
