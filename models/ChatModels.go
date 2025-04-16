package models

import "time"

type ChatStruct struct {
	ID        int    `json:"id"`
	DealID    int    `json:"dealID"`
	Message   string `json:"message"`
	IsMedia   int    `json:"isMedia"`
	Media     string `json:"media"`
	CreatedBy string `json:"createdBy"`
	CreatedAt string `json:"createdAt"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedAt string `json:"updatedAt"`
	IsDeleted int    `json:"isDeleted"`
}

func InsertChat(DealID, IsMedia int, Message, Media, CreatedBY string) error {
	query := `INSERT INTO DealChat(DealID, Message, IsMedia, Media, CreatedBy, CreatedAt, IsDeleted) 
			VALUES(?, ?, ?, ?, ?, Now(), ?)`
	res, err := db.Exec(query, DealID, Message, IsMedia, Media, CreatedBY, 0)
	if err != nil {
		return err
	}
	ID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	err = UpdateAudit("Create", CreatedBY, int(ID), "127.0.0.1", "Deal_Chat created")
	if err != nil {
		return err
	}
	return nil
}

func GetChat(ID int) ([]ChatStruct, error) {
	query := `SELECT ID, DealID, Message, IsMedia, Media, CreatedBy, CreatedAt FROM DealChat WHERE DealID = ? AND IsDeleted = ?`
	rows, err := db.Query(query, ID, 0)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var chats []ChatStruct
	for rows.Next() {
		var createdAt string
		var chat ChatStruct
		err := rows.Scan(&chat.ID, &chat.DealID, &chat.Message, &chat.IsMedia, &chat.Media, &chat.CreatedBy, &createdAt)
		if err != nil {
			return nil, err
		}
		if len(createdAt) > 0 {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", string(createdAt))
			if err != nil {
				return nil, err
			}
			chat.CreatedAt = parsedTime.Format("02 January 2006 15:04")
		} else {
			chat.CreatedAt = ""
		}
		chats = append(chats, chat)
	}
	return chats, nil
}
