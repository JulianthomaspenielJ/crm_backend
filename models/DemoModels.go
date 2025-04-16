package models

import (
	"database/sql"
	"fmt"
	"time"
)

var timeNow = time.Now().Format("02 Jan 2006")

func AddDemo(Data Demo) error {
	Status := "Requested"
	Query := `INSERT INTO Demo(SalesPerson, EstimatedDate, ProductType, Processor, RAM, Storage, Notes, Status, Commants, SalesPersonID, IsApproved, IsDeleted, CreatedAt)
				VALUES(?,?,?,?,?,?,?,?,?,?,?,?,Now())`
	res, err := db.Exec(Query, Data.SalesPerson, Data.EstimatedDate, Data.ProductType, Data.Processor, Data.RAM, Data.Storage, Data.Notes, Status, Data.Comments, Data.SalesPersonID, 0, 0)
	if err != nil {
		return err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	err = DemoSummary(int(ID), fmt.Sprintf("New demo has been requested by %s at %s. And Now it was awaiting for the Permission to the Products Team.", Data.SalesPerson, timeNow), Data.SalesPerson)
	if err != nil {
		return err
	}
	return nil
}

func GetDemo(Team, Role string, ID int) ([]Demo, error) {
	var demos []Demo
	var Comment sql.NullString
	Status := "Approved"
	var rows *sql.Rows
	var err error

	if Team == "Admin" || Team == "Product" {
		Query := `SELECT ID, SalesPerson, EstimatedDate, ProductType, Processor, RAM, Storage, Notes, Status, Commants, SalesPersonID, CreatedAt FROM Demo ORDER BY CreatedAt DESC`
		rows, err = db.Query(Query)
	} else {
		Query := `SELECT ID, SalesPerson, EstimatedDate, ProductType, Processor, RAM, Storage, Notes, Status, Commants, SalesPersonID, CreatedAt FROM Demo WHERE Status = ? ORDER BY CreatedAt DESC`
		rows, err = db.Query(Query, Status)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var data Demo
		err := rows.Scan(&data.ID, &data.SalesPerson, &data.EstimatedDate, &data.ProductType, &data.Processor, &data.RAM, &data.Storage, &data.Notes, &data.Status, &Comment, &data.SalesPersonID, &data.CreatedAt)
		if err != nil {
			return nil, err
		}

		if Comment.Valid {
			data.Comments = Comment.String
		} else {
			data.Comments = ""
		}

		demos = append(demos, data)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return demos, nil
}

func GetDemoData(ID int) (Demo, error) {
	var Data Demo
	var Comment sql.NullString

	Query := `SELECT ID, SalesPerson, EstimatedDate, ProductType, Processor, RAM, Storage, Notes, Status, Commants, SalesPersonID, CreatedAt FROM Demo WHERE ID = ?`
	row := db.QueryRow(Query, ID)
	err := row.Scan(&Data.ID, &Data.SalesPerson, &Data.EstimatedDate, &Data.ProductType, &Data.Processor, &Data.RAM, &Data.Storage, &Data.Notes, &Data.Status, &Comment, &Data.SalesPersonID, &Data.CreatedAt)
	if err != nil {
		return Demo{}, err
	}
	if Comment.Valid {
		Data.Comments = Comment.String
	} else {
		Data.Comments = "N/A"
	}
	return Data, nil
}

func GetDemoSummary(DemoID int) ([]DemoSummaryStruct, error) {
	var summaries []DemoSummaryStruct

	Query := `SELECT ID, DemoID, Description, ByPerson, HappenedAT FROM DemoSummary WHERE DemoID = ? ORDER BY HappenedAT ASC`
	rows, err := db.Query(Query, DemoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var summary DemoSummaryStruct
		if err := rows.Scan(&summary.ID, &summary.DemoID, &summary.Description, &summary.ByPerson, &summary.HappenedAT); err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func UpdateDemo(ID int, Status, UpdatedBy string) error {
	Query := `UPDATE Demo SET Status = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(Query, Status, UpdatedBy, ID)
	if err != nil {
		return err
	}
	var description string

	switch Status {
	case "Approved":
		description = fmt.Sprintf("Demo request has been approved by %s at %s. The demo product is now awaiting further action.", UpdatedBy, timeNow)
	case "Rejected":
		description = fmt.Sprintf("Demo request has been rejected by %s at %s. No further action will be taken.", UpdatedBy, timeNow)
	case "Ready to Delivery":
		description = fmt.Sprintf("Demo request has been marked 'Ready to Deliver' by %s at %s. The product will be delivered soon.", UpdatedBy, timeNow)
	case "Delivered":
		description = fmt.Sprintf("Demo product has been delivered by %s at %s. The process is now complete.", UpdatedBy, timeNow)
	default:
		description = fmt.Sprintf("Demo request status updated to %s by %s at %s.", Status, UpdatedBy, timeNow)
	}

	err = DemoSummary(int(ID), description, UpdatedBy)
	if err != nil {
		return err
	}
	return nil
}

func DemoSummary(ID int, Description, ByPerson string) error {
	SummaryQuery := `INSERT INTO DemoSummary(DemoID, Description, ByPerson, HappenedAT) 
					VALUES(?, ?, ?, Now())`
	_, err := db.Exec(SummaryQuery, ID, Description, ByPerson)
	if err != nil {
		return err
	}
	return nil
}
