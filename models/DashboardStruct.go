package models

import "database/sql"

type ProgressStruct struct {
	Month      string `json:"month"`
	Win        sql.NullString
	Loss       sql.NullString
	OnProgress sql.NullString
}

type RecentDealsStruct struct {
	ID                   int    `json:"id"`
	IsCompany            int    `json:"isCompany"`
	CompanyID            int    `json:"companyID"`
	Name                 string `json:"name"`
	DealID               int    `json:"dealID"`
	OwnerName            string `json:"ownerName"`
	OwnerID              int    `json:"ownerID"`
	ApproverName         string `json:"approverName"`
	ApproverID           int    `json:"approverID"`
	QuotationUpdaterName string `json:"quotationUpdaterName"`
	QuotationUpdaterID   int    `json:"QuotationUpdaterID"`
	ProductType          string `json:"ProductType"`
	Price                string `json:"Price"`
	Quantity             string `json:"quantity"`
	Status               string `json:"Status"`
}
