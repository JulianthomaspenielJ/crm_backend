package models

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func AddDeal(deal DealsStruct) (DealsStruct, error) {
	var lastTicketNo string
	err := db.QueryRow("SELECT TicketNo FROM Deals ORDER BY ID DESC LIMIT 1").Scan(&lastTicketNo)
	if err != nil && err != sql.ErrNoRows {
		return deal, fmt.Errorf("failed to get last TicketNo: %v", err)
	}

	newTicketNo := "HD00001"
	if lastTicketNo != "" {
		lastNumberStr := lastTicketNo[2:]
		lastNumber, err := strconv.Atoi(lastNumberStr)
		if err != nil {
			return deal, fmt.Errorf("failed to parse last TicketNo: %v", err)
		}
		newNumber := lastNumber + 1
		newTicketNo = fmt.Sprintf("HD%05d", newNumber)
	}
	status := "Waiting for Approval"
	approver := "-"
	Team := deal.Team
	query := `INSERT INTO Deals(TicketNo, IsCompany, CompanyID, ProductType, ProductDetails, DealOwner, DealOwnerID, Quantity, DealStartDate, ApproximateClosingDate, Priority, CreatedBy, CreatedAt, IsDeleted, IsApproved, Status, DealApproverID, DealApprover, FinalOutCome)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, Now(), ?, ?, ?, ?, ?, ?)`

	if deal.DealsStartDate == "" || deal.ApproximateEndDate == "" {
		return deal, fmt.Errorf("error occurred: empty date")
	}

	res, err := db.Exec(query, newTicketNo, deal.IsCompany, deal.CompanyID, deal.ProductType, deal.ProductDetails, deal.DealOwner, deal.DealOwnerID, deal.Quantity, deal.DealsStartDate, deal.ApproximateEndDate, deal.Priority, deal.CreatedBy, 0, 0, status, 0, approver, "On Progress")
	if err != nil {
		return deal, err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return deal, err
	}
	summaryTime := time.Now().Format("02 Jan 2006")
	err = DealSummary(int(ID), fmt.Sprintf("The deal was created By %s at %s and now it was waiting for the pricing team approval", deal.CreatedBy, summaryTime), deal.CreatedBy)
	if err != nil {
		return deal, err
	}
	err = UpdateAudit("Create", deal.CreatedBy, int(ID), "127.0.0.1", "Deal created successfully")
	if err != nil {
		return deal, err
	}
	contact := ""
	var GetQuery string
	if deal.IsCompany == 1 {
		GetQuery = `SELECT CompanyName FROM company WHERE ID = ?`
		contact = "Company Name"
	} else if deal.IsCompany == 0 {
		GetQuery = `SELECT Name FROM contacts WHERE ID = ?`
		contact = "Contact Name"
	}

	err = db.QueryRow(GetQuery, deal.CompanyID).Scan(&deal.CompanyName)
	if err != nil {
		return deal, fmt.Errorf("error retrieving company/contact name: %v", err)
	}
	var To []string
	//To := []string{"marketing@holoware.co", "pricing@holoware.co"}
	// Cc := []string{"marketing@holoware.co"}
	Cc := []string{deal.From}
	if Team == "Sales" {
		To = []string{"pricing@holoware.co", "sathish.pk@holoware.co"}
	} else if Team == "Marketing" {
		To = []string{"marketing@holoware.co", "pricing@holoware.co"}
	}

	currentDate := time.Now().Format("02 Jan 2006")
	Subject := fmt.Sprintf("New Deal Created for %s on %s and the Ticket ID is %s", deal.CompanyName, currentDate, newTicketNo)
	Body := fmt.Sprintf(`
		<html>
		<body>
			<p>Hello Team,</p>
			<p>A new deal has been created by <strong>%s</strong> and is currently awaiting admin approval.</p>
			<p><strong>Details of the deal:</strong></p>
			<ul>
				<li><strong>Product Type:</strong> %s</li>
				<li><strong>Description:</strong> %s</li>
				<li><strong>Quantity:</strong> %s</li>
				<li><strong>%s:</strong> %s</li>
			</ul>
			<p>Thanks & Regards,<br>Team CRM</p>
		</body>
		</html>
		`,
		deal.DealOwner, deal.ProductType, deal.ProductDetails, deal.Quantity, contact, deal.CompanyName,
	)
	deal.To = To
	deal.CC = Cc
	//deal.BCC = Bcc
	deal.Subject = Subject
	deal.Body = Body

	err = AddNotification(int(ID), deal.DealOwnerID, "Deal Created, Now It's Waiting for approval.")
	if err != nil {
		fmt.Print("Deal was not created")
		return deal, nil
	}

	GetManagerQuery := `SELECT ID FROM users WHERE (Team = ? AND Role = ?) AND IsDeleted = ? AND IsApproved = ?`

	result, err := db.Query(GetManagerQuery, "Pricing", "Manager", 0, 1)
	if err != nil {
		return deal, err
	}

	for result.Next() {
		var id int
		err := result.Scan(&id)
		if err != nil {
			return deal, err
		}

		err = AddNotification(int(ID), id, fmt.Sprintf("Deal Created by %s, Now It's Waiting for your approval.", deal.DealOwner))
		if err != nil {
			fmt.Print("Deal was not created")
			return deal, nil
		}
	}

	GetAdminQuery := `SELECT ID FROM users WHERE (Role = ? OR Role = ?) AND IsDeleted = ? AND IsApproved = ?`

	adminres, err := db.Query(GetAdminQuery, "Admin", "Super Admin", 0, 1)
	if err != nil {
		fmt.Print("While executing admin query", err)
		return deal, err
	}

	for adminres.Next() {
		var id int
		err := adminres.Scan(&id)
		if err != nil {
			fmt.Print("While scanning admin query", err)
			return deal, err
		}

		err = AddNotification(int(ID), id, fmt.Sprintf("Deal Created by %s, Now It's Waiting for the Pricing team approval.", deal.DealOwner))
		if err != nil {
			fmt.Print("Deal was not created")
			return deal, nil
		}
	}

	return deal, nil
}

func GetApproveDeal() ([]DealsStruct, error) {
	query := `SELECT ID, IsCompany, CompanyID, ProductType, DealOwner, Quantity, ApproximateClosingDate, Priority, CreatedAt FROM Deals WHERE IsApproved = ? AND IsDeleted = ? ORDER BY CreatedAt DESC`
	res, err := db.Query(query, 0, 0)
	if err != nil {
		return nil, err
	}
	var deals []DealsStruct
	for res.Next() {
		var createdAt string
		var deal DealsStruct
		err := res.Scan(&deal.ID, &deal.IsCompany, &deal.CompanyID, &deal.ProductType, &deal.DealOwner, &deal.Quantity, &deal.ApproximateEndDate, &deal.Priority, &createdAt)
		if err != nil {
			return nil, err
		}
		createAt, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err != nil {
			return nil, err
		}
		deal.CreatedAt = createAt.Format("2006-01-02 15:04:05")
		if deal.IsCompany == 0 {
			getCompanyName := `SELECT ID, Name FROM contacts WHERE ID = ? `
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, fmt.Errorf("error occurred while fetching contacts details %s", err)
			}

		} else if deal.IsCompany == 1 {
			getCompanyName := `SELECT ID, CompanyName FROM company WHERE ID = ? `
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, err
			}
		} else {
			deal.CompanyName = "Undefine"
			deal.ContactID = 0
		}
		deals = append(deals, deal)
	}
	return deals, nil
}

func GetDetailedDeal(ID, NotificationID int) (DealsStruct, []DealSummaryStruct, error) {
	var deal DealsStruct
	var summary []DealSummaryStruct
	var quotation, Reason, StatusCommand, MailStatus, QuotationUpdater sql.NullString
	var createdAtStr string
	var QuotationUpdaterID, ApproverID sql.NullInt16

	query := `SELECT ID, IsCompany, CompanyID, ProductType, ProductDetails, DealOwner, DealOwnerID, Quantity, Priority, DealApprover, DealStartDate, ApproximateClosingDate, DealsEndDate, Status, StatusCommands, DeletedReason, CreatedBy, CreatedAt, Quotation, FinalOutCome, MailStatus, QuotationUpdaterID, QuotationUpdater, DealApproverID, TicketNo FROM Deals WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(
		&deal.ID,
		&deal.IsCompany,
		&deal.CompanyID,
		&deal.ProductType,
		&deal.ProductDetails,
		&deal.DealOwner,
		&deal.DealOwnerID,
		&deal.Quantity,
		&deal.Priority,
		&deal.DealApprover,
		&deal.DealsStartDate,
		&deal.ApproximateEndDate,
		&deal.DealEndDate,
		&deal.Status,
		&StatusCommand,
		&Reason,
		&deal.CreatedBy,
		&createdAtStr,
		&quotation,
		&deal.FinalOutCome,
		&MailStatus,
		&QuotationUpdaterID,
		&QuotationUpdater,
		&ApproverID,
		&deal.TicketID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return deal, summary, nil
		}
		return deal, summary, fmt.Errorf("error occurred while executing query: %v", err)
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr)
	if err != nil {
		return deal, summary, fmt.Errorf("error occurred while parsing createdAtStr: %v", err)
	}
	deal.CreatedAt = createdAt.Format("02 Jan 2006")

	if QuotationUpdaterID.Valid {
		deal.QuotationUpdaterID = int(QuotationUpdaterID.Int16)
	} else {
		deal.QuotationUpdaterID = 0
	}
	if ApproverID.Valid {
		deal.DealApproverID = int(ApproverID.Int16)
	} else {
		deal.DealApproverID = 0
	}

	if QuotationUpdater.Valid {
		deal.QuotationUpdaterName = QuotationUpdater.String
	} else {
		deal.QuotationUpdaterName = "N/A"
	}
	EmailQuery := `SELECT Email FROM users WHERE ID = ?`
	err = db.QueryRow(EmailQuery, deal.DealOwnerID).Scan(&deal.DealOwnerEmail)
	if err != nil {
		return deal, summary, err
	}

	var Mobile sql.NullString
	if deal.IsCompany == 0 {
		getContactName := `SELECT ID, Name, Email, Mobile FROM contacts WHERE ID = ?`
		err = db.QueryRow(getContactName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName, &deal.Email, &Mobile)
		if err != nil {
			return deal, summary, err
		}
	} else if deal.IsCompany == 1 {
		getCompanyName := `SELECT ID, CompanyName, CompanyEmail, CONCAT(City, ', ', State, ', ', Country) AS Address, CompanyMobile FROM company WHERE ID = ?`
		err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName, &deal.Email, &deal.Address, &Mobile)
		if err != nil {
			return deal, summary, err
		}
	} else {
		deal.CompanyName = "Individual"
		deal.ContactID = 0
	}
	if MailStatus.Valid {
		deal.MailStatus = MailStatus.String
	} else {
		deal.MailStatus = "N/A"
	}
	if Reason.Valid {
		deal.Reason = Reason.String
	} else {
		deal.Reason = "N/A"
	}

	if Mobile.Valid {
		deal.Mobile = Mobile.String
	} else {
		deal.Mobile = "N/A"
	}
	if StatusCommand.Valid {
		deal.StatusCommand = StatusCommand.String
	} else {
		deal.StatusCommand = "N/A"
	}
	if quotation.Valid {
		fmt.Print("Processing quotation files...\n")
		quotationFiles := strings.Split(quotation.String, ",")
		var base64Quotes []string

		for _, fileName := range quotationFiles {
			parts := strings.Split(fileName, "_")
			if len(parts) < 4 {
				fmt.Printf("Invalid file path format: %s\n", fileName)
				base64Quotes = append(base64Quotes, "No Data")
				continue
			}

			datePart := fmt.Sprintf("%s_%s", parts[0], parts[1])
			directory := filepath.Join("Assets", "Deals", "Quotation", fmt.Sprintf("%d", deal.ID), datePart)
			fullPath := filepath.Join(directory, fileName)

			if _, err := os.Stat(fullPath); err == nil {
				fileBytes, err := ioutil.ReadFile(fullPath)
				if err != nil {
					return deal, summary, err
				}
				base64Quotes = append(base64Quotes, base64.StdEncoding.EncodeToString(fileBytes)) // Convert file to Base64
			} else {
				fmt.Printf("File not found: %s\n", fullPath)
				base64Quotes = append(base64Quotes, "No Data")
			}
		}

		if len(base64Quotes) > 0 {
			deal.Quotation = base64Quotes
		} else {
			deal.Quotation = nil
		}
	} else {
		deal.Quotation = nil
	}

	err = UpdateNotification(NotificationID)
	if err != nil {
		fmt.Printf("Data was not updated , %s %d", err, NotificationID)
		return deal, summary, err
	}
	fmt.Printf("Updated Notification %d", NotificationID)
	summary, err = GetDealData(ID)
	if err != nil {
		return deal, summary, err
	}

	ProductQuery := `SELECT Products, RefarenceName FROM product_data WHERE DealID = ? ORDER BY ID DESC LIMIT 1`
	var productDetails string
	err = db.QueryRow(ProductQuery, deal.ID).Scan(&productDetails, &deal.ModelName)
	if err != nil {
		if err == sql.ErrNoRows {
			return deal, summary, nil
		}
		return deal, summary, fmt.Errorf("error occurred while fetching product details: %v", err)
	}

	// Assuming the product details are stored in `Products` as JSON
	var productData DealProductData

	// Unmarshal the product details into the `productData` struct
	err = json.Unmarshal([]byte(productDetails), &productData)
	if err != nil {
		return deal, summary, fmt.Errorf("error occurred while parsing product data: %v", err)
	}

	// Set the values from `productData` into the `deal` struct
	deal.MarginType = productData.MarginType
	deal.MarginTypePercentage = productData.MarginTypePercentage

	// Fetch MarginID from the Margin table based on MarginType
	var marginID int
	marginQuery := `SELECT ID FROM Margin WHERE MarginType = ? AND IsActive = 1 LIMIT 1`
	err = db.QueryRow(marginQuery, productData.MarginType).Scan(&marginID)
	if err != nil {
		if err == sql.ErrNoRows {
			return deal, summary, nil
		}
		return deal, summary, fmt.Errorf("error occurred while fetching MarginID: %v", err)
	}

	// Fetch the Tier from MarginValues based on MarginID
	var tier string
	tierQuery := `SELECT Tier FROM MarginValues WHERE MarginID = ? AND IsActive = 1 LIMIT 1`
	err = db.QueryRow(tierQuery, marginID).Scan(&tier)
	if err != nil {
		if err == sql.ErrNoRows {
			return deal, summary, nil
		}
		return deal, summary, fmt.Errorf("error occurred while fetching tier: %v", err)
	}

	// Set the retrieved values in the deal struct
	deal.Tier = tier

	return deal, summary, nil
}

func GetQuotationDeal(ID int, Team string) ([]DealsStruct, error) {
	status := "Quotation Required"
	var query string
	var res *sql.Rows
	var err error
	var quotaionUpdaterStr sql.NullString

	switch Team {
	case "Sales", "Marketing":
		userIDs, err := GetUserIDsByTeam(Team)
		if err != nil {
			return nil, err
		}

		if len(userIDs) == 0 {
			return nil, fmt.Errorf("no users found for team %s", Team)
		}

		placeholders := strings.Repeat("?,", len(userIDs))
		placeholders = placeholders[:len(placeholders)-1]

		query = fmt.Sprintf(`SELECT ID, IsCompany, CompanyID, ProductType, DealOwner, QuotationUpdater, DealApprover, Quantity,ApproximateClosingDate, Priority, CreatedAt FROM Deals WHERE Status = ? AND IsDeleted = ? AND DealOwnerID IN (%s) ORDER BY CreatedAt DESC`, placeholders)

		args := make([]interface{}, len(userIDs)+2)
		args[0] = status
		args[1] = 0
		for i, id := range userIDs {
			args[i+2] = id
		}

		res, err = db.Query(query, args...)
		if err != nil {
			return nil, err
		}
	default:
		query = `SELECT ID, IsCompany, CompanyID, ProductType, DealOwner, QuotationUpdater, DealApprover, Quantity, ApproximateClosingDate, Priority, CreatedAt FROM Deals WHERE Status = ? AND IsDeleted = ? ORDER BY CreatedAt DESC`
		res, err = db.Query(query, status, 0)
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var deals []DealsStruct
	for res.Next() {
		var deal DealsStruct
		err := res.Scan(&deal.ID, &deal.IsCompany, &deal.CompanyID, &deal.ProductType, &deal.DealOwner, &quotaionUpdaterStr, &deal.DealApprover, &deal.Quantity,
			&deal.ApproximateEndDate, &deal.Priority, &deal.CreatedAt)
		if err != nil {
			return nil, err
		}
		if quotaionUpdaterStr.Valid {
			deal.QuotationUpdaterName = quotaionUpdaterStr.String
		} else {
			deal.QuotationUpdaterName = "N/A"
		}
		if deal.IsCompany == 0 {
			getCompanyName := `SELECT ID, Name FROM contacts WHERE ID = ?`
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, fmt.Errorf("error occurred while fetching contacts details %s", err)
			}
		} else if deal.IsCompany == 1 {
			getCompanyName := `SELECT ID, CompanyName FROM company WHERE ID = ?`
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, err
			}
		} else {
			deal.CompanyName = "Undefined"
			deal.ContactID = 0
		}

		deals = append(deals, deal)
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return deals, nil
}

func ApproveDeal(ID, ApproverID int, Approver, QuotationUpdater string, QuotationUpdaterID int, status string) (DealsStruct, error) {
	var deal DealsStruct

	tx, err := db.Begin()
	if err != nil {
		return deal, err
	}
	summaryTime := time.Now().Format("02 Jan 2006")

	updateQuery := `UPDATE Deals SET DealApprover = ?, DealApproverID = ?, QuotationUpdater = ?, QuotationUpdaterID = ?, IsApproved = ?, Status = ?, ApprovedAt = Now() WHERE ID = ?`
	_, err = tx.Exec(updateQuery, Approver, ApproverID, QuotationUpdater, QuotationUpdaterID, 1, status, ID)
	if err != nil {
		tx.Rollback()
		return deal, err
	}

	err = DealSummary(ID, fmt.Sprintf("This deal has been approved by %s at %s, And it has been assigned to %s for Quotation", Approver, summaryTime, QuotationUpdater), Approver)
	if err != nil {
		return deal, err
	}

	getQuery := `SELECT DealOwner ,DealOwnerID FROM Deals WHERE ID = ?`
	err = tx.QueryRow(getQuery, ID).Scan(&deal.DealOwner, &deal.DealOwnerID)
	if err != nil {
		return deal, err
	}
	err = AddNotification(int(ID), int(ApproverID), fmt.Sprintf("You have approved the deal successfully, and it was assigned to %s", QuotationUpdater))
	if err != nil {
		fmt.Print("error in approving", err)
		return deal, err
	}

	err = AddNotification(ID, int(QuotationUpdaterID), fmt.Sprintf("%s assigned to you a deal from %s.", Approver, deal.DealOwner))
	if err != nil {
		fmt.Print("error in quote updater", err)
		return deal, err
	}

	err = AddNotification(ID, int(deal.DealOwnerID), fmt.Sprintf("Your deal has been approved by %s and it was assigned to %s.", Approver, QuotationUpdater))
	if err != nil {
		fmt.Print("Error in owner", err)
		return deal, err
	}

	audit, err := InsertAudit("Update", Approver, int64(ID), "127.0.0.1", "Deal Approved Successfully")
	if err != nil {
		tx.Rollback()
		return deal, err
	}
	fmt.Print(audit)

	getDealQuery := `
		SELECT 
			d.IsCompany,
			d.CompanyID,
			d.ProductType, 
			d.ProductDetails, 
			d.Quantity,
			d.DealOwner, 
			d.DealOwnerID,
			u1.email AS DealOwnerEmail,
			u2.email AS DealApproverEmail,
			u3.email AS QuotationUpdaterEmail
		FROM 
			Deals d
		LEFT JOIN 
			users u1 ON d.DealOwnerID = u1.ID
		LEFT JOIN 
			users u2 ON d.DealApproverID = u2.ID
		LEFT JOIN 
			users u3 ON d.QuotationUpdaterID = u3.ID
		WHERE 
			d.ID = ?
	`

	var dealOwnerEmail, dealApproverEmail, quotationUpdaterEmail string
	err = tx.QueryRow(getDealQuery, ID).Scan(
		&deal.IsCompany,
		&deal.CompanyID,
		&deal.ProductType,
		&deal.ProductDetails,
		&deal.Quantity,
		&deal.DealOwner,
		&deal.DealOwnerID,
		&dealOwnerEmail,
		&dealApproverEmail,
		&quotationUpdaterEmail,
	)
	if err != nil {
		tx.Rollback()
		return deal, fmt.Errorf("error retrieving deal details and emails: %v", err)
	}

	var contactLabel string
	if deal.IsCompany == 1 {
		contactLabel = "Company Name"
		getCompanyQuery := `SELECT CompanyName FROM company WHERE ID = ?`
		err = tx.QueryRow(getCompanyQuery, deal.CompanyID).Scan(&deal.CompanyName)
	} else {
		contactLabel = "Contact Name"
		getContactQuery := `SELECT Name FROM contacts WHERE ID = ?`
		err = tx.QueryRow(getContactQuery, deal.CompanyID).Scan(&deal.CompanyName)
	}
	if err != nil {
		tx.Rollback()
		return deal, fmt.Errorf("error retrieving company/contact name: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return deal, err
	}

	To := []string{dealOwnerEmail, dealApproverEmail, quotationUpdaterEmail}
	Cc := []string{"marketing@holoware.co", "pricing@holoware.co"}

	currentDate := time.Now().Format("02 Jan 2006")
	Subject := fmt.Sprintf("Deal Approved for %s on %s", deal.CompanyName, currentDate)
	Body := fmt.Sprintf(`
		<html>
			<body>
				<p>Hello Team,</p>
				<p>A deal has been approved by <strong>%s</strong> and is now available to Upload quotation.</p>
				<p><strong>Details of the deal:</strong></p>
				<ul
					<li><strong>%s:</strong> %s</li>
					<li><strong>Product Type:</strong> %s</li>
					<li><strong>Description:</strong> %s</li>
					<li><strong>Quantity:</strong> %s</li>
				</ul>
				<p>The deal was created by <strong>%s</strong>, approved by <strong>%s</strong>, and the quotation updater is <strong>%s</strong>.</p>
				<p>Thanks & Regards,<br>Team CRM</p>
			</body>
		</html>
		`,
		Approver, contactLabel, deal.CompanyName, deal.ProductType, deal.ProductDetails, deal.Quantity, deal.DealOwner, Approver, QuotationUpdater,
	)
	deal.From = dealApproverEmail
	deal.To = To
	deal.CC = Cc
	deal.Subject = Subject
	deal.Body = Body

	return deal, nil
}

func UpdateQuotation(ID, PricingID int, FileName string, QuotationUpdaterID int, QuotationUpdaterName, TotalPrice string) (DealsStruct, error) {
	var deal DealsStruct
	fmt.Print("Deal ID", ID)
	tx, err := db.Begin()
	if err != nil {
		return deal, err
	}

	status := "Quotation Sent"
	query := `UPDATE Deals SET Quotation = ?, Status = ?, QuotationUpdater= ?, QuotationUpdaterID = ?, QuotationUpdatedAt = Now(), PricingID = ?, TotalPrice = ? WHERE ID = ?`
	res, err := db.Exec(query, FileName, status, QuotationUpdaterName, QuotationUpdaterID, PricingID, TotalPrice, ID)
	if err != nil {
		return deal, err
	}
	fmt.Print(res)
	audit, err := InsertAudit("Update", QuotationUpdaterName, int64(ID), "127.0.0.1", "Quotation updated in Deal")
	if err != nil {
		return deal, err
	}
	summaryTime := time.Now().Format("02 Jan 2006")
	summary := fmt.Sprintf("Quotation Updated by %s on %s", QuotationUpdaterName, summaryTime)
	err = DealSummary(ID, summary, QuotationUpdaterName)
	if err != nil {
		fmt.Print("error in summary", err)
		return deal, err
	}

	fmt.Print(audit)
	GetQuery := `
		SELECT 
			d.IsCompany,
			d.CompanyID,
			d.ProductType, 
			d.ProductDetails, 
			d.Quantity,
			d.DealOwner, 
			d.DealApprover,
			d.QuotationUpdater,
			u1.email AS DealOwnerEmail,
			u2.email AS DealApproverEmail,
			u3.email AS QuotationUpdaterEmail
		FROM 
			Deals d
		LEFT JOIN 
			users u1 ON d.DealOwnerID = u1.ID
		LEFT JOIN 
			users u2 ON d.DealApproverID = u2.ID
		LEFT JOIN 
			users u3 ON d.QuotationUpdaterID = u3.ID
		WHERE 
			d.ID = ?
	`

	var dealOwnerEmail, dealApproverEmail, quotationUpdaterEmail string
	err = db.QueryRow(GetQuery, ID).Scan(
		&deal.IsCompany,
		&deal.CompanyID,
		&deal.ProductType,
		&deal.ProductDetails,
		&deal.Quantity,
		&deal.DealOwner,
		&deal.DealApprover,
		&deal.QuotationUpdaterName,
		&dealOwnerEmail,
		&dealApproverEmail,
		&quotationUpdaterEmail,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return deal, fmt.Errorf("no matching deal found with ID: %d", ID)
		}
		return deal, fmt.Errorf("error retrieving deal details and emails: %v", err)
	}
	var contactLabel string
	if deal.IsCompany == 1 {
		contactLabel = "Company Name"
		getCompanyQuery := `SELECT CompanyName FROM company WHERE ID = ?`
		err = tx.QueryRow(getCompanyQuery, deal.CompanyID).Scan(&deal.CompanyName)
	} else {
		contactLabel = "Contact Name"
		getContactQuery := `SELECT Name FROM contacts WHERE ID = ?`
		err = tx.QueryRow(getContactQuery, deal.CompanyID).Scan(&deal.CompanyName)
	}
	if err != nil {
		tx.Rollback()
		return deal, fmt.Errorf("error retrieving company/contact name: %v", err)
	}
	To := []string{dealOwnerEmail, dealApproverEmail, quotationUpdaterEmail}
	Cc := []string{"marketing@holoware.co", "pricing@holoware.co"}

	currentDate := time.Now().Format("02 Jan 2006")
	Subject := fmt.Sprintf("Deal Quotation updated for %s on %s", deal.CompanyName, currentDate)
	Body := fmt.Sprintf(`
		<html>
			<body>
			<p>Hello Team,</p>
			<p>Quotation has been sent by <strong>%s</strong> for the deal.</p>
			<p><strong>Details of the deal:</strong></p>
				<ul>
					<li><strong>%s:</strong> %s</li>
					<li><strong>Product Type:</strong> %s</li>
					<li><strong>Description:</strong> %s</li>
					<li><strong>Quantity:</strong> %s</li>
				</ul>
			<p>The deal was created by <strong>%s</strong>, approved by <strong>%s</strong>, and the quotation updater is <strong>%s</strong>.</p>
				<p>Thanks & Regards,<br>Team CRM</p>
			</body>
		</html>
		`, QuotationUpdaterName, contactLabel, deal.CompanyName, deal.ProductType, deal.ProductDetails, deal.Quantity, deal.DealOwner, deal.DealApprover, deal.QuotationUpdaterName,
	)
	deal.From = quotationUpdaterEmail
	deal.To = To
	deal.CC = Cc
	deal.Subject = Subject
	deal.Body = Body

	return deal, nil
}
func FetchDeals(ID, IsCompany int) ([]DealsStruct, bool, error) {
	var Deals []DealsStruct
	query := `SELECT ID, DealOwner, ProductType, DealStartDate, Status, Quantity, Priority, ProductDetails, ApproximateClosingDate, DealApprover, FinalOutCome, TicketNo FROM Deals WHERE IsCompany = ? AND CompanyID = ? AND IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, IsCompany, ID, 0)
	if err != nil {
		return Deals, false, err
	}
	defer rows.Close()

	for rows.Next() {
		var Deal DealsStruct

		if err := rows.Scan(&Deal.ID, &Deal.DealOwner, &Deal.ProductType, &Deal.DealsStartDate, &Deal.Status, &Deal.Quantity, &Deal.Priority, &Deal.ProductDetails, &Deal.ApproximateEndDate, &Deal.DealApprover, &Deal.FinalOutCome, &Deal.TicketID); err != nil {
			return Deals, false, err
		}

		Deals = append(Deals, Deal)
	}

	return Deals, true, nil
}

func GetUserIDsByTeam(team string) ([]int, error) {
	query := `SELECT ID FROM users WHERE Team = ? AND IsDeleted = 0`
	rows, err := db.Query(query, team)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}
func GetAllDealsData(ID int, Team, Role string, UserID int, MyDeals bool) ([]DealsStruct, error) {
	var query string
	var res *sql.Rows
	var err error
	var QuotationUpdaterID, ApproverID sql.NullInt16
	var QuotationUploaderName sql.NullString

	if MyDeals {
		// If MyDeals filter is active, show deals where the logged-in user is involved
		query = `
		SELECT 
			ID, IsCompany, CompanyID, ProductType, DealOwner, Quantity,  
			Priority, CreatedAt, Status, DealApprover, QuotationUpdater, ProductDetails, 
			DealStartDate, ApproximateClosingDate, FinalOutCome, QuotationUpdaterID, DealApproverID, TicketNo
		FROM 
			Deals 
		WHERE 
			AND (DealOwner = ? OR QuotationUpdaterID = ? OR DealApproverID = ?)
		ORDER BY 
			CreatedAt DESC`
		res, err = db.Query(query, UserID, UserID, UserID)
	} else {
		// If user is Admin or Pricing, show all deals
		if Role == "Admin" || Role == "Super Admin" || Team == "Pricing" {
			query = `
			SELECT 
				ID, IsCompany, CompanyID, ProductType, DealOwner, Quantity, 
				Priority, CreatedAt, Status, DealApprover, QuotationUpdater, ProductDetails, 
				DealStartDate, ApproximateClosingDate, FinalOutCome, QuotationUpdaterID, DealApproverID, TicketNo
			FROM 
				Deals 
			ORDER BY 
				CreatedAt DESC`
			res, err = db.Query(query)
		} else {
			// Otherwise, show deals based on team or role
			query = `
			SELECT 
				ID, IsCompany, CompanyID, ProductType, DealOwner, Quantity,
				Priority, CreatedAt, Status, DealApprover, QuotationUpdater, ProductDetails, 
				DealStartDate, ApproximateClosingDate, FinalOutCome, QuotationUpdaterID, DealApproverID, TicketNo
			FROM 
				Deals 
			WHERE 
				IsDeleted = 0 
				AND (DealOwnerID = ? OR DealApproverID = ? 
				OR (DealOwnerID IN (SELECT ID FROM users WHERE Team = ? AND Role = 'Manager')))
			ORDER BY 
				CreatedAt DESC`
			res, err = db.Query(query, ID, ID, Team)
		}
	}

	if err != nil {
		return nil, err
	}
	defer res.Close()

	var deals []DealsStruct

	for res.Next() {
		var deal DealsStruct
		err := res.Scan(
			&deal.ID, &deal.IsCompany, &deal.CompanyID, &deal.ProductType, &deal.DealOwner,
			&deal.Quantity, &deal.Priority, &deal.CreatedAt, &deal.Status,
			&deal.DealApprover, &QuotationUploaderName, &deal.ProductDetails,
			&deal.DealsStartDate, &deal.ApproximateEndDate, &deal.FinalOutCome, &QuotationUpdaterID, &ApproverID, &deal.TicketID,
		)
		if err != nil {
			return nil, err
		}

		if QuotationUpdaterID.Valid {
			deal.QuotationUpdaterID = int(QuotationUpdaterID.Int16)
		} else {
			deal.QuotationUpdaterID = 0
		}

		if ApproverID.Valid {
			deal.DealApproverID = int(ApproverID.Int16)
		} else {
			deal.DealApproverID = 0
		}
		if QuotationUploaderName.Valid {
			deal.QuotationUpdaterName = QuotationUploaderName.String
		} else {
			deal.QuotationUpdaterName = "N/A"
		}

		// Fetch company name and other related data
		err = fetchCompanyName(&deal)
		if err != nil {
			return nil, err
		}
		//fmt.Printf("The User id %d and the deal ID %d", UserID, deal.ID)
		// Fetch notifications for the user and the deal
		notification, err := GetNotifications(UserID, deal.ID)
		if err != nil {
			return nil, err
		}

		deal.NotificationID = notification.NotificationID
		deal.NotificationMessage = notification.NotificationMessage
		deal.IsNotificationRead = notification.IsNotificationRead
		deal.NotificationCreatedAt = notification.NotificationCreatedAt

		deals = append(deals, deal)
	}

	if err = res.Err(); err != nil {
		return nil, err
	}

	return deals, nil
}

func fetchCompanyName(deal *DealsStruct) error {
	var query string

	if deal.IsCompany == 0 {
		query = `SELECT ID, Name, Email FROM contacts WHERE ID = ?`
	} else if deal.IsCompany == 1 {
		query = `SELECT ID, CompanyName, CompanyEmail FROM company WHERE ID = ?`
	} else {
		deal.CompanyName = "Undefined"
		deal.ContactID = 0
		return nil
	}

	err := db.QueryRow(query, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName, &deal.Email)
	if err != nil {
		return fmt.Errorf("error occurred while fetching company/contact details: %v", err)
	}

	return nil
}

func UpdateDealStatus(ID int, Status, UpdatedBy, Reason string) (DealsStruct, error) {
	var deal DealsStruct
	query := `UPDATE Deals SET Status = ?, StatusCommands = ?, UpdatedBy = ?, UpdatedAt = Now()  WHERE ID = ?`
	_, err := db.Query(query, Status, Reason, UpdatedBy, ID)
	if err != nil {
		return deal, err
	}

	summaryTime := time.Now().Format("02 Jan 2006")
	err = DealSummary(ID, fmt.Sprintf("The deal status has been updated to %s by %s at %s", Status, UpdatedBy, summaryTime), UpdatedBy)
	if err != nil {
		return deal, err
	}
	SelectQuery := `SELECT DealApproverID, QuotationUpdaterID, DealOwnerID FROM Deals WHERE ID = ?`
	err = db.QueryRow(SelectQuery, ID).Scan(&deal.DealApproverID, &deal.QuotationUpdaterID, &deal.DealOwnerID)
	if err != nil {
		return deal, err
	}
	err = AddNotification(int(ID), deal.DealApproverID, fmt.Sprintf("Deal status updated to %s by %s", Status, UpdatedBy))
	if err != nil {
		return deal, err
	}

	err = AddNotification(int(ID), deal.QuotationUpdaterID, fmt.Sprintf("Deal status updated to %s by %s", Status, UpdatedBy))
	if err != nil {
		return deal, err
	}

	err = AddNotification(int(ID), deal.DealOwnerID, fmt.Sprintf("Deal status updated to %s by %s", Status, UpdatedBy))
	if err != nil {
		return deal, err
	}
	audit, err := InsertAudit("Update", UpdatedBy, int64(ID), "127.0.0.1", "Status updated in Deal")
	if err != nil {
		return deal, err
	}

	if Status == "Closed Won" {
		UpdateFinalStatusQuery := `UPDATE Deals SET FinalOutCome = ?, UpdatedBy = ?, UpdatedAt = Now(), DealsEndDate = Now() WHERE ID = ?`
		_, err := db.Query(UpdateFinalStatusQuery, "WON", UpdatedBy, ID)
		if err != nil {
			return deal, err
		}
	} else if Status == "Closed Loss" {
		UpdateFinalStatusQuery := `UPDATE Deals SET FinalOutCome = ?, UpdatedBy = ?, UpdatedAt = Now(), DealsEndDate = Now() WHERE ID = ?`
		_, err := db.Query(UpdateFinalStatusQuery, "Loss", UpdatedBy, ID)
		if err != nil {
			return deal, err
		}
	}
	fmt.Print(audit)
	return deal, nil
}

func DenyDeal(ID int, Reason, DeletedBy string) error {
	FinalOutCome := "Invalid"
	Status := "Deal Rejected"
	query := `UPDATE Deals SET IsDeleted = ?, UpdatedBy = ?, UpdatedAt = Now(), DeletedReason = ?, Status = ?, FinalOutCome = ? WHERE ID = ? `
	_, err := db.Query(query, 1, DeletedBy, Reason, Status, FinalOutCome, ID)
	if err != nil {
		return err
	}
	summaryTime := time.Now().Format("02 Jan 2006")
	err = DealSummary(ID, fmt.Sprintf("The deal has been denied by %s at %s, because of %s", DeletedBy, summaryTime, Reason), DeletedBy)
	if err != nil {
		return err
	}
	err = AuditFunction("Update", DeletedBy, int64(ID), "127.0.0.1", "Deal Denied")
	if err != nil {
		return err
	}
	return nil
}

func GetQuotationUploadedData() ([]DealsStruct, error) {
	var deals []DealsStruct
	query := `SELECT ID, IsCompany, CompanyID, ProductType, Quantity, DealOwner, DealApprover, QuotationUpdater, CreatedAt, Status,ApproximateClosingDate, Priority FROM Deals WHERE Status = ? AND IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, "Quotation Sent", 0)
	if err != nil {
		return deals, err
	}
	defer rows.Close()
	for rows.Next() {
		var deal DealsStruct
		err := rows.Scan(&deal.ID, &deal.IsCompany, &deal.CompanyID, &deal.ProductType, &deal.Quantity, &deal.DealOwner, &deal.DealApprover, &deal.QuotationUpdaterName, &deal.CreatedAt, &deal.Status, &deal.ApproximateEndDate, &deal.Priority)
		if err != nil {
			return deals, err
		}
		if deal.IsCompany == 0 {
			getCompanyName := `SELECT ID, Name FROM contacts WHERE ID = ? `
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, fmt.Errorf("error occurred while fetching contacts details %s", err)
			}

		} else if deal.IsCompany == 1 {
			getCompanyName := `SELECT ID, CompanyName FROM company WHERE ID = ? `
			err := db.QueryRow(getCompanyName, deal.CompanyID).Scan(&deal.ContactID, &deal.CompanyName)
			if err != nil {
				return nil, err
			}
		} else {
			deal.CompanyName = "Undefine"
			deal.ContactID = 0
		}
		deals = append(deals, deal)
	}
	return deals, nil
}

func DealSummary(DealID int, Description, CreatedBy string) error {
	query := `INSERT INTO DealSummary(DealID, Description, CreatedBy, CreatedAt) VALUES (?, ?, ?, Now())`
	_, err := db.Query(query, DealID, Description, CreatedBy)
	if err != nil {
		return err
	}
	return nil
}

func GetDealData(ID int) ([]DealSummaryStruct, error) {
	var DealSummary []DealSummaryStruct
	query := `SELECT ID, Description, CreatedBy FROM DealSummary WHERE DealID = ? ORDER BY CreatedAt ASC`
	rows, err := db.Query(query, ID)
	if err != nil {
		return DealSummary, err
	}
	defer rows.Close()
	for rows.Next() {
		var summary DealSummaryStruct
		err := rows.Scan(&summary.ID, &summary.Description, &summary.ByPerson)
		if err != nil {
			return DealSummary, err
		}
		DealSummary = append(DealSummary, summary)
	}

	return DealSummary, nil
}
