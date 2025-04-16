package models

import "database/sql"

type PartnerStruct struct {
	ID                    int          `json:"id"`
	PartnerName           string       `json:"partner"`
	AuthCode              string       `json:"authCode"`
	AreaCovered           string       `json:"areaCovered"`
	CreatedBy             string       `json:"createdBy"`
	CreatedAt             string       `json:"createdAt"`
	UpdatedBy             string       `json:"updatedBy"`
	UpdatedAt             string       `json:"updatedAt"`
	Address               string       `json:"address"`
	Contacts              []GemContact `json:"contacts"`
	TurnOver              string       `json:"turnOver"`
	DueDate               string       `json:"dueDate"`
	RespectiveSalesPerson string       `json:"respectiveSalesPerson"`
	RegisteredDate        string       `json:"registeredDate"`
	Email                 string       `json:"email"`
	Phone                 string       `json:"phoneNumber"`
}

type GemContact struct {
	ContactID int    `json:"contactId"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
}

type ContactPerson struct {
	ID           int    `json:"id"`
	PartnerID    int    `json:"partnerId"`
	PartnerName  string `json:"partnerName"`
	ContactName  string `json:"contactPersonName"`
	ContactEmail string `json:"email"`
	ContactPhone string `json:"mobile"`
	CreatedBy    string `json:"createdBy"`
	CreatedAt    string `json:"createdAt"`
}

type GemStructure struct {
	ID                 int            `json:"id"`
	GemNumber          string         `json:"gemNumber"`
	DueDate            string         `json:"dueDate"`
	Quantity           int            `json:"quantity"`
	DealerName         string         `json:"dealerName"`
	DealerAddress      string         `json:"dealerAddress"`
	DealerEmail        string         `json:"dealerEmail"`
	DealerMobile       string         `json:"dealerMobile"`
	DealerID           int            `json:"dealerId"`
	Customer           string         `json:"customerName"`
	Description        string         `json:"description"`
	GemCertificate     []string       `json:"gemCertificate"`
	Quotation          []string       `json:"quotations"`
	GemRaisedBy        string         `json:"createdBy"`
	GemRaisedAt        string         `json:"gemRaisedAt"`
	QuotationUpdatedBy sql.NullString `json:"quotationUpdatedBy"`
	QuotationUpdaterID int            `json:"quotationUpdaterID"`
	QuotationUpdatedAt string         `json:"quotationUpdatedAt"`
	UpdatedBy          string         `json:"updatedBy"`
	UpdatedAt          string         `json:"updatedAt"`
	Status             string         `json:"status"`
	Link               string         `json:"link"`
	Team               string         `json:"team"`
	Level              string         `json:"level"`
	IsApproved         int            `json:"isApproved"`
	IsQuotationUpdated int            `json:"isQuotationUpdated"`
	ApprovedBy         string         `json:"approvedBy"`
	OwnerID            int            `json:"ownerId"`
	PricingID          int            `json:"pricingID"`
	ApproverID         int            `json:"approverID"`
	GemUrl             string         `json:"gemURL"`
	GemDate            string         `json:"GEMDate"`
	ProductType        string         `json:"productType"`
	TicketNo           string         `json:"ticketNo"`
	GEMFinalURL        string         `json:"gemFinalURL"`
	DeaOwnerEmail      string         `json:"dealOwnerEmail"`
}

type GeMApprove struct {
	ID                  int    `json:"id"`
	GemApprover         string `json:"gemApprover"`
	GemApproverID       int    `json:"gemApproverID"`
	QuotationApproverID int    `json:"quotationUpdaterID"`
	QuotationApprover   string `json:"quotationUpdaterName"`
}

type GeMDeny struct {
	ID           int    `json:"id"`
	GemRemover   string `json:"deletedBy"`
	GemRemoverID int    `json:"deletedByID"`
	Reason       string `json:"reason"`
}

type AllGeM struct {
	UserID int    `json:"userID"`
	Team   string `json:"team"`
	Role   string `json:"role"`
}

type UpdateStatus struct {
	ID        int    `json:"id"`
	Status    string `json:"status"`
	UpdatedBy string `json:"editedBy"`
}

type MailStruct struct {
	From    string   `json:"fromEmail"`
	To      []string `json:"toMail"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
}

type DetailedPartnerData struct {
	PartnerID             int              `json:"partnerId"`
	PartnerName           string           `json:"partner"`
	AuthCode              string           `json:"authCode"`
	AreaCovered           string           `json:"areaCovered"`
	PartnerCreatedBy      string           `json:"partnerCreatedBy"`
	PartnerCreatedAt      string           `json:"partnerCreatedAt"`
	Address               string           `json:"address"`
	Contacts              []GemContact     `json:"contacts"`
	TurnOver              string           `json:"turnOver"`
	DueDate               string           `json:"dueDate"`
	RespectiveSalesPerson string           `json:"respectiveSalesPerson"`
	RegisteredDate        string           `json:"registeredDate"`
	Email                 string           `json:"email"`
	Phone                 string           `json:"phoneNumber"`
	LastActivityDate      string           `json:"lastActivityData"`
	Gem                   []GetProductList `json:"gem"`
}
