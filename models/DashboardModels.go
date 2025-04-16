package models

import (
	"database/sql"
)

func Progress(Year string, ID int) (ProgressStruct, bool, error) {
	var datas ProgressStruct
	var err error

	if ID == 0 {
		query := `SELECT     
				SUM(CASE WHEN FinalOutCome = 'On Progress' THEN 1 ELSE 0 END) AS OnProgressCount,     
				SUM(CASE WHEN FinalOutCome = 'WON' THEN 1 ELSE 0 END) AS WonCount,
				SUM(CASE WHEN FinalOutCome = 'Loss' THEN 1 ELSE 0 END) AS LoseCount 
			FROM 
				Deals  
			WHERE 
				YEAR(DealStartDate) = ? AND IsDeleted = ?`
		err = db.QueryRow(query, Year, 0).Scan(&datas.OnProgress, &datas.Win, &datas.Loss)
	} else {
		query := `SELECT     
				SUM(CASE WHEN FinalOutCome = 'On Progress' THEN 1 ELSE 0 END) AS OnProgressCount,     
				SUM(CASE WHEN FinalOutCome = 'WON' THEN 1 ELSE 0 END) AS WonCount,
				SUM(CASE WHEN FinalOutCome = 'Loss' THEN 1 ELSE 0 END) AS LoseCount 
			FROM 
				Deals  
			WHERE 
				YEAR(DealStartDate) = ? AND IsDeleted = ? AND
				(DealOwnerID = ? OR DealApproverID = ? OR QuotationUpdaterID = ?)`
		err = db.QueryRow(query, Year, 0, ID, ID, ID).Scan(&datas.OnProgress, &datas.Win, &datas.Loss)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return datas, false, nil
		}
		return datas, false, err
	}

	return datas, true, nil
}
func ContactsCount(year string, ID int) ([]struct {
	Month    int
	Contacts int
}, bool, error) {
	var query string
	var rows *sql.Rows
	var err error

	query = `
		SELECT
			m.Month AS Month,
			COALESCE(COUNT(c.ID), 0) AS Contacts
		FROM
			MonthNumbers m
			LEFT JOIN contacts c ON m.Month = MONTH(c.CreatedAt) AND YEAR(c.CreatedAt) = ? AND IsDeleted = ? 
		GROUP BY
			m.Month
		ORDER BY
			m.Month;
	`
	rows, err = db.Query(query, year, 0)

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var contactCounts []struct {
		Month    int
		Contacts int
	}

	dataFound := false

	for rows.Next() {
		var month int
		var contacts int
		if err := rows.Scan(&month, &contacts); err != nil {
			return nil, false, err
		}
		contactCounts = append(contactCounts, struct {
			Month    int
			Contacts int
		}{
			Month:    month,
			Contacts: contacts,
		})
		if contacts > 0 {
			dataFound = true
		}
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return contactCounts, dataFound, nil
}

func CompanyCount(year string, ID int) ([]struct {
	Month    int
	Contacts int
}, bool, error) {
	var query string
	var rows *sql.Rows
	var err error

	query = `
		SELECT
			m.Month AS Month,
			COALESCE(COUNT(c.ID), 0) AS Contacts
		FROM
			MonthNumbers m
			LEFT JOIN company c ON m.Month = MONTH(c.CreatedAt) AND YEAR(c.CreatedAt) = ? AND IsDeleted = ?
		GROUP BY
			m.Month
		ORDER BY
			m.Month;
	`
	rows, err = db.Query(query, year, 0)

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var contactCounts []struct {
		Month    int
		Contacts int
	}

	dataFound := false

	for rows.Next() {
		var month int
		var contacts int
		if err := rows.Scan(&month, &contacts); err != nil {
			return nil, false, err
		}
		contactCounts = append(contactCounts, struct {
			Month    int
			Contacts int
		}{
			Month:    month,
			Contacts: contacts,
		})
		dataFound = true
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return contactCounts, dataFound, nil
}

func QuotationCount(year string, ID int) ([]struct {
	Month     int
	Quotation int
}, bool, error) {
	var query string
	var rows *sql.Rows
	var err error

	if ID == 0 {
		query = `
			SELECT
				m.Month AS Month,
				COALESCE(COUNT(c.QuotationUpdaterID), 0) AS Quotation
			FROM
				MonthNumbers m
				LEFT JOIN Deals c ON m.Month = MONTH(c.CreatedAt) 
					AND YEAR(c.CreatedAt) = ? 
					AND c.IsDeleted = ? 
					AND c.Quotation != ''
			GROUP BY
				m.Month
			ORDER BY
				m.Month;
		`
		rows, err = db.Query(query, year, 0)
	} else {
		query = `
			SELECT
				m.Month AS Month,
				COALESCE(COUNT(c.QuotationUpdaterID), 0) AS Quotation
			FROM
				MonthNumbers m
				LEFT JOIN Deals c ON m.Month = MONTH(c.CreatedAt) 
					AND YEAR(c.CreatedAt) = ? 
					AND c.IsDeleted = ? 
					AND c.Quotation != ''
					AND (c.DealOwnerID = ? OR c.DealApproverID = ? OR c.QuotationUpdaterID = ?)
			GROUP BY
				m.Month
			ORDER BY
				m.Month;
		`
		rows, err = db.Query(query, year, 0, ID, ID, ID)
	}

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var quotationCounts []struct {
		Month     int
		Quotation int
	}

	dataFound := false

	for rows.Next() {
		var month int
		var quotation int
		if err := rows.Scan(&month, &quotation); err != nil {
			return nil, false, err
		}
		quotationCounts = append(quotationCounts, struct {
			Month     int
			Quotation int
		}{
			Month:     month,
			Quotation: quotation,
		})
		dataFound = true
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return quotationCounts, dataFound, nil
}

func DealCounts(year string, ID int) ([]struct {
	Month     int
	WonCount  int
	LoseCount int
}, bool, error) {

	var query string
	var rows *sql.Rows
	var err error

	if ID == 0 {
		query = `
			SELECT
				m.Month AS Month,
				COALESCE(SUM(CASE WHEN d.FinalOutCome = 'WON' THEN 1 ELSE 0 END), 0) AS WonCount,
				COALESCE(SUM(CASE WHEN d.FinalOutCome = 'Loss' THEN 1 ELSE 0 END), 0) AS LoseCount
			FROM
				MonthNumbers m
				LEFT JOIN Deals d ON m.Month = MONTH(d.DealsEndDate) AND YEAR(d.DealsEndDate) = ?
			GROUP BY
				m.Month
			ORDER BY
				m.Month;	
		`
		rows, err = db.Query(query, year)
	} else {
		query = `
			SELECT
				m.Month AS Month,
				COALESCE(SUM(CASE WHEN d.FinalOutCome = 'WON' THEN 1 ELSE 0 END), 0) AS WonCount,
				COALESCE(SUM(CASE WHEN d.FinalOutCome = 'Loss' THEN 1 ELSE 0 END), 0) AS LoseCount
			FROM
				MonthNumbers m
				LEFT JOIN Deals d ON m.Month = MONTH(d.DealsEndDate) 
					AND YEAR(d.DealsEndDate) = ? 
					AND (d.DealOwnerID = ? OR d.DealApproverID = ? OR d.QuotationUpdaterID = ?)
			GROUP BY
				m.Month
			ORDER BY
				m.Month;		
		`
		rows, err = db.Query(query, year, ID, ID, ID)
	}

	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	var dealCounts []struct {
		Month     int
		WonCount  int
		LoseCount int
	}

	dataFound := false
	for rows.Next() {
		var month int
		var wonCount int
		var loseCount int
		if err := rows.Scan(&month, &wonCount, &loseCount); err != nil {
			return nil, false, err
		}
		dealCounts = append(dealCounts, struct {
			Month     int
			WonCount  int
			LoseCount int
		}{
			Month:     month,
			WonCount:  wonCount,
			LoseCount: loseCount,
		})
		dataFound = true
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return dealCounts, dataFound, nil
}

func RecentDeals(Status string, ID int) ([]RecentDealsStruct, bool, error) {
	var Datas []RecentDealsStruct
	var query string
	var rows *sql.Rows
	var err error

	if ID == 0 {
		query = `SELECT ID, IsCompany, CompanyID, ProductType, Pricing, Quantity, DealOwner, DealOwnerID, DealApprover, DealApproverID, QuotationUpdater, QuotationUpdaterID, FinalOutCome FROM Deals WHERE FinalOutCome = ? AND IsDeleted = ? ORDER BY DealsEndDate DESC LIMIT 5`
		rows, err = db.Query(query, Status, 0)
	} else {
		query = `SELECT ID, IsCompany, CompanyID, ProductType, Pricing, Quantity, DealOwner, DealOwnerID, DealApprover, DealApproverID, QuotationUpdater, QuotationUpdaterID, FinalOutCome FROM Deals WHERE FinalOutCome = ? AND IsDeleted = ? AND (DealOwnerID = ? OR DealApproverID = ? OR QuotationUpdaterID = ?) ORDER BY DealsEndDate DESC LIMIT 5`
		rows, err = db.Query(query, Status, 0, ID, ID, ID)
	}

	if err != nil {
		return Datas, false, err
	}

	for rows.Next() {
		var data RecentDealsStruct
		if err := rows.Scan(&data.ID, &data.IsCompany, &data.CompanyID, &data.ProductType, &data.Price, &data.Quantity, &data.OwnerName, &data.OwnerID, &data.ApproverName, &data.ApproverID, &data.QuotationUpdaterName, &data.QuotationUpdaterID, &data.Status); err != nil {
			return Datas, false, err
		}
		if data.IsCompany == 0 {
			query := `SELECT ID, Name FROM contacts WHERE ID = ? `
			err := db.QueryRow(query, data.CompanyID).Scan(&data.DealID, &data.Name)
			if err != nil {

				return Datas, false, err
			}
		} else if data.IsCompany == 1 {
			query := `SELECT ID, CompanyName FROM company WHERE ID = ? `
			err := db.QueryRow(query, data.CompanyID).Scan(&data.DealID, &data.Name)
			if err != nil {

				return Datas, false, err
			}
		}
		Datas = append(Datas, data)
	}
	return Datas, true, nil
}

// BY week data
func ContactsCountByWeek(year string, month int) ([]struct {
	Week     int
	Contacts int
}, bool, error) {
	query := `
        SELECT
            WEEK(c.CreatedAt, 3) AS Week,
            COALESCE(COUNT(c.ID), 0) AS Contacts
        FROM
            contacts c
        WHERE
            YEAR(c.CreatedAt) = ? AND MONTH(c.CreatedAt) = ?
        GROUP BY
            WEEK(c.CreatedAt, 3)
        ORDER BY
            WEEK(c.CreatedAt, 3);
    `

	rows, err := db.Query(query, year, month)
	if err != nil {
		return nil, false, err
	}
	defer rows.Close()

	contactCounts := make([]struct {
		Week     int
		Contacts int
	}, 5)

	for i := range contactCounts {
		contactCounts[i].Week = i + 1
		contactCounts[i].Contacts = 0
	}

	dataFound := false

	for rows.Next() {
		var week int
		var contacts int
		if err := rows.Scan(&week, &contacts); err != nil {
			return nil, false, err
		}
		if week > 0 && week <= 5 {
			contactCounts[week-1].Contacts = contacts
		}
		dataFound = true
	}
	if err := rows.Err(); err != nil {
		return nil, false, err
	}

	return contactCounts, dataFound, nil
}
