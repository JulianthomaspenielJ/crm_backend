package models

import (
	"time"
)

func GetUsersData() ([]UserDetails, error) {
	var users []UserDetails

	query := `SELECT ID, Name, Email, Mobile, Role, Designation, Team, CreatedBy, CreatedAt FROM users WHERE IsDeleted = ? AND IsApproved = ? ORDER BY CreatedAt DESC`
	res, err := db.Query(query, 0, 1)
	if err != nil {
		return users, err
	}
	defer res.Close()

	for res.Next() {
		var CreatedDateRaw []uint8
		var user UserDetails

		if err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Mobile, &user.Role, &user.Designation, &user.Team, &user.CreatedBy, &CreatedDateRaw); err != nil {
			return users, err
		}
		CreatedStr := string(CreatedDateRaw)
		CreatedDate, err := time.Parse("2006-01-02 15:04:05", CreatedStr)
		if err != nil {
			return users, err
		}
		user.CreatedAt = CreatedDate.Format("Jan 02 2006")
		users = append(users, user)
	}
	return users, nil
}

func GetDeletedUsersData() ([]UserDetails, error) {
	var users []UserDetails

	query := `SELECT ID, Name, Email, Mobile, Role, Designation, Team, UpdatedBy, UpdatedAt FROM users WHERE IsDeleted = ? AND IsApproved = ? ORDER BY CreatedAt DESC`

	res, err := db.Query(query, 1, 0)
	if err != nil {
		return users, err
	}
	defer res.Close()

	for res.Next() {
		var user UserDetails
		var updatedAtRaw []uint8

		if err := res.Scan(&user.ID, &user.Name, &user.Email, &user.Mobile, &user.Role, &user.Designation, &user.Team, &user.UpdatedBy, &updatedAtRaw); err != nil {
			return users, err
		}
		updatedAtStr := string(updatedAtRaw)

		updatedAt, err := time.Parse("2006-01-02 15:04:05", updatedAtStr)
		if err != nil {
			return users, err
		}

		user.CreatedAt = updatedAt.Format("Jan 02 2006")

		users = append(users, user)
	}

	return users, nil
}

func UpdateUserManagementData(ID int, Name, Email, Role, Designation, Team, Mobile, UpdatedBy string) error {
	query := `UPDATE users SET Name = ?, Email = ?, Mobile = ?, Role = ?, Designation = ?, Team = ?, UpdatedBy = ? WHERE ID = ?`

	_, err := db.Query(query, Name, Email, Mobile, Role, Designation, Team, UpdatedBy, ID)
	if err != nil {
		return err
	}

	err = AuditFunction("Update", UpdatedBy, int64(ID), "127.0.0.1", "User Data updated")
	if err != nil {
		return err
	}
	return nil
}

func RestoreUserManagementData(ID int, UpdatedBy string) error {
	query := `UPDATE users SET IsDeleted = ?, IsApproved = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Query(query, 0, 1, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = AuditFunction("Recover", UpdatedBy, int64(ID), "127.0.0.1", "User Data restored")
	if err != nil {
		return err
	}
	return nil
}

func RemoveUser(ID int, UpdatedBy string) error {
	query := `UPDATE users SET IsDeleted = ?, IsApproved = ?, UpdatedBy = ?, UpdatedAt = Now() WHERE ID = ?`
	_, err := db.Query(query, 1, 0, UpdatedBy, ID)
	if err != nil {
		return err
	}
	err = AuditFunction("Delete", UpdatedBy, int64(ID), "127.0.0.1", "User Data removed")
	if err != nil {
		return err
	}
	return nil
}

func UserNameList() ([]UserDetails, error) {
	Team := "Pricing"
	query := `SELECT ID, Name, Email FROM users WHERE IsDeleted = ? AND IsApproved = ? AND Team = ? ORDER BY CreatedAt DESC`
	res, err := db.Query(query, 0, 1, Team)
	if err != nil {
		return nil, err
	}

	defer res.Close()
	var Users []UserDetails
	for res.Next() {
		var user UserDetails
		err = res.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		Users = append(Users, user)
	}
	return Users, nil
}

func UserNameListDashboards() ([]UserDetails, error) {
	query := `SELECT ID, Name, Email FROM users WHERE IsDeleted = ? AND IsApproved = ? ORDER BY CreatedAt DESC`
	res, err := db.Query(query, 0, 1)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var Users []UserDetails
	for res.Next() {
		var user UserDetails
		err = res.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		Users = append(Users, user)
	}
	return Users, nil
}

func StoreFCMToken(UserID int, FCMToken string) error {
	query := `UPDATE users SET FCMToken = ? WHERE ID = ?`
	_, err := db.Exec(query, FCMToken, UserID)
	if err != nil {
		return err
	}
	return nil
}

func GetTeam() ([]TeamStructure, error) {
	var teams []TeamStructure
	query := `SELECT ID, TeamName FROM Team WHERE IsDeleted = ? ORDER BY CreatedAt DESC`
	res, err := db.Query(query, 0)
	if err != nil {
		return teams, err
	}
	for res.Next() {
		var team TeamStructure
		err = res.Scan(&team.ID, &team.TeamName)
		if err != nil {
			return teams, err
		}
		teams = append(teams, team)
	}
	defer res.Close()
	return teams, nil
}
