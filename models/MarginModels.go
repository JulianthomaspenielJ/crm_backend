package models

import (
	"database/sql"
	"fmt"
)

func AddMargin(MarginType string, Percentage int, CreatedBY string) (string, error) {
	var Margin string
	CheckQuery := `SELECT MarginType FROM Margin WHERE MarginType = ?`

	err := db.QueryRow(CheckQuery, MarginType).Scan(&Margin)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error checking existing margin type: %v", err)
	}

	if Margin != "" {
		res := "Margin Type already Exist"
		return res, err
	}

	query := `INSERT INTO Margin(MarginType, BPPercentage, IsActive, CreatedBy, CreatedAt, IsDelete) VALUES (?, ?, ?, ?, Now(), ?)`
	result, err := db.Exec(query, MarginType, Percentage, 1, CreatedBY, 0)

	if err != nil {
		return "", fmt.Errorf("error executing margin insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting last insert ID: %v", err)
	}

	res, err := InsertAudit("Create", CreatedBY, id, "127.1.0", "Data inserted to the Margin table")
	if err != nil {
		return "", fmt.Errorf("error occurred while audit data inserting: %v", err)
	}
	fmt.Print(res)

	return "", nil
}

func AddMarginValues(MarginID int, Tier string, TierPercentage int, CreatedBY string) (string, error) {
	var ExistTier string
	CheckQuery := `SELECT Tier FROM MarginValues WHERE Tier = ? AND MarginID = ?`

	err := db.QueryRow(CheckQuery, Tier, MarginID).Scan(&ExistTier)
	if err != nil && err != sql.ErrNoRows {
		return "", fmt.Errorf("error checking existing margin type: %v", err)
	}

	if ExistTier != "" {
		res := "Tier already Exist"
		return res, err
	}

	query := `INSERT INTO MarginValues(MarginID, Tier, TierPercentage, IsActive, CreatedBy, CreatedAt) VALUES (?, ?, ?, ?, ?, Now())`
	result, err := db.Exec(query, MarginID, Tier, TierPercentage, 1, CreatedBY)

	if err != nil {
		return "", fmt.Errorf("error executing margin insert: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return "", fmt.Errorf("error getting last insert ID: %v", err)
	}

	res, err := InsertAudit("Create", CreatedBY, id, "127.1.0", "Data inserted to the Margin_Values table")
	if err != nil {
		return "", fmt.Errorf("error occurred while audit data inserting: %v", err)
	}
	fmt.Print(res)

	return "", nil
}

func GetMargin() ([]Margin, error) {
	var margin []Margin

	query := `SELECT ID, MarginType FROM Margin WHERE IsActive = ? AND IsDelete = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("error querying margin table: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mar Margin
		if err := rows.Scan(&mar.ID, &mar.MarginType); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data %s", err)
		}
		margin = append(margin, mar)
	}
	return margin, nil
}

func GetMarginTableData() ([]Margin, error) {
	var margin []Margin

	query := `SELECT ID, MarginType, BPPercentage, CreatedBy, CreatedAt FROM Margin WHERE IsActive = ? AND IsDelete = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 1, 0)
	if err != nil {
		return nil, fmt.Errorf("error querying margin table: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mar Margin
		if err := rows.Scan(&mar.ID, &mar.MarginType, &mar.BPPercentage, &mar.CreatedBy, &mar.CreatedAt); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data %s", err)
		}
		margin = append(margin, mar)
	}
	return margin, nil
}

func GetMarginValuesTableData() ([]MarginValues, error) {
	var margin []MarginValues

	query := `SELECT ID, MarginID, Tier, TierPercentage, CreatedBy, CreatedAt FROM MarginValues WHERE IsActive = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 1)
	if err != nil {
		return nil, fmt.Errorf("error querying margin table: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var mar MarginValues
		if err := rows.Scan(&mar.ID, &mar.MarginID, &mar.Tier, &mar.TierPercentage, &mar.CreatedBy, &mar.CreatedAt); err != nil {
			return nil, fmt.Errorf("error occurred while scanning data %s", err)
		}
		GetMarginName := `SELECT MarginType FROM Margin WHERE ID = ?`
		err := db.QueryRow(GetMarginName, mar.MarginID).Scan(&mar.MarginName)
		if err != nil {
			return margin, fmt.Errorf("error occurred while query executing %s", err)
		}
		if err == sql.ErrNoRows {
			return margin, fmt.Errorf("no contact found with ID %s", err)
		}
		margin = append(margin, mar)
	}
	return margin, nil
}

func UpdateMargin(MarginType, UpdatedBy string, Percentage, ID int) error {
	query := `UPDATE Margin SET MarginType = ?, BPPercentage = ?, UpdatedBy =?, UpdatedAt = Now() WHERE ID = ?`
	res, err := db.Exec(query, MarginType, Percentage, UpdatedBy, ID)
	if err != nil {
		return err
	}
	fmt.Printf("data %s", res)
	result, err := InsertAudit("Update", UpdatedBy, int64(ID), "127.0.0.0", "Margin data updated successfully")
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func UpdateMarginValues(Tier, UpdatedBy string, Percentage, ID int) (string, error) {
	// checkQuery := `SELECT ID FROM MarginValues WHERE Tier = ? AND IsActive = ?`
	// var existingID int
	// err := db.QueryRow(checkQuery, Tier, 1).Scan(&existingID)
	// if err != nil {
	// 	if err != sql.ErrNoRows {
	// 		return "", err
	// 	}
	// }

	// if existingID != 0 {
	// 	return "exist", fmt.Errorf("tier name already exists")
	// }

	updateQuery := `UPDATE MarginValues SET Tier = ?, TierPercentage = ?, UpdatedBy = ?, UpdatedAt = NOW() WHERE ID = ?`
	_, err := db.Exec(updateQuery, Tier, Percentage, UpdatedBy, ID)
	if err != nil {
		return "", err
	}

	_, err = InsertAudit("Update", UpdatedBy, int64(ID), "127.0.0.0", "Margin_Values data updated successfully")
	if err != nil {
		return "", err
	}

	return "", nil
}

func GetMarginValuesData(ID int) ([]MarginValues, error) {
	var margins []MarginValues
	query := `SELECT ID, Tier, TierPercentage FROM MarginValues WHERE MarginID = ? AND IsActive = ? ORDER BY Tier DESC`
	res, err := db.Query(query, ID, 1)
	if err != nil {
		return nil, err
	}

	for res.Next() {
		var margin MarginValues
		if err := res.Scan(&margin.ID, &margin.Tier, &margin.TierPercentage); err != nil {
			return nil, err
		}
		margins = append(margins, margin)
	}
	return margins, nil
}
func DeleteMargin(UpdatedBy string, ID int) error {
	query := `UPDATE Margin SET IsDelete = ? , UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, 1, UpdatedBy, ID)
	if err != nil {
		return err
	}
	result, err := InsertAudit("Delete", UpdatedBy, int64(ID), "127.0.0.0", "Margin data deleted successfully")
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func DeleteMarginValues(UpdatedBy string, ID int) error {
	query := `UPDATE MarginValues SET IsActive = ? , UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, 0, UpdatedBy, ID)
	if err != nil {
		return err
	}
	result, err := InsertAudit("Delete", UpdatedBy, int64(ID), "127.0.0.0", "Margin_Values data deleted successfully")
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func RecoverMargin(UpdatedBy string, ID int) error {
	query := `UPDATE Margin SET IsDelete = ? , UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, 0, UpdatedBy, ID)
	if err != nil {
		return err
	}
	result, err := InsertAudit("Recovered", UpdatedBy, int64(ID), "127.0.0.0", "Margin data Recovered successfully")
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}

func RecoverMarginValues(UpdatedBy string, ID int) error {
	query := `UPDATE MarginValues SET IsActive = ? , UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, 1, UpdatedBy, ID)
	if err != nil {
		return err
	}
	result, err := InsertAudit("Recovered", UpdatedBy, int64(ID), "127.0.0.0", "Margin_Values data Recovered successfully")
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}
