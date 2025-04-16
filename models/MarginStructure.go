package models

import "database/sql"

type Margin struct {
	ID           int    `json:"id"`
	MarginType   string `json:"marginType"`
	BPPercentage int    `json:"bpPercentage"`
	CreatedBy    string `json:"createdBy"`
	CreatedAt    string `json:"createdAt"`
	UpdatedBy    sql.NullString
	UpdatedAt    sql.NullString
	IsActive     int `json:"isActive"`
}

type MarginValues struct {
	ID             int            `json:"id"`
	MarginID       int            `json:"marginID"`
	Tier           string         `json:"tier"`
	TierPercentage int            `json:"tierPercentage"`
	CreatedBy      string         `json:"createdBy"`
	CreatedAt      string         `json:"createdAt"`
	UpdatedBy      sql.NullString `json:"updatedBy"`
	UpdatedAt      sql.NullString `json:"updatedAt"`
	IsActive       int            `json:"isActive"`
	MarginName     string         `json:"MarginName"`
}

type DealsStruct struct {
	ID                    int            `json:"id"`
	Team                  string         `json:"team"`
	Email                 string         `json:"email"`
	Role                  string         `json:"role"`
	IsCompany             int            `json:"isCompany"`
	CompanyID             int            `json:"companyID"`
	ProductType           string         `json:"productType"`
	ProductDetails        string         `json:"productDescription"`
	Price                 string         `json:"pricing"`
	DealOwner             string         `json:"dealOwnerName"`
	DealOwnerID           int            `json:"dealOwnerId"`
	DealApprover          string         `json:"dealApprover"`
	DealApproverID        int            `json:"dealApproverID"`
	Quantity              string         `json:"quantity"`
	DealsStartDate        string         `json:"dealStartDate"`
	ApproximateEndDate    string         `json:"approximateClosingDate"`
	DealEndDate           sql.NullString `json:"dealEndDate"`
	StatusCommand         string         `json:"statusCommand"`
	Status                string         `json:"status"`
	Priority              string         `json:"priority"`
	IsApproved            int            `json:"isApproved"`
	CreatedBy             string         `json:"createdBy"`
	CreatedAt             string         `json:"createdAt"`
	UpdatedBy             sql.NullString `json:"updatedBy"`
	UpdatedAt             sql.NullString `json:"updatedAt"`
	IsDeleted             int            `json:"isDeleted"`
	ApprovedAt            sql.NullString `json:"approvedAt"`
	CompanyName           string         `json:"companyName"`
	ContactID             int            `json:"contactID"`
	Quotation             []string       `json:"quotation"`
	QuotationUpdaterID    int            `json:"quotationUpdaterID"`
	QuotationUpdaterName  string         `json:"quotationUpdaterName"`
	EditedBy              string         `json:"editedBy"`
	FinalOutCome          string         `json:"finalOutCome"`
	From                  string         `json:"fromEmail"`
	To                    []string       `json:"toMail"`
	CC                    []string       `json:"cc"`
	BCC                   []string       `json:"bcc"`
	Subject               string         `json:"subject"`
	Body                  string         `json:"body"`
	NotificationID        int            `json:"notificationID"`
	NotificationMessage   string         `json:"notificationMessage"`
	IsNotificationRead    int            `json:"isNotificationRead"`
	NotificationCreatedAt string         `json:"notificationCreatedAt"`
	UserID                int            `json:"userId"`
	MyDeals               bool           `json:"myDeals"`
	Reason                string         `json:"reason"`
	MailStatus            string         `json:"mailStatus"`
	Mobile                string         `json:"mobile"`
	Address               string         `json:"address"`
	TicketID              string         `json:"ticketID"`
	MarginType            string         `json:"marginType"`
	MarginTypePercentage  int            `json:"marginTypePercentage"`
	Tier                  string         `json:"tier"`
	ModelName             string         `json:"modelName"`
	DealOwnerEmail        string         `json:"dealOwnerEmail"`
}

type DenyDeals struct {
	ID        int    `json:"id"`
	Reason    string `json:"reason"`
	DeletedBy string `json:"deletedBy"`
}

type Demo struct {
	ID            int    `json:"id"`
	SalesPerson   string `json:"salesPerson"`
	EstimatedDate string `json:"estimatedDate"`
	Processor     string `json:"processor"`
	RAM           string `json:"ram"`
	Storage       string `json:"storage"`
	ProductType   string `json:"productType"`
	Status        string `json:"status"`
	Comments      string `json:"comments"`
	Notes         string `json:"notes"`
	SalesPersonID int    `json:"salesPersonID"`
	CreatedAt     string `json:"createdAt"`
	UpdatedAt     string `json:"updatedAt"`
	UpdatedBy     string `json:"updatedBy"`
	Team          string `json:"team"`
	Role          string `json:"role"`
	UserID        int    `json:"userID"`
}

type DemoSummaryStruct struct {
	ID          int    `json:"ID"`
	DemoID      int    `json:"demoId"`
	Description string `json:"description"`
	ByPerson    string `json:"byPerson"`
	HappenedAT  string `json:"happenedAt"`
}

type DealSummaryStruct struct {
	ID          int    `json:"id"`
	DealID      int    `json:"dealID"`
	Description string `json:"description"`
	ByPerson    string `json:"byPerson"`
}
