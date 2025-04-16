package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ConvertToMySQLDate(dateStr string) (string, error) {
	parsedDate, err := time.Parse("02/01/2006", dateStr)
	if err != nil {
		return "", err
	}
	return parsedDate.Format("2006-01-02"), nil
}

func AddPartner(PartnerName, AuthCode, AreaCovered, CreatedBy, Address, RespectiveSalesPerson, DueDate, RegisteredDate, Email, Phone string, Contacts []GemContact) (string, error) {
	dueDateMySQL, err := ConvertToMySQLDate(DueDate)
	if err != nil {
		return "", fmt.Errorf("invalid DueDate format: %v", err)
	}

	registeredDateMySQL, err := ConvertToMySQLDate(RegisteredDate)
	if err != nil {
		return "", fmt.Errorf("invalid RegisteredDate format: %v", err)
	}

	// Now insert into the database using the correct date format
	query := `INSERT INTO Partners(PartnerName, Address, AuthCode, AreaCovered, RespectiveSalesPerson, DueDate, RegisteredDate, Email, Mobile, LastActivityDate, CreatedBy, CreatedAt, IsDeleted) 
						  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, Now(), ?, Now(), ?)`
	res, err := db.Exec(query, PartnerName, Address, AuthCode, AreaCovered, RespectiveSalesPerson, dueDateMySQL, registeredDateMySQL, Email, Phone, CreatedBy, 0)
	if err != nil {
		return "", err
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	fmt.Println("Contacts Information:")
	for _, contact := range Contacts {
		fmt.Printf("Name: %s, Email: %s, Phone: %s\n", contact.Name, contact.Email, contact.Phone)
		InsertQuery := `INSERT INTO ContactPerson(PartnerID, Name, Mobile, Email, CreatedBy, CreatedAt) VALUES(?, ?, ?, ?, ?, Now())`
		_, err = db.Exec(InsertQuery, ID, contact.Name, contact.Phone, contact.Email, CreatedBy)
		if err != nil {
			return "", err
		}
	}

	err = AuditFunction("Create", CreatedBy, ID, "127.0.0.1", "New Partner created successfully")
	if err != nil {
		return "", err
	}
	return "", nil
}

func GetPartnersData() ([]PartnerStruct, error) {
	var Partners []PartnerStruct
	query := `SELECT ID, PartnerName, AuthCode, AreaCovered, CreatedBy, CreatedAt FROM Partners WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return Partners, err
	}
	defer rows.Close()
	for rows.Next() {
		var Partner PartnerStruct
		err = rows.Scan(&Partner.ID, &Partner.PartnerName, &Partner.AuthCode, &Partner.AreaCovered, &Partner.CreatedBy, &Partner.CreatedAt)
		if err != nil {
			return Partners, err
		}
		Partners = append(Partners, Partner)
	}
	return Partners, nil
}

func GetPartnerListData() ([]PartnerStruct, error) {
	var Partners []PartnerStruct
	var Address sql.NullString
	query := `SELECT ID, PartnerName, Address FROM Partners WHERE IsDeleted = ? 
	ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return Partners, err
	}
	defer rows.Close()
	for rows.Next() {
		var Partner PartnerStruct
		err = rows.Scan(&Partner.ID, &Partner.PartnerName, &Address)
		if err != nil {
			return Partners, err
		}
		if Address.Valid {
			Partner.Address = Address.String
		} else {
			Partner.Address = "N/A"
		}
		Partners = append(Partners, Partner)
	}
	return Partners, nil
}

func AddContactPerson(ID int, Name, Email, Phone, CreatedBy string) (string, error) {
	var ExistID int

	query := `SELECT ID 
				FROM ContactPerson 
				WHERE PartnerID = ? 
				AND (Mobile = ? OR Email = ?)
				AND IsDeleted = ?`

	err := db.QueryRow(query, ID, Phone, Email, 0).Scan(&ExistID)
	if err != nil && err != sql.ErrNoRows {
		return "", err
	}

	if ExistID > 0 {
		return "exist", nil
	}

	query = `INSERT INTO ContactPerson (PartnerID, Name, Mobile, Email, CreatedBy, CreatedAt, IsDeleted) 
			 VALUES (?, ?, ?, ?, ?, NOW(), ?)`

	res, err := db.Exec(query, ID, Name, Phone, Email, CreatedBy, 0)
	if err != nil {
		return "", err
	}

	PersonID, err := res.LastInsertId()
	if err != nil {
		return "", err
	}

	err = AuditFunction("Create", CreatedBy, PersonID, "127.0.0.1", "Contact Person data added successfully")
	if err != nil {
		return "", err
	}

	return "success", nil
}

func GetDataContactPerson() ([]ContactPerson, error) {
	var ContactPersonList []ContactPerson
	query := `SELECT ID, PartnerID, Name, Mobile, Email, CreatedBy, CreatedAt
			FROM ContactPerson
			WHERE IsDeleted = ? 
			ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var ContactPerson ContactPerson
		err := rows.Scan(&ContactPerson.ID, &ContactPerson.PartnerID, &ContactPerson.ContactName, &ContactPerson.ContactPhone, &ContactPerson.ContactEmail, &ContactPerson.CreatedBy, &ContactPerson.CreatedAt)
		if err != nil {
			return nil, err
		}
		GetQuery := `SELECT PartnerName FROM Partners WHERE ID = ?`
		err = db.QueryRow(GetQuery, ContactPerson.PartnerID).Scan(&ContactPerson.PartnerName)
		if err != nil {
			return nil, err
		}
		ContactPersonList = append(ContactPersonList, ContactPerson)
	}
	return ContactPersonList, nil
}

func AddGem(DealerID, OwnerID int, GemNumber, DueDate string, Quantity int, DealerName, Customer, Description, GemRaisedBy, Status, gemCertificateStr, Link, Team, GemUrl, GemDate, ProductType string) (MailStruct, error) {
	var data MailStruct
	var lastTicketNo string
	err := db.QueryRow("SELECT TicketNo FROM GeM ORDER BY ID DESC LIMIT 1").Scan(&lastTicketNo)
	if err != nil && err != sql.ErrNoRows {
		return data, fmt.Errorf("failed to get last TicketNo: %v", err)
	}

	newTicketNo := "HG00001"
	if lastTicketNo != "" {
		lastNumberStr := lastTicketNo[2:]
		lastNumber, err := strconv.Atoi(lastNumberStr)
		if err != nil {
			return data, fmt.Errorf("failed to parse last TicketNo: %v", err)
		}

		newNumber := lastNumber + 1
		newTicketNo = fmt.Sprintf("HG%05d", newNumber) // Ensures at least 3 digits (e.g., HD001, HD010, HD100)
	}
	Level := ""
	if Quantity >= 1 && Quantity < 25 {
		Level = "Bronze"
	} else if Quantity >= 25 && Quantity < 50 {
		Level = "Silver"
	} else if Quantity >= 50 {
		Level = "Gold"
	}

	query := `INSERT INTO GeM (GemNumber, DueDate, Quantity, DealerName, DealerID, CustomerName, Description, GemCertificate, GemRaisedBy, Status, IsApproved, IsDeleted, IsQuotationUpdated, Level, Link, GemRaisedAt, OwnerID, GemUrl, GemDate, ProductType, TicketNo, GemFinalURL)
						VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, Now(), ?, ?, ?, ?, ?, ?)`
	res, err := db.Exec(query, GemNumber, DueDate, Quantity, DealerName, DealerID, Customer, Description, gemCertificateStr, GemRaisedBy, Status, 0, 0, 0, Level, Link, OwnerID, GemUrl, GemDate, ProductType, newTicketNo, "N/A")
	if err != nil {
		fmt.Print(err)
		return data, err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		fmt.Print(err)
		return data, err
	}
	LastActivityQuery := `UPDATE Partners SET LastActivityDate = Now() WHERE ID = ?`
	_, err = db.Exec(LastActivityQuery, DealerID)
	if err != nil {
		fmt.Print(err)
		return data, err
	}
	err = AddGemNotification(OwnerID, int(ID), fmt.Sprintf("New GEM raised by %s and now it was waiting for Pricing team approval", GemRaisedBy))
	if err != nil {
		fmt.Print(err)
		return data, err
	}
	GetManagerQuery := `SELECT ID FROM users WHERE (Team = ? AND Role = ?) AND IsDeleted = ? AND IsApproved = ?`

	result, err := db.Query(GetManagerQuery, "Pricing", "Manager", 0, 1)
	if err != nil {
		return data, err
	}

	for result.Next() {
		var id int
		err := result.Scan(&id)
		if err != nil {
			return data, err
		}

		err = AddGemNotification(id, int(ID), fmt.Sprintf("New GEM Created by %s, Now It's Waiting for your approval.", GemRaisedBy))
		if err != nil {
			fmt.Print("Deal was not created")
			return data, nil
		}
	}
	GetAdminQuery := `SELECT ID FROM users WHERE (Role = ? OR Role = ?) AND IsDeleted = ? AND IsApproved = ?`

	adminres, err := db.Query(GetAdminQuery, "Admin", "Super Admin", 0, 1)
	if err != nil {
		fmt.Print("While executing admin query", err)
		return data, err
	}

	for adminres.Next() {
		var id int
		err := adminres.Scan(&id)
		if err != nil {
			fmt.Print("While scanning admin query", err)
			return data, err
		}

		err = AddGemNotification(id, int(ID), fmt.Sprintf("New GEM has been Created by %s, Now It's Waiting for the Pricing team approval.", GemRaisedBy))
		if err != nil {
			fmt.Print("Deal was not created")
			return data, nil
		}
	}
	err = AuditFunction("Create", GemRaisedBy, ID, "127.0.0.1", "New GEM created Successfully")
	if err != nil {
		fmt.Print(err)
		return data, err
	}
	var To []string
	var From string

	GetEmailQuery := `SELECT Email FROM users WHERE ID = ?`
	err = db.QueryRow(GetEmailQuery, OwnerID).Scan(&From)
	if err != nil {
		return data, err
	}

	if Team == "Sales" {
		To = []string{"pricing@holoware.co", "sathish.pk@holoware.co"}
	} else if Team == "Marketing" {
		To = []string{"marketing@holoware.co", "pricing@holoware.co"}
	} else {
		To = []string{"pricing@holoware.co"}
	}

	currentDate := time.Now().Format("02 Jan 2006")
	Subject := fmt.Sprintf("New GeM Created for %s on %s", DealerName, currentDate)
	Body := fmt.Sprintf(`
		<html>
		<body>
			<p>Hello Team,</p>
			<p>A new GeM has been created by <strong>%s</strong> and is currently awaiting admin approval.</p>
			<p><strong>Details of the GeM:</strong></p>
			<ul>
				<li><strong>GeM Number:</strong> %s</li>
				<li><strong>Quantity:</strong> %d</li>
				<li><strong>Due Date:</strong> %s</li>
				<li><strong>Customer Name:</strong> %s</li>
			</ul>
			<p>Thanks & Regards,<br>Team CRM</p>
		</body>
		</html>
		`,
		GemRaisedBy, GemNumber, Quantity, DueDate, Customer,
	)
	data.To = To
	data.From = From
	data.Subject = Subject
	data.Body = Body
	return data, nil
}

func GetGemWaitingForApproval() ([]GemStructure, error) {
	var Gems []GemStructure
	query := `SELECT ID, GemNumber, DueDate, Quantity, DealerName,
	CustomerName, Description, GemRaisedBy, GemRaisedAt, TicketNo FROM GeM WHERE IsApproved = ? AND IsDeleted = ? ORDER BY ID DESC`
	rows, err := db.Query(query, 0, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var Gem GemStructure
		err := rows.Scan(&Gem.ID, &Gem.GemNumber, &Gem.DueDate, &Gem.Quantity, &Gem.DealerName, &Gem.Customer, &Gem.Description, &Gem.GemRaisedBy, &Gem.GemRaisedAt, &Gem.TicketNo)
		if err != nil {
			return nil, err
		}
		Gems = append(Gems, Gem)
	}
	return Gems, nil
}

func GetGem(ID int) (GemStructure, error) {
	var gem GemStructure
	var gemCertificateStr sql.NullString
	var QuotationStr, GemDate, ProductType sql.NullString
	query := `SELECT ID, DealerID, GemNumber, DueDate, Quantity, DealerName, CustomerName, Description, GemCertificate, Quotation, GemRaisedBy, Status, GemRaisedAt, IsApproved, GemUrl, GemDate, ProductType, TicketNo, OwnerID FROM GeM WHERE ID = ?`
	err := db.QueryRow(query, ID).Scan(
		&gem.ID, &gem.DealerID, &gem.GemNumber, &gem.DueDate, &gem.Quantity,
		&gem.DealerName, &gem.Customer, &gem.Description,
		&gemCertificateStr, &QuotationStr, &gem.GemRaisedBy, &gem.Status, &gem.GemRaisedAt, &gem.IsApproved, &gem.GemUrl, &GemDate, &ProductType, &gem.TicketNo, &gem.OwnerID,
	)
	if err != nil {
		return gem, err
	}
	GetAddress := `SELECT Address, Email, Mobile FROM Partners WHERE ID = ?`
	err = db.QueryRow(GetAddress, gem.DealerID).Scan(&gem.DealerAddress, &gem.DealerEmail, &gem.DealerMobile)
	if err != nil {
		return gem, err
	}
	fmt.Printf("Owner ID is %d", gem.OwnerID)
	GetEmailQuery := `SELECT Email FROM users WHERE ID = ?`
	err = db.QueryRow(GetEmailQuery, gem.OwnerID).Scan(&gem.DeaOwnerEmail)
	if err != nil {
		if err == sql.ErrNoRows {
			return gem, nil
		}
		return gem, err
	}
	if GemDate.Valid {
		gem.GemDate = GemDate.String
	} else {
		gem.GemDate = "N/A"
	}
	if ProductType.Valid {
		gem.ProductType = ProductType.String
	} else {
		gem.ProductType = "N/A"
	}
	if gemCertificateStr.Valid {
		gem.GemCertificate = strings.Split(gemCertificateStr.String, ",")
	} else {
		gem.GemCertificate = []string{}
	}
	if QuotationStr.Valid {
		gem.Quotation = strings.Split(QuotationStr.String, ",")
		fmt.Print(gem.Quotation)
	} else {
		gem.Quotation = []string{}
	}
	return gem, nil
}

func GetAwaitingPricingGEMList() ([]GemStructure, error) {
	var GemList []GemStructure
	Status := "Gem Requested"
	query := `SELECT ID, GemNumber, DueDate, Quantity, DealerName, CustomerName, Description, GemCertificate, GemRaisedBy, Status, GemRaisedAt, TicketNo FROM GeM WHERE IsApproved = ? AND IsDeleted = ? AND Status = ? ORDER BY ID DESC`
	rows, err := db.Query(query, 1, 0, Status)
	if err != nil {
		return GemList, err
	}
	defer rows.Close()
	for rows.Next() {
		var gem GemStructure
		var gemCertificateStr string
		err = rows.Scan(
			&gem.ID, &gem.GemNumber, &gem.DueDate, &gem.Quantity,
			&gem.DealerName, &gem.Customer, &gem.Description,
			&gemCertificateStr, &gem.GemRaisedBy, &gem.Status, &gem.GemRaisedAt, &gem.TicketNo,
		)
		if err != nil {
			return GemList, err
		}
		gem.GemCertificate = strings.Split(gemCertificateStr, ",")
		GemList = append(GemList, gem)
	}
	return GemList, nil
}

func GetUpdatedPricingGEMList() ([]GemStructure, error) {
	var GemList []GemStructure
	query := `SELECT ID, GemNumber, DueDate, Quantity, DealerName, CustomerName, Description, GemCertificate, GemRaisedBy, Status, GemRaisedAt, TicketNo FROM GeM WHERE IsApproved = ? AND IsDeleted = ? AND IsQuotationUpdated = ? ORDER BY ID DESC`
	rows, err := db.Query(query, 1, 0, 1)
	if err != nil {
		return GemList, err
	}
	defer rows.Close()
	for rows.Next() {
		var gem GemStructure
		var gemCertificateStr string
		err = rows.Scan(
			&gem.ID, &gem.GemNumber, &gem.DueDate, &gem.Quantity,
			&gem.DealerName, &gem.Customer, &gem.Description,
			&gemCertificateStr, &gem.GemRaisedBy, &gem.Status, &gem.GemRaisedAt, &gem.TicketNo,
		)
		if err != nil {
			return GemList, err
		}
		gem.GemCertificate = strings.Split(gemCertificateStr, ",")
		GemList = append(GemList, gem)
	}
	return GemList, nil
}

func UploadQuotationGem(ID int, Quotation, UpdatedBy, Status, GEMFinalURL string) error {
	fmt.Printf("Quotation %s ID %d UpdatedBy %s Status %s", Quotation, ID, UpdatedBy, Status)
	query := `UPDATE GeM SET Quotation = ?, IsQuotationUpdated = ?, UpdatedBy = ?, UpdatadAt = Now(), Status = ?, GemFinalURL = ? WHERE ID = ?`
	_, err := db.Exec(query, Quotation, 1, "Arun AK", Status, GEMFinalURL, ID)
	if err != nil {
		fmt.Print("while executing", err)
		return err
	}

	// var DealerID int
	// GetDealerID := `SELECT DealerID FROM Gem WHERE ID = ? `
	// err = db.QueryRow(GetDealerID, ID).Scan(&DealerID)
	// if err != nil {
	// 	return err
	// }
	// UpdateLastActivity := `UPDATE LastActivityDate SET Partners WHERE ID = ?`
	// _, err = db.Exec(UpdateLastActivity, DealerID)
	// if err != nil {
	// 	return err
	// }

	GetDataQuery := `SELECT GemRaisedBy, OwnerID, ApprovedBy, ApproverID, QuotationUpdaterID FROM GeM WHERE ID = ?`
	var OwnerID, ApproverID, QuotationUpdaterID int
	var GemRaisedBy, ApprovedBy string

	err = db.QueryRow(GetDataQuery, ID).Scan(&GemRaisedBy, &OwnerID, &ApprovedBy, &ApproverID, &QuotationUpdaterID)
	if err != nil {
		return err
	}

	// err = AddGemNotification(OwnerID, ID, fmt.Sprintf("The documents for your GEM were updated by %s ", UpdatedBy))
	// if err != nil {
	// 	fmt.Print("while executing", err)
	// 	return err
	// }

	// err = AddGemNotification(ApproverID, ID, fmt.Sprintf("The documents for your GEM were updated by %s ", UpdatedBy))
	// if err != nil {
	// 	fmt.Print("while executing", err)
	// 	return err
	// }

	// err = AddGemNotification(QuotationUpdaterID, ID, "You have successfully updated the documents for this GEM")
	// if err != nil {
	// 	return nil
	// }

	err = AuditFunction("Update", UpdatedBy, int64(ID), "127.0.0.1", "GEM Quotation updated Successfully")
	if err != nil {
		fmt.Print("while audit", err)
		return err
	}
	fmt.Print("Audit Data Updated Successfully")
	return nil
}

func ApproveGeM(Data GeMApprove) error {
	Status := "Gem Requested"
	query := `UPDATE GeM SET IsApproved = ?, ApprovedBy = ?, QuotationUpdatedBy = ?, QuotationUpdaterID = ?, ApproverID = ?, ApprovedAt = Now(), Status = ? WHERE ID = ?`
	_, err := db.Exec(query, 1, Data.GemApprover, Data.QuotationApprover, Data.QuotationApproverID, Data.GemApproverID, Status, Data.ID)
	if err != nil {
		return err
	}
	GetDataQuery := `SELECT GemRaisedBy, OwnerID FROM GeM WHERE ID = ?`
	var OwnerID int
	var GemRaisedBy string

	err = db.QueryRow(GetDataQuery, Data.ID).Scan(&GemRaisedBy, &OwnerID)
	if err != nil {
		return err
	}
	err = AuditFunction("Update", Data.GemApprover, int64(Data.ID), "127.0.0.1", "GEM Approved Successfully")
	if err != nil {
		return err
	}

	err = AddGemNotification(Data.GemApproverID, Data.ID, fmt.Sprintf("You have successfully assigned the GEM document to %s", Data.QuotationApprover))
	if err != nil {
		return err
	}
	err = AddGemNotification(Data.QuotationApproverID, Data.ID, fmt.Sprintf("%s has assigned you a GEM from %s ", Data.GemApprover, GemRaisedBy))
	if err != nil {
		return err
	}
	err = AddGemNotification(OwnerID, Data.ID, fmt.Sprintf("Your GEM has been approved by %s and has now been assigned to %s", Data.GemApprover, Data.QuotationApprover))
	if err != nil {
		return err
	}
	//GetMail := `SELECT Email FROM users WHERE ID = ?`
	//err = db.QueryRow(GetMail,)
	return nil
}

func DenyGeM(Data GeMDeny) error {
	Status := "Denied"
	query := `UPDATE GeM SET IsApproved = ?, IsDeleted = ?, ApprovedBy = ?, ApproverID = ?, Reason = ?, GemEndDate = Now(), Status = ? WHERE ID = ?`
	_, err := db.Exec(query, 0, 1, Data.GemRemover, Data.GemRemoverID, Data.Reason, Status, Data.ID)
	if err != nil {
		return err
	}
	err = AuditFunction("Update", Data.GemRemover, int64(Data.ID), "127.0.0.1", "GEM Denied Successfully")
	if err != nil {
		return err
	}
	GetDataQuery := `SELECT GemRaisedBy, OwnerID FROM GeM WHERE ID = ?`
	var OwnerID int
	var GemRaisedBy string
	err = db.QueryRow(GetDataQuery, Data.ID).Scan(&GemRaisedBy, &OwnerID)
	if err != nil {
		return err
	}

	err = AddGemNotification(Data.GemRemoverID, Data.ID, fmt.Sprintf("This GEM document has been Rejected by %s", Data.GemRemover))
	if err != nil {
		return err
	}
	err = AddGemNotification(OwnerID, Data.ID, fmt.Sprintf("Your GEM document has been Rejected by %s", Data.GemRemover))
	if err != nil {
		return err
	}
	return nil
}

func GetMyGem(Data AllGeM) ([]GemStructure, error) {
	var Gems []GemStructure
	var QuotationUpdaterID, ApproverID, OwnerID sql.NullInt16
	var query string
	var res *sql.Rows
	var err error
	var ApprovedBy sql.NullString

	if Data.Team == "Admin" || Data.Team == "Pricing" {
		query = `SELECT ID, GemNUmber, DueDate, Quantity, DealerName, CustomerName, Description, 
		GemRaisedBy, GemRaisedAt, IsApproved, ApprovedBy, IsQuotationUpdated, Level, OwnerID, ApproverID, QuotationUpdaterID, TicketNo,DealerID FROM GeM WHERE IsDeleted = 0 ORDER BY GemRaisedAt Desc`
		res, err = db.Query(query)
	} else {
		query = `
		SELECT 
			ID, GemNUmber, DueDate, Quantity, DealerName, CustomerName, Description, GemRaisedBy, 
			GemRaisedAt, IsApproved, ApprovedBy, IsQuotationUpdated, Level, OwnerID, ApproverID, QuotationUpdaterID, TicketNo,DealerID
		FROM 
			GeM
		WHERE 
			IsDeleted = 0 AND OwnerID = ? OR (OwnerID IN (SELECT ID FROM users WHERE Team = ? AND Role = 'Manager'))
		ORDER BY 
			GemRaisedAt DESC
		`
		res, err = db.Query(query, Data.UserID, Data.Team)
	}
	if err != nil {
		fmt.Print(err)
		return Gems, err
	}
	defer res.Close()
	for res.Next() {
		var Gem GemStructure
		err := res.Scan(
			&Gem.ID, &Gem.GemNumber, &Gem.DueDate, &Gem.Quantity, &Gem.DealerName,
			&Gem.Customer, &Gem.Description, &Gem.GemRaisedBy, &Gem.GemRaisedAt,
			&Gem.IsApproved, &ApprovedBy, &Gem.IsQuotationUpdated, &Gem.Level, &OwnerID, &ApproverID, &QuotationUpdaterID, &Gem.TicketNo, &Gem.DealerID,
		)
		if err != nil {
			fmt.Print(err)
			return Gems, err
		}

		if QuotationUpdaterID.Valid {
			Gem.QuotationUpdaterID = int(QuotationUpdaterID.Int16)
		} else {
			Gem.QuotationUpdaterID = 0
		}

		if ApproverID.Valid {
			Gem.ApproverID = int(ApproverID.Int16)
		} else {
			Gem.ApproverID = 0
		}
		if OwnerID.Valid {
			Gem.OwnerID = int(OwnerID.Int16)
		} else {
			Gem.OwnerID = 0
		}
		if ApprovedBy.Valid {
			Gem.ApprovedBy = ApprovedBy.String
		}
		if Gem.IsApproved == 1 && Gem.IsQuotationUpdated == 0 {
			Gem.Status = "GeM Approved Waiting for Quotation"
		} else if Gem.IsApproved == 0 {
			Gem.Status = "Waiting for Approval"
		} else if Gem.IsQuotationUpdated == 1 {
			Gem.Status = "Quotation Updated"
		}
		Gems = append(Gems, Gem)
	}
	return Gems, nil
}
func UpdateGemStatus(Data UpdateStatus) error {
	var EndDate bool
	var query string
	if Data.Status == "Won" || Data.Status == "Loss" {
		EndDate = true
	}
	if EndDate {
		query = `UPDATE GeM SET Status = ?, UpdatedBy = ?, UpdatadAt = Now(), GemEndDate = Now() WHERE ID = ?`
	} else {
		query = `UPDATE GeM SET Status = ?, UpdatedBy = ?, UpdatadAt = Now() WHERE ID = ?`
	}
	_, err := db.Exec(query, Data.Status, Data.UpdatedBy, Data.ID)
	if err != nil {
		return err
	}
	err = AuditFunction("Update", Data.UpdatedBy, int64(Data.ID), "127.0.0.1", "This GeM status has been updated")
	if err != nil {
		return err
	}
	return nil
}

func GetPartner(ID int) (DetailedPartnerData, error) {
	var data DetailedPartnerData

	// Query for partner data
	Query := `SELECT ID, PartnerName, AuthCode, AreaCovered, CreatedBy, Address, RespectiveSalesPerson, DueDate, RegisteredDate, Email, Mobile, LastActivityDate FROM Partners WHERE ID = ?`
	err := db.QueryRow(Query, ID).Scan(
		&data.PartnerID, &data.PartnerName, &data.AuthCode, &data.AreaCovered,
		&data.PartnerCreatedBy, &data.Address, &data.RespectiveSalesPerson,
		&data.DueDate, &data.RegisteredDate, &data.Email, &data.Phone, &data.LastActivityDate,
	)
	if err != nil {
		fmt.Print("Error while partner", err)
		return data, err
	}

	// Query for contact data
	ContactQuery := `SELECT ID, Name, Mobile, Email, CreatedBy, CreatedAt FROM ContactPerson WHERE PartnerID = ?`
	ContactRes, err := db.Query(ContactQuery, ID)
	if err != nil {
		fmt.Print("Error while Contact", err)
		return data, err
	}
	defer ContactRes.Close()

	// Loop through contacts and append them to the Contacts slice in DetailedPartnerData
	for ContactRes.Next() {
		var contact GemContact
		err := ContactRes.Scan(
			&contact.ContactID, &contact.Name,
			&contact.Phone, &contact.Email, &contact.CreatedBy, &contact.CreatedAt,
		)
		if err != nil {
			fmt.Print("Error while Contact Scanning", err)
			return data, err
		}
		data.Contacts = append(data.Contacts, contact)
	}
	var Reason, GemEndDate sql.NullString

	GemQuery := `SELECT ID, GemNumber, DueDate, Quantity, DealerName, CustomerName, Description, GemRaisedAt, Status, Level, Reason, GemEndDate FROM GeM WHERE DealerID = ? ORDER BY GemRaisedAt DESC`
	GemRes, err := db.Query(GemQuery, ID)
	if err != nil {
		fmt.Print("Error while Gem", err)
		return data, err
	}
	defer GemRes.Close()
	for GemRes.Next() {
		var gem GetProductList
		err := GemRes.Scan(
			&gem.ID, &gem.GemNumber, &gem.DueDate, &gem.Quantity, &gem.DealerName,
			&gem.Customer, &gem.Description, &gem.GemRaisedAt, &gem.Status, &gem.Level, &Reason, &GemEndDate,
		)
		if err != nil {
			fmt.Print("Error while gem scanning", err)
			return data, err
		}
		if Reason.Valid {
			gem.Reason = Reason.String
		} else {
			gem.Reason = "N/A"
		}
		if GemEndDate.Valid {
			gem.GemEndDate = GemEndDate.String
		} else {
			gem.GemEndDate = "N/A"
		}
		data.Gem = append(data.Gem, gem)
	}
	// Return the filled data
	return data, nil
}

// SELECT
// 			ID, GemNUmber, DueDate, Quantity, DealerName, CustomerName, Description, GemRaisedBy,
// 			GemRaisedAt, IsApproved, ApprovedBy, IsQuotationUpdated, Level, OwnerID, ApproverID, QuotationUpdaterID
// 		FROM
// 			GeM
// 		WHERE
// 			IsDeleted = 0 AND OwnerID = 17 OR (OwnerID IN (SELECT ID FROM users WHERE Team = "Sales" AND Role = 'Manager'))
// 		ORDER BY
// 			GemRaisedAt DESC
