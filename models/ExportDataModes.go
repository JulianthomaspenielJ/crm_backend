package models

import (
	"database/sql"
	"time"
)

type ExportStruct struct {
	ID       int    `json:"id"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

type ExportDataStruct struct {
	ID        int    `json:"id"`
	Priority  string `json:"priority"`
	CreatedAt string `json:"createdAt"`
	CreatedBy string `json:"createdBy"`
	Status    string `json:"status"`
	Customer  string `json:"company"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	//Location           string `json:"location"`
	QuotationUpdatedAt string `json:"quotationUpdatedAt"`
	QuotationUpdater   string `json:"quotationUpdater"`
	ProductType        string `json:"productType"`
	ProductDetails     string `json:"productDetails"`
	Quantity           int    `json:"quantity"`
	DeletedReason      string `json:"reason"`
}

func ExportData(FromDate, ToDate string) ([]ExportDataStruct, error) {
	var Datas []ExportDataStruct
	var Status string
	var IsCompany, CompanyID int
	var CreatedAt string
	var QuotationUpdatedAt, QuotationUpdater, DeletedReason sql.NullString
	query := `SELECT ID, Priority, CreatedAt, CreatedBy, Status, IsCompany, CompanyID, QuotationUpdatedAt, QuotationUpdater, ProductType, ProductDetails, Quantity, DeletedReason FROM Deals WHERE CreatedAt BETWEEN ? AND ? AND DealOwnerID  = ? `
	rows, err := db.Query(query, FromDate, ToDate, 12)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var Data ExportDataStruct
		err := rows.Scan(&Data.ID, &Data.Priority, &CreatedAt, &Data.CreatedBy, &Status, &IsCompany, &CompanyID, &QuotationUpdatedAt, &QuotationUpdater, &Data.ProductType, &Data.ProductDetails, &Data.Quantity, &DeletedReason)
		if err != nil {
			return nil, err
		}
		parsedTime, err := time.Parse("2006-01-02 15:04:05", CreatedAt)
		if err != nil {
			return nil, err
		}
		Data.CreatedAt = parsedTime.Format("Jan 02 2006 15:04")
		if QuotationUpdatedAt.Valid {
			parsedQuoteTime, err := time.Parse("2006-01-02 15:04:05", QuotationUpdatedAt.String)
			if err != nil {
				return nil, err
			}
			Data.QuotationUpdatedAt = parsedQuoteTime.Format("Jan 02 2006 15:04")
		} else {
			Data.QuotationUpdatedAt = "N/A"
		}

		if QuotationUpdater.Valid {
			Data.QuotationUpdater = QuotationUpdater.String
		} else {
			Data.QuotationUpdater = "N/A"
		}

		if DeletedReason.Valid {
			Data.DeletedReason = DeletedReason.String
		} else {
			Data.DeletedReason = "N/A"
		}

		if Status == "Invalid" {
			Data.Status = "Rejected"
		} else if Status == "Waiting for Approval" || Status == "Quotation Required" || Status == "Deal Approved" {
			Data.Status = "On Progress"
		} else {
			Data.Status = "QUOTED"
		}

		if IsCompany == 0 {
			CustomerQuery := `SELECT Name, Email, Mobile FROM contacts WHERE ID = ?`
			err := db.QueryRow(CustomerQuery, CompanyID).Scan(&Data.Customer, &Data.Email, &Data.Phone)
			if err != nil {
				return Datas, err
			}
		} else if IsCompany == 1 {
			CompanyQuery := `SELECT CompanyName, CompanyEmail, CompanyMobile FROM company WHERE ID = ?`
			err := db.QueryRow(CompanyQuery, CompanyID).Scan(&Data.Customer, &Data.Email, &Data.Phone)
			if err != nil {
				return Datas, err
			}
		}
		Datas = append(Datas, Data)
	}
	return Datas, nil
}
