package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

type SalesPersons struct {
	ID          int    `json:"id"`
	SalesPerson string `json:"salesPerson"`
}

type DataSheetData struct {
	ID               int    `json:"id"`
	ChipSet          string `json:"chipSet"`
	Memory           string `json:"memory"`
	MemorySlots      string `json:"memorySlots"`
	Storage          string `json:"storage"`
	StorageSupport   string `json:"storageSupport"`
	StorageInterface string `json:"storageInterface"`
	Graphics         string `json:"graphics"`
	Audio            string `json:"audio"`
	Ports            string `json:"ports"`
	LAN              string `json:"lan"`
	PCI              string `json:"pci"`
	WLANSupport      string `json:"wlanSupport"`
	Monitor          string `json:"monitor"`
	Power            string `json:"power"`
	Wifi             string `json:"wifiBt"`
	Security         string `json:"security"`
	Dimension        string `json:"dimensions"`
	Weight           string `json:"weight"`
	OS               string `json:"availableOS"`
	Warranty         string `json:"warranty"`
	Camera           string `json:"camera"`
	Battery          string `json:"battery"`
	PowerAdapter     string `json:"powerAdaptor"`
	CreatedBy        string `json:"createdBy"`
	CreatedAt        string `json:"createdAt"`
	UpdatedBy        string `json:"updatedBy"`
	UpdatedAt        string `json:"updatedAt"`
	CabinetColor     string `json:"cabinetColor"`
	CabinetLitre     string `json:"cabinetLitres"`
	Sensor           string `json:"sensor"`
	Connectivity     string `json:"connectivity"`
	TouchPanel       string `json:"touchPanel"`
	Model            string `json:"laptopModel"`
}

type ProcessorChipsetData struct {
	ChipSet   string `json:"chipSet"`
	Processor string `json:"processor"`
	Product   string `json:"product"`
}

type ProcessorData struct {
	ID              int    `json:"id"`
	Processor       string `json:"processor"`
	Graphics        string `json:"graphics"`
	GraphicsSupport string `json:"graphicsSupport"`
	MonitorSupport  string `json:"monitorSupport"`
	Screen          string `json:"screen"`
	CreatedBy       string `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
}

type DataSheetName struct {
	ID      int    `json:"id"`
	ChipSet string `json:"chipSet"`
	Product string `json:"product"`
}
type GemNotificationStruct struct {
	ID        int    `json:"id"`
	GemID     int    `json:"gemID"`
	UserID    int    `json:"userID"`
	Message   string `json:"message"`
	IsRead    int    `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}
type Notifications struct {
	ID        int    `json:"id"`
	DealID    int    `json:"dealID"`
	UserID    int    `json:"userID"`
	Message   string `json:"message"`
	IsRead    int    `json:"isRead"`
	CreatedAt string `json:"createdAt"`
}

type ProcessorStruct struct {
	ID              int    `json:"id"`
	Processor       string `json:"processor"`
	Graphics        string `json:"graphics"`
	GraphicsSupport string `json:"graphicsSupport"`
	MonitorSupport  string `json:"monitorSupport"`
	Screen          string `json:"screen"`
	CreatedBy       string `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
	UpdatedBy       string `json:"updatedBy"`
	UpdatedAt       string `json:"updatedAt"`
}

type DataSheetDatas struct {
	ID               int    `json:"id"`
	ChipSet          string `json:"chipSet"`
	Memory           string `json:"memory"`
	MemorySlots      string `json:"memorySlots"`
	Storage          string `json:"storage"`
	StorageSupport   string `json:"storageSupport"`
	StorageInterface string `json:"storageInterface"`
	Graphics         string `json:"graphics"`
	Audio            string `json:"audio"`
	Ports            string `json:"ports"`
	LAN              string `json:"lan"`
	PCI              string `json:"pci"`
	WLANSupport      string `json:"wlanSupport"`
	Monitor          string `json:"monitor"`
	Power            string `json:"power"`
	Wifi             string `json:"wifiBt"`
	Security         string `json:"security"`
	Dimension        string `json:"dimensions"`
	Weight           string `json:"weight"`
	OS               string `json:"availableOS"`
	Warranty         string `json:"warranty"`
	Camera           string `json:"camera"`
	Battery          string `json:"battery"`
	PowerAdapter     string `json:"powerAdaptor"`
	CreatedBy        string `json:"createdBy"`
	CreatedAt        string `json:"createdAt"`
	UpdatedBy        string `json:"updatedBy"`
	UpdatedAt        string `json:"updatedAt"`
	CabinetColor     string `json:"cabinetColor"`
	CabinetLitre     string `json:"cabinetLitres"`
	Sensor           string `json:"sensor"`
	Connectivity     string `json:"connectivity"`
	TouchPanel       string `json:"touchPanel"`
	Model            string `json:"laptopModel"`
}

var redisClient *redis.Client

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
	})
	fmt.Printf("Connecting to Redis system: Host=%s Port=%s Password=%s\n", redisHost, redisPort, redisPassword)

}

func NormalizeChipset(chipset string) string {
	chipset = strings.ToLower(chipset)
	chipset = strings.TrimPrefix(chipset, "intel ")
	chipset = strings.TrimPrefix(chipset, "amd ")

	chipset = strings.ReplaceAll(chipset, " ", "")

	return chipset
}

func AddDataSheet(Data DataSheetData) error {
	var ExistID int

	if Data.ChipSet == "N/A" {
		fmt.Print("model name", Data.Model)
		Data.ChipSet = Data.Model
		fmt.Printf("The data switched, Chipset value is %s", Data.ChipSet)
	}
	CheckQuery := `SELECT ID FROM DataSheet WHERE ChipSet = ? `
	err := db.QueryRow(CheckQuery, Data.ChipSet).Scan(&ExistID)

	if err != nil {
		if err == sql.ErrNoRows {
			Query := `INSERT INTO DataSheet(ChipSet, MemorySlots, StorageSupport, Audio, PortsAndConnectors, LAN, PCI, WifiBt, Power, Dimension, Weight, StorageInterface, CabinetColor, CabinetLitre, WLANSupport, Camera, Battery, PowerAdapter, Sensor, Connectivity, TouchPanel, CreatedBy, CreatedAt, Security)
            VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, Now(), ?)`
			res, err := db.Exec(Query, Data.ChipSet, Data.MemorySlots, Data.StorageSupport, Data.Audio, Data.Ports, Data.LAN, Data.PCI, Data.Wifi, Data.Power, Data.Dimension, Data.Weight, Data.StorageInterface, Data.CabinetColor, Data.CabinetLitre, Data.WLANSupport, Data.Camera, Data.Battery, Data.PowerAdapter, Data.Sensor, Data.Connectivity, Data.TouchPanel, Data.CreatedBy, Data.Security)
			if err != nil {
				return err
			}
			ID, err := res.LastInsertId()
			if err != nil {
				fmt.Print("Data inserted successfully")
				return err
			}
			err = AuditFunction("Create", Data.CreatedBy, ID, "127.0.0.1", "New DataSheet added successfully")
			if err != nil {
				return err
			}
			return nil
		} else {
			fmt.Print("Data exist")
		}
		return err
	}

	return nil
}

func GetDataSheetName() ([]DataSheetName, error) {
	var dataSheetNames []DataSheetName
	query := `SELECT ID, ChipSet FROM DataSheet WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, 0)
	if err != nil {
		return dataSheetNames, err
	}
	defer rows.Close()
	for rows.Next() {
		var data DataSheetName
		err = rows.Scan(&data.ID, &data.ChipSet)
		if err != nil {
			return dataSheetNames, err
		}
		dataSheetNames = append(dataSheetNames, data)
	}
	return dataSheetNames, nil
}

func GetDataSheetData(ID int) (DataSheetData, error) {
	var DataSheet DataSheetData
	var Storage, Power sql.NullString

	query := `SELECT ID, ChipSet, MemorySlots, Storage, Audio, PortsAndConnectors, LAN, PCI, Power, WIFI FROM DataSheet WHERE ID = ? AND IsDeleted = ?`
	err := db.QueryRow(query, ID, 0).Scan(&DataSheet.ID, &DataSheet.ChipSet, &DataSheet.MemorySlots, &Storage, &DataSheet.Audio, &DataSheet.Ports, &DataSheet.LAN, &DataSheet.PCI, &Power, &DataSheet.Wifi)
	if err != nil {
		return DataSheet, err
	}
	if Storage.Valid {
		DataSheet.Storage = Storage.String
	} else {
		DataSheet.Storage = ""
	}
	if Power.Valid {
		DataSheet.Power = Power.String
	} else {
		DataSheet.Power = ""
	}
	return DataSheet, nil
}

func GetDataSheetByChipSet(Data ProcessorChipsetData) (DataSheetDatas, error) {
	var data DataSheetDatas
	fmt.Printf("Model %s", Data.ChipSet)
	query := `SELECT ID, ChipSet, MemorySlots, StorageSupport, StorageInterface, 
			  Audio, PortsAndConnectors, LAN, PCI, WLANSupport, Power, WifiBt, 
			  Security, Dimension, Weight, Camera, Battery, PowerAdapter, CreatedBy, 
			  CreatedAt, CabinetColor, CabinetLitre, Sensor, Connectivity, 
			  TouchPanel 
			  FROM DataSheet WHERE ChipSet = ? AND IsDeleted = ? ORDER BY CreatedAt DESC`

	err := db.QueryRow(query, Data.ChipSet, 0).Scan(
		&data.ID, &data.ChipSet, &data.MemorySlots,
		&data.StorageSupport, &data.StorageInterface, &data.Audio,
		&data.Ports, &data.LAN, &data.PCI, &data.WLANSupport, &data.Power,
		&data.Wifi, &data.Security, &data.Dimension, &data.Weight,
		&data.Camera, &data.Battery, &data.PowerAdapter, &data.CreatedBy, &data.CreatedAt,
		&data.CabinetColor, &data.CabinetLitre, &data.Sensor,
		&data.Connectivity, &data.TouchPanel)
	if err != nil {
		return data, err
	}
	return data, nil
}

func GetDataSheetDatas(Data ProcessorChipsetData) (DataSheetData, ProcessorStruct, error) {
	var ProcessorData ProcessorStruct
	var ChipSetData DataSheetData
	var PowerAdapter, Camera, Battery, StorageSupport, StorageInterface, WLANSupport, Security, Monitor, Power, Dimension, Weight, CabinetColor, CabinetLitre, Sensor, Connectivity, TouchPanel sql.NullString

	fmt.Printf("Chipset data: %s, Processor data: %s, Product type %s\n", Data.ChipSet, Data.Processor, Data.Product)

	// Query for the chipset
	query := `SELECT ID, ChipSet, MemorySlots, StorageSupport, StorageInterface, Audio, PortsAndConnectors, LAN, PCI, Monitor, Power, WIFI, WLANSupport, Security, Camera, Battery, PowerAdapter, Dimension, Weight, CabinetColor, CabinetLitre, Sensor, Connectivity, TouchPanel 
			  FROM DataSheet WHERE ChipSet = ? AND IsDeleted = ?`
	err := db.QueryRow(query, Data.ChipSet, 0).Scan(
		&ChipSetData.ID,
		&ChipSetData.ChipSet,
		&ChipSetData.MemorySlots,
		&StorageSupport,
		&StorageInterface,
		&ChipSetData.Audio,
		&ChipSetData.Ports,
		&ChipSetData.LAN,
		&ChipSetData.PCI,
		&Monitor,
		&Power,
		&ChipSetData.Wifi,
		&WLANSupport,
		&Security,
		&Camera,
		&Battery,
		&PowerAdapter,
		&Dimension,
		&Weight,
		&CabinetColor,
		&CabinetLitre,
		&Sensor,
		&Connectivity,
		&TouchPanel,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			ChipSetData = DataSheetData{}
		} else {
			return ChipSetData, ProcessorData, err
		}
	}

	fmt.Printf("Raw values from DB:\nPower: %v, Monitor: %v, Camera: %v, Battery: %v, PowerAdapter: %v\n",
		Power, Monitor, Camera, Battery, PowerAdapter)

	if Power.Valid && Power.String != "" {
		ChipSetData.Power = Power.String
	} else {
		ChipSetData.Power = "N/A"
	}
	if Monitor.Valid && Monitor.String != "" {
		ChipSetData.Monitor = Monitor.String
	} else {
		ChipSetData.Monitor = "N/A"
	}
	if Camera.Valid && Camera.String != "" {
		ChipSetData.Camera = Camera.String
	} else {
		ChipSetData.Camera = Camera.String
	}
	if Battery.Valid && Battery.String != "" {
		ChipSetData.Battery = Battery.String
	} else {
		ChipSetData.Battery = "N/A"
	}
	if PowerAdapter.Valid && PowerAdapter.String != "" {
		ChipSetData.PowerAdapter = PowerAdapter.String
	} else {
		ChipSetData.PowerAdapter = "N/A"
	}
	if StorageSupport.Valid && StorageSupport.String != "" {
		ChipSetData.StorageSupport = StorageSupport.String
	} else {
		ChipSetData.StorageSupport = "N/A"
	}
	if StorageInterface.Valid && StorageInterface.String != "" {
		ChipSetData.StorageInterface = StorageInterface.String
	} else {
		ChipSetData.StorageInterface = "N/A"
	}
	if WLANSupport.Valid && WLANSupport.String != "" {
		ChipSetData.WLANSupport = WLANSupport.String
	} else {
		ChipSetData.WLANSupport = "N/A"
	}
	if Security.Valid && Security.String != "" {
		ChipSetData.Security = Security.String
	} else {
		ChipSetData.Security = "N/A"
	}
	if Dimension.Valid && Dimension.String != "" {
		ChipSetData.Dimension = Dimension.String
	} else {
		ChipSetData.Dimension = "N/A"
	}
	if Weight.Valid && Weight.String != "" {
		ChipSetData.Weight = Weight.String
	} else {
		ChipSetData.Weight = "N/A"
	}
	if CabinetColor.Valid && CabinetColor.String != "" {
		ChipSetData.CabinetColor = CabinetColor.String
	} else {
		ChipSetData.CabinetColor = "N/A"
	}
	if CabinetLitre.Valid && CabinetLitre.String != "" {
		ChipSetData.CabinetLitre = CabinetLitre.String
	} else {
		ChipSetData.CabinetLitre = "N/A"
	}
	if Sensor.Valid && Sensor.String != "" {
		ChipSetData.Sensor = Sensor.String
	} else {
		ChipSetData.Sensor = "N/A"
	}
	if Connectivity.Valid && Connectivity.String != "" {
		ChipSetData.Connectivity = Connectivity.String
	} else {
		ChipSetData.Connectivity = "N/A"
	}
	if TouchPanel.Valid && TouchPanel.String != "" {
		ChipSetData.TouchPanel = TouchPanel.String
	} else {
		ChipSetData.TouchPanel = "N/A"
	}

	fmt.Printf("Test data %s\n", ChipSetData.Camera)

	// Query for the processor
	query = `SELECT ID, Procesors, Graphics, GraphicsSupport, MonitorSupport, Screen FROM Processors WHERE Procesors = ? AND IsDeleted = ?`
	err = db.QueryRow(query, Data.Processor, 0).Scan(
		&ProcessorData.ID,
		&ProcessorData.Processor,
		&ProcessorData.Graphics,
		&ProcessorData.GraphicsSupport,
		&ProcessorData.MonitorSupport,
		&ProcessorData.Screen,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result found for the processor, return a nil struct
			ProcessorData = ProcessorStruct{} // Return empty struct
		} else {
			return ChipSetData, ProcessorData, err // Return error for any other issue
		}
	}

	return ChipSetData, ProcessorData, nil
}

// func AddProcessors(Data ProcessorData) error {
// 	var ExistID int
// 	CheckQuery := `SELECT ID FROM Processors WHERE Procesors = ? `
// 	err := db.QueryRow(CheckQuery, Data.Processor).Scan(&ExistID)

// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			Query := `INSERT INTO Processors(Procesors, Graphics, GraphicsSupport, MonitorSupport, Screen, CreatedBy, CreatedAt) VALUES (?, ?, ?, ?, ?, Now())`
// 			_, err = db.Exec(Query, Data.Processor, Data.Graphics, Data.GraphicsSupport, Data.MonitorSupport, Data.Screen, Data.CreatedBy)
// 			if err != nil {
// 				return err
// 			}
// 			return nil
// 		}
// 	}
// 	return nil
// }

func AddNotification(DealID, UserID int, Message string) error {
	Query := `INSERT INTO DealNotification(DealID, UserID, Message) VALUES (?, ?, ?)`
	_, err := db.Exec(Query, DealID, UserID, Message)
	if err != nil {
		return err
	}
	return nil
}

func GetNotifications(UserID, DealID int) (DealsStruct, error) {
	var notifications DealsStruct
	query := `SELECT ID, Message, IsRead, CreatedAt FROM DealNotification WHERE UserID = ? AND DealID = ? ORDER BY CreatedAt DESC`
	err := db.QueryRow(query, UserID, DealID).Scan(&notifications.NotificationID, &notifications.NotificationMessage, &notifications.IsNotificationRead, &notifications.NotificationCreatedAt)
	if err == sql.ErrNoRows {
		// Return a default notification struct if no notification is found
		notifications.IsNotificationRead = 1
		return notifications, nil
	}
	if err != nil {
		return notifications, err
	}
	return notifications, nil
}

func UpdateNotification(NotificationID int) error {
	Query := `UPDATE DealNotification SET IsRead = ?, ReadAt = Now()  WHERE ID = ?`
	_, err := db.Exec(Query, 1, NotificationID)
	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}

func UpdateNotificationNew(DealID, UserID int, Message string) error {
	// insertQuery := `INSERT INTO Notifications (DealID, UserID, Message) VALUES (?, ?, ?)`
	// _, err := db.Exec(insertQuery, DealID, UserID, Message)
	// if err != nil {
	// 	fmt.Print("Error inserting new notification:", err)
	// 	return err
	// }
	// fmt.Printf("New notification created for Deal ID %d and User ID %d\n", DealID, UserID)

	return nil
}

func GetNotificationsNumber(UserID int) ([]Notifications, error) {
	var NotificationData []Notifications
	fmt.Print("UserID", UserID)
	query := `SELECT ID, DealID, UserID, Message, CreatedAt FROM DealNotification WHERE UserID = ? AND IsRead = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, UserID, 0)
	if err != nil {
		return NotificationData, err
	}
	defer rows.Close()
	for rows.Next() {
		var data Notifications
		err = rows.Scan(&data.ID, &data.DealID, &data.UserID, &data.Message, &data.CreatedAt)
		if err != nil {
			return NotificationData, err
		}
		NotificationData = append(NotificationData, data)
	}
	return NotificationData, nil
}
func GetGemNotificationNumber(UserID int) ([]GemNotificationStruct, error) {
	var NotificationData []GemNotificationStruct
	query := `SELECT ID, GemID, UserID, Message, CreatedAt FROM GemNotification WHERE UserID = ? AND IsRead = ? ORDER BY CreatedAt DESC`
	rows, err := db.Query(query, UserID, 0)
	if err != nil {
		return NotificationData, err
	}
	defer rows.Close()
	for rows.Next() {
		var data GemNotificationStruct
		err = rows.Scan(&data.ID, &data.GemID, &data.UserID, &data.Message, &data.CreatedAt)
		if err != nil {
			return NotificationData, err
		}
		NotificationData = append(NotificationData, data)
	}
	return NotificationData, nil
}

type Remainder struct {
	To      []string `json:"toMail"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	From    string   `json:"fromEmail"`
	CC      []string `json:"cc"`
	BCC     []string `json:"bcc"`
}

func CheckAndNotifyDueTasks() ([]Remainder, error) {
	fmt.Println("Got that")
	now := time.Now()
	nextHour := now.Add(1 * time.Hour)

	query := `
		SELECT t.ID, t.TaskDescription, t.TaskType, t.DueDate, t.IsCompany, t.CompanyID, u.Email, u.Name 
		FROM tasks t
		JOIN users u ON t.CreatedUser = u.ID
		WHERE t.DueDate >= ? AND t.DueDate < ?;
	`

	// Calculate start and end time of the next minute in the next hour
	startTime := nextHour.Format("2006-01-02 15:04:00")
	endTime := nextHour.Add(1 * time.Minute).Format("2006-01-02 15:04:00")

	// Execute the query
	rows, err := db.Query(query, startTime, endTime)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, fmt.Errorf("error executing query: %v", err)
	}
	defer rows.Close()

	var remainders []Remainder

	// Iterate over the tasks and process them
	for rows.Next() {
		var taskID, isCompany, companyID int
		var taskDescription, taskType, dueDate, ownerEmail, ownerName, customer string

		err := rows.Scan(&taskID, &taskDescription, &taskType, &dueDate, &isCompany, &companyID, &ownerEmail, &ownerName)
		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// Fetch customer name based on IsCompany flag
		if isCompany == 1 {
			customerQuery := `SELECT CompanyName FROM company WHERE ID = ?`
			err := db.QueryRow(customerQuery, companyID).Scan(&customer)
			if err != nil {
				log.Println("Error fetching company name:", err)
			}
		} else if isCompany == 0 {
			customerQuery := `SELECT Name FROM contacts WHERE ID = ?`
			err := db.QueryRow(customerQuery, companyID).Scan(&customer)
			if err != nil {
				log.Println("Error fetching contact name:", err)
			}
		}

		// Prepare the email subject and body
		Subject := fmt.Sprintf("Task Reminder: Task #%d (%s) Due on %s", taskID, taskType, dueDate)
		Body := fmt.Sprintf(`
			<html>
				<body>
					<p>Hello %s,</p>
					<p>This is a reminder that the following task is due:</p>
					<ul>
						<li><strong>Customer:</strong> %s</li>
						<li><strong>Description:</strong> %s</li>
						<li><strong>Due Date:</strong> %s</li>
					</ul>
					<p>Please take the necessary actions.</p>
					<p>Thanks & Regards,<br>Team CRM</p>
				</body>
			</html>
		`, ownerName, customer, taskDescription, dueDate)

		// Prepare the email struct similar to DealsMail
		remainder := Remainder{
			To:      []string{ownerEmail},
			Subject: Subject,
			Body:    Body,
			From:    ownerEmail,
			CC:      []string{"thedeveloper.ak@gmail.com"},
		}

		// Add remainder to the list
		remainders = append(remainders, remainder)

		// Push email to Redis queue
		err = DealsMail(remainder)
		if err != nil {
			log.Printf("Failed to queue email for Task #%d: %v", taskID, err)
		}
	}

	// Check if there are any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	fmt.Println(remainders)
	return remainders, nil
}

func DealsMail(deal Remainder) error {
	To := []string{deal.To[0]}
	Cc := []string{"thedeveloper.ak@gmail.com"}

	Subject := deal.Subject
	Body := deal.Body

	dealMail := DealsStruct{
		From:    deal.From,
		To:      To,
		CC:      Cc,
		Subject: Subject,
		Body:    Body,
	}

	mailData, err := json.Marshal(dealMail)
	if err != nil {
		return fmt.Errorf("failed to marshal mail data: %v", err)
	}

	err = redisClient.Publish("ticket_queue", string(mailData)).Err()
	if err != nil {
		return fmt.Errorf("failed to publish mail to Redis: %v", err)
	}

	return nil
}
func AddProcessor(data ProcessorStruct) (string, error) {
	var existId int
	query := `SELECT ID FROM Processors WHERE Procesors = ?`
	err := db.QueryRow(query, data.Processor).Scan(&existId)

	if err != nil {
		if err == sql.ErrNoRows {
			query := `INSERT INTO Processors (Procesors, Graphics, GraphicsSupport, MonitorSupport, Screen, CreatedBy, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, Now())`
			_, err := db.Exec(query, data.Processor, data.Graphics, data.GraphicsSupport, data.MonitorSupport, data.Screen, data.CreatedBy)
			if err != nil {
				return "", err
			}
			return "added", nil
		}
		return "", err
	}

	if existId != 0 {
		return "exist", nil
	}

	return "", nil
}

func GetProcessor() ([]ProcessorStruct, error) {
	var Processor []ProcessorStruct
	query := `SELECT ID, Procesors FROM Processors WHERE IsDeleted = ?`
	rows, err := db.Query(query, 0)
	if err != nil {
		return Processor, err
	}
	defer rows.Close()
	for rows.Next() {
		var data ProcessorStruct
		err = rows.Scan(&data.ID, &data.Processor)
		if err != nil {
			return Processor, err
		}
		Processor = append(Processor, data)
	}
	return Processor, nil
}

type EmailStruct struct {
	Recipient    string `json:"recipient"`
	Subject      string `json:"subject"`
	Body         string `json:"body"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Team         string `json:"team"`
	DealID       int    `json:"dealID"`
}

func UpdateMailStatus(ID int) error {
	query := `UPDATE Deals SET MailStatus = ? WHERE ID = ?`
	_, err := db.Exec(query, "Mail Sending", ID)
	if err != nil {
		return err
	}
	return nil
}

func AddGemNotification(UserID, GemID int, Message string) error {
	query := `INSERT INTO GemNotification (UserID, GemID, Message) VALUES (?, ?, ?)`
	_, err := db.Exec(query, UserID, GemID, Message)
	if err != nil {
		return err
	}
	return nil
}

func UpdateGemNotification(ID int) error {
	query := `UPDATE GemNotification SET IsRead = ?, ReadAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, 1, ID)
	if err != nil {
		return err
	}
	return nil
}

func GetProcessorData(Data ProcessorStruct) (ProcessorStruct, error) {
	var Res ProcessorStruct
	Query := `SELECT Procesors, Graphics, GraphicsSupport, MonitorSupport, Screen FROM Processors WHERE ID = ?`
	err := db.QueryRow(Query, Data.ID).Scan(&Res.Processor, &Res.Graphics, &Res.GraphicsSupport, &Res.MonitorSupport, &Res.Screen)
	if err != nil {
		return Res, err
	}
	return Res, nil
}

func RespectiveSalesPerson() ([]SalesPersons, error) {
	var salesPersons []SalesPersons
	query := `SELECT ID, Name FROM users WHERE IsDeleted = ? AND (Team = ? OR Team = ?)`
	rows, err := db.Query(query, 0, "Sales", "Marketing")
	if err != nil {
		return salesPersons, err
	}
	defer rows.Close()
	for rows.Next() {
		var data SalesPersons
		err = rows.Scan(&data.ID, &data.SalesPerson)
		if err != nil {
			return salesPersons, err
		}
		salesPersons = append(salesPersons, data)
	}
	return salesPersons, nil
}

func MarkAsRead(UserID int) error {
	fmt.Println("UserID:", UserID)

	query := `UPDATE GemNotification SET IsRead = ?, ReadAt = Now() WHERE UserID = ? AND IsRead = ?`

	// Execute the query
	result, err := db.Exec(query, 1, UserID, 0)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("No rows were updated. Possibly no unread notifications for this UserID.")
	} else {
		fmt.Printf("%d rows were updated.\n", rowsAffected)
	}

	DealQuery := `UPDATE DealNotification SET IsRead = ?, ReadAt = Now() WHERE UserID = ? AND IsRead = ?`

	_, err = db.Exec(DealQuery, 1, UserID, 0)
	if err != nil {
		return err
	}

	return nil
}
