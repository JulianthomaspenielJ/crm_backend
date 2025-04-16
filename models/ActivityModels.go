package models

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func StoreMailContent(From, To, Subject, Body, CC, CreatedBy string, IsCompany, CompanyID, CreatedUser int) error {
	query := `INSERT INTO mail(MailFrom, MailTo, MailCC, MailSubject, MailContent, CreatedBy, CreatedAt, IsCompany, CompanyID, CreatedUser)
				VALUES(?, ?, ?, ?, ?, ?, Now(), ?, ?, ?)`
	res, err := db.Exec(query, From, To, CC, Subject, Body, CreatedBy, IsCompany, CompanyID, CreatedUser)
	if err != nil {
		return fmt.Errorf("error occurred while query executing %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error occurred while getting last inserting id %s", err)
	}

	err = AuditFunction("Update", CreatedBy, ID, "127.0.0.1", "Email data inserted successfully")
	if err != nil {
		return fmt.Errorf("error occurred while audit table executing %s ", err)
	}

	return nil
}

func FetchMailData(ID, IsCompany int) ([]MailingStruct, bool, error) {
	var mails []MailingStruct
	var status sql.NullString
	query := `SELECT MailFrom, MailTo, MailCC, MailSubject, MailContent, CreatedBy, CreatedAt, Status FROM mail WHERE CompanyID = ? AND IsCompany = ? ORDER BY CreatedAt DESC`

	res, err := db.Query(query, ID, IsCompany)
	if err != nil {
		return mails, false, fmt.Errorf("error occurred while scanning data: %v", err)
	}
	for res.Next() {
		var mail MailingStruct
		err := res.Scan(&mail.MailFrom, &mail.MailTo, &mail.MailCC, &mail.MailSubject, &mail.MailContent, &mail.CreatedBy, &mail.CreatedAt, &status)
		if err != nil {
			if err == sql.ErrNoRows {
				return mails, false, nil
			}
			return mails, false, fmt.Errorf("error scanning data %s ", err)
		}
		if status.Valid {
			mail.Status = status.String
		} else {
			mail.Status = "N/A"
		}
		mails = append(mails, mail)
	}
	return mails, true, nil
}

func FetchNotesData(ID, IsCompany int) ([]NotesStruct, bool, error) {
	var note []NotesStruct
	query := `SELECT ID, NoteDescription, CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, CreatedUser FROM notes WHERE CompanyID = ? AND IsCompany = ? ORDER BY CreatedAt DESC`

	res, err := db.Query(query, ID, IsCompany)
	if err != nil {
		return note, false, fmt.Errorf("error occurred while scanning data: %v", err)
	}

	for res.Next() {
		var notes NotesStruct
		err := res.Scan(&notes.ID, &notes.NoteDescription, &notes.CreatedBy, &notes.CreatedAt, &notes.UpdatedBy, &notes.UpdatedAt, &notes.CreatedUser)
		if err != nil {
			if err == sql.ErrNoRows {
				return note, false, nil
			}
			return note, false, fmt.Errorf("error scanning data %s ", err)
		}
		note = append(note, notes)
	}
	return note, true, nil
}
func FetchCallsData(ID, IsCompany int) ([]CallsStruct, bool, error) {
	var calls []CallsStruct
	query := `SELECT ID, CallFrom, CallTo, CallAt, CallDescription, CallOutCome, CallDirection, CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, IsCompany, CompanyID, CreatedUser FROM call_logs WHERE CompanyID = ? AND IsCompany = ? ORDER BY CreatedAt DESC`

	resCall, err := db.Query(query, ID, IsCompany)
	if err != nil {
		return calls, false, fmt.Errorf("error occurred while executing query: %v", err)
	}
	defer resCall.Close()

	for resCall.Next() {
		var call CallsStruct
		err := resCall.Scan(&call.ID, &call.CallFrom, &call.CallTo, &call.CallAt, &call.CallDescription, &call.CallOutCome, &call.CallDirection, &call.CreatedBy, &call.CreatedAt, &call.UpdatedBy, &call.UpdatedAt, &call.IsCompany, &call.CompanyID, &call.CreatedUser)
		if err != nil {
			return calls, false, fmt.Errorf("error scanning data: %v", err)
		}

		var contactName, contactEmail string
		var contactID int
		contactQuery := `SELECT Name, Email FROM contacts WHERE ID = ?`
		contactRes := db.QueryRow(contactQuery, call.CallTo)

		err = contactRes.Scan(&contactName, &contactEmail)
		if err != nil {
			if err == sql.ErrNoRows {
				contactName = "Unknown"
				contactID = 0
				contactEmail = "Unknown"
			} else {
				return calls, false, fmt.Errorf("error scanning contact data: %v", err)
			}
		}

		call.CallToPersonName = contactName
		call.CallToPersonEmail = contactEmail
		call.CallToPersonID = contactID

		calls = append(calls, call)
	}

	if len(calls) == 0 {
		return calls, false, nil
	}
	return calls, true, nil
}

func FetchTaskData(ID, IsCompany int) ([]TaskStruct, bool, error) {
	var tasks []TaskStruct
	query := `SELECT ID, TaskDescription,TaskType, DueDate, CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, CreatedUser FROM tasks WHERE CompanyID = ? AND IsCompany = ? ORDER BY CreatedAt DESC`

	resCall, err := db.Query(query, ID, IsCompany)
	if err != nil {
		return tasks, false, fmt.Errorf("error occurred while executing query: %v", err)
	}
	defer resCall.Close()

	for resCall.Next() {
		var task TaskStruct
		err := resCall.Scan(&task.ID, &task.TaskDescription, &task.TaskType, &task.DueDate, &task.CreatedBy, &task.CreatedAt, &task.UpdatedBy, &task.UpdatedAt, &task.CreatedUser)
		if err != nil {
			return tasks, false, fmt.Errorf("error scanning data: %v", err)
		}
		tasks = append(tasks, task)
	}

	if len(tasks) == 0 {
		return tasks, false, nil
	}
	return tasks, true, nil
}
func FetchMeetingData(ID, IsCompany int) ([]MeetingStruct, bool, error) {
	var meetings []MeetingStruct
	query := `SELECT ID, MeetingScheduledBy, Attendees, OutCome, MeetingAt, MeetingDescription, CreatedBy, CreatedAt, UpdatedBy, UpdatedAt, Duration, CreatedUser FROM meeting_logs WHERE CompanyID = ? AND IsCompany = ? ORDER BY CreatedAt DESC`

	resMeetings, err := db.Query(query, ID, IsCompany)
	if err != nil {
		return meetings, false, fmt.Errorf("error occurred while executing query: %v", err)
	}
	defer resMeetings.Close()

	for resMeetings.Next() {
		var meet MeetingStruct
		var attendeesStr string
		err := resMeetings.Scan(&meet.ID, &meet.CreatedBy, &attendeesStr, &meet.OutCome, &meet.MeetingAt, &meet.MeetingDescription, &meet.CreatedBy, &meet.CreatedAt, &meet.UpdatedBy, &meet.UpdatedAt, &meet.Duration, &meet.CreatedUser)
		if err != nil {
			return meetings, false, fmt.Errorf("error scanning data: %v", err)
		}

		// Process attendees
		attendeesIDs := strings.Split(attendeesStr, ",")
		var attendeesInfo []string

		for _, idStr := range attendeesIDs {
			idStr = strings.TrimSpace(idStr)

			if idStr == "" {
				continue
			}

			attendeeID, err := strconv.Atoi(idStr)
			if err != nil {
				log.Printf("Invalid attendee ID %s: %v", idStr, err)
				continue
			}

			var attendeeName, attendeeEmail string
			attendeeQuery := `SELECT Name, Email FROM contacts WHERE ID = ?`
			err = db.QueryRow(attendeeQuery, attendeeID).Scan(&attendeeName, &attendeeEmail)
			if err != nil {
				if err == sql.ErrNoRows {
					attendeeName = "Unknown"
					attendeeEmail = "Unknown"
				} else {
					return meetings, false, fmt.Errorf("error fetching attendee data: %v", err)
				}
			}
			attendeesInfo = append(attendeesInfo, fmt.Sprintf("%s(%s)", attendeeName, attendeeEmail))
		}
		meet.AttendeesList = strings.Join(attendeesInfo, ", ")

		meetings = append(meetings, meet)
	}

	if len(meetings) == 0 {
		return meetings, false, nil
	}
	return meetings, true, nil
}

func InsertCall(CallFrom, CallDescription, CallOutCome, CallAt, CallDirection, CreatedBy string, CallTo, CreatedUser, IsCompany, CompanyID int) error {
	query := `INSERT INTO call_logs(CallFrom, CallTo, CallAt, CallDescription, CallOutCome, CallDirection, CreatedBy, CreatedAt, IsCompany, CompanyID, CreatedUser) 
				VALUES (?, ?, ?, ?, ?, ?, ?, Now(), ?, ?, ?)`
	res, err := db.Exec(query, CallFrom, CallTo, CallAt, CallDescription, CallOutCome, CallDirection, CreatedBy, IsCompany, CompanyID, CreatedUser)
	if err != nil {
		return fmt.Errorf("error occurred while call log inserting query %s", err)
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = UpdateAudit("Create", CreatedBy, int(ID), "127.0.0.1", "New Call log inserted")
	if err != nil {
		return err
	}

	return nil
}

func InsertTask(Description, TaskType, DueDate, CreatedBy string, CreatedUser, IsCompany, CompanyID int) error {
	query := `INSERT INTO tasks(TaskDescription,TaskType, DueDate, IsCompany, CompanyID, CreatedUser, CreatedBy, CreatedAt, Status)
			VALUES(?, ?, ?, ?, ?, ?,?, Now(), ?)`
	res, err := db.Exec(query, Description, TaskType, DueDate, IsCompany, CompanyID, CreatedUser, CreatedBy, 0)
	if err != nil {
		return fmt.Errorf("error occurred while call log inserting query %s", err)
	}

	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = UpdateAudit("Create", CreatedBy, int(ID), "127.0.0.1", "New Task inserted")

	if err != nil {
		return err
	}

	return nil
}

func InsertNotes(Description, CreatedBy string, IsCompany, CompanyID, CreatedUser int) error {
	query := `INSERT INTO notes(NoteDescription,CreatedBy, CreatedAt, IsCompany, CompanyID, CreatedUser) 
				VALUES (?, ?, Now(), ?, ?, ?)`
	res, err := db.Exec(query, Description, CreatedBy, IsCompany, CompanyID, CreatedUser)
	if err != nil {
		return fmt.Errorf("error occurred while call log inserting query %s", err)
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = UpdateAudit("Create", CreatedBy, int(ID), "127.0.0.1", "New Call log inserted")
	if err != nil {
		return err
	}

	return nil
}

func InsertMeeting(Duration, ScheduledBy, Attendees, OutCome, MeetingAt, MeetingDescription, CreatedBy string, CreatedUser, IsCompany, CompanyID int) error {
	query := `INSERT INTO meeting_logs(MeetingScheduledBy, Attendees, OutCome, MeetingAt, MeetingDescription, CreatedBy, IsCompany, CompanyID, CreatedAt, Status, CreatedUser, Duration) 
				VALUES(?, ?, ?, ?, ?, ?, ?, ?, Now(), ?, ?, ?)`
	res, err := db.Exec(query, ScheduledBy, Attendees, OutCome, MeetingAt, MeetingDescription, CreatedBy, IsCompany, CompanyID, 0, CreatedUser, Duration)
	if err != nil {
		return err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return nil
	}
	err = UpdateAudit("Create", CreatedBy, int(ID), "127.0.0.1", "New meeting log inserted")
	if err != nil {
		return err
	}
	return nil
}

func UpdateNotes(NoteDescription, UpdatedBy string, ID int) error {
	query := `UPDATE notes SET NoteDescription = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, NoteDescription, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0.1", "Notes data updated")
	if err != nil {
		return err
	}
	return nil
}

func UpdateCall(CallDescription, CallOutcome, CallsAt, CallDirection, UpdatedBy string, ID int) error {
	query := `UPDATE call_logs SET CallDescription = ?, CallOutCome = ?, CallAt = ?, CallDirection = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, CallDescription, CallOutcome, CallsAt, CallDirection, UpdatedBy, ID)
	if err != nil {
		return err
	}

	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0.1", "Calls data updated")
	if err != nil {
		return err
	}
	return nil
}

func UpdateTask(TaskDescription, TaskType, DueDate, UpdatedBy string, Status, ID int) error {
	query := `UPDATE call_logs SET TaskDescription = ?, TaskType = ?, DueDate = ?, Status = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Exec(query, TaskDescription, TaskType, DueDate, Status, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = UpdateAudit("Update", UpdatedBy, ID, "127.0.0.1", "Task data updated")
	if err != nil {
		return err
	}
	return nil
}

func GetAttendees() ([]Attendees, error) {
	var attendees []Attendees

	query := `SELECT ID, Name, Email FROM users WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	row, err := db.Query(query, 0)
	if err != nil {
		return nil, fmt.Errorf("database error: %s", err)
	}
	defer row.Close()

	for row.Next() {
		var items Attendees
		err := row.Scan(&items.ID, &items.Name, &items.Email)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		attendees = append(attendees, items)
	}
	return attendees, nil
}
