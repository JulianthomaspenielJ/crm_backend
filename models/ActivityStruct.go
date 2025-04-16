package models

import (
	"database/sql"
)

type MailingStruct struct {
	ID          int    `json:"id"`
	MailFrom    string `json:"mailFrom"`
	MailTo      string `json:"mailTo"`
	MailCC      sql.NullString
	MailSubject string `json:"mailSubject"`
	MailContent string `json:"mailBody"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	IsCompany   int    `json:"isCompany"`
	CompanyID   int    `json:"companyID"`
	CreatedUser int    `json:"UserId"`
	Status      string `json:"status"`
}

type NotesStruct struct {
	ID              int    `json:"id"`
	NoteDescription string `json:"notes"`
	CreatedBy       string `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
	UpdatedBy       sql.NullString
	UpdatedAt       sql.NullString
	IsCompany       int `json:"isCompany"`
	CompanyID       int `json:"companyID"`
	CreatedUser     int `json:"UserId"`
}

type CallsStruct struct {
	ID                int    `json:"id"`
	CallFrom          string `json:"callFrom"`
	CallTo            int    `json:"callTo"`
	CallAt            string `json:"callAt"`
	CallDescription   string `json:"notes"`
	CallOutCome       string `json:"callOutcome"`
	CallDirection     string `json:"callDirection"`
	CreatedBy         string `json:"createdBy"`
	CreatedAt         string `json:"createdAt"`
	UpdatedBy         sql.NullString
	UpdatedAt         sql.NullString
	IsCompany         int    `json:"isCompany"`
	CompanyID         int    `json:"companyID"`
	CreatedUser       int    `json:"UserId"`
	CallToPersonName  string `json:"CallToPersonName"`
	CallToPersonEmail string `json:"CallToPersonEmail"`
	CallToPersonID    int    `json:"CallToPersonID"`
}

type TaskStruct struct {
	ID              int    `json:"id"`
	TaskDescription string `json:"notes"`
	TaskType        string `json:"task"`
	DueDate         string `json:"date"`
	IsCompany       int    `json:"isCompany"`
	CompanyID       int    `json:"companyID"`
	CreatedUser     int    `json:"UserId"`
	CreatedBy       string `json:"createdBy"`
	CreatedAt       string `json:"createdAt"`
	UpdatedBy       sql.NullString
	UpdatedAt       sql.NullString
	Status          int `json:"status"`
}

type MeetingStruct struct {
	ID int `json:"id"`
	//MeetingScheduledBy string `json:"createdBy"`
	Attendees          string `json:"Attendees"`
	OutCome            string `json:"outCome"`
	MeetingAt          string `json:"date"`
	MeetingDescription string `json:"notes"`
	Duration           string `json:"duration"`
	Status             int    `json:"status"`
	IsCompany          int    `json:"isCompany"`
	CompanyID          int    `json:"companyID"`
	CreatedBy          string `json:"createdBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedBy          sql.NullString
	UpdatedAt          sql.NullString
	CreatedUser        int    `json:"UserId"`
	AttendeesList      string `json:"attendeesList"`
}

type Attendees struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
