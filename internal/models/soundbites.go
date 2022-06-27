package models

import (
	"database/sql"
	"errors"
	"time"
)

type Soundbite struct {
	ID       int
	Name     string
	Username string
	UserID   string
	FilePath string
	FileHash string
	Created  time.Time
}

type SoundbiteModel struct {
	DB *sql.DB
}

// Insert Soundbite's metadata into the 'soundbites' table
func (m *SoundbiteModel) Insert(name, username, uid, filepath string) (int, error) {
	stmt := `INSERT INTO soundbites (name, username, user_id, filepath, filehash, created) as 
	values(?, ?, ?, ?, ?, UTC_TIMESTAMP())`

	res, err := m.DB.Exec(stmt, name, username, uid, filepath)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Gets Soundbite based on the name command
func (m *SoundbiteModel) Get(name string) (*Soundbite, error) {
	stmt := `SELECT id, name, username, user_id, filepath, filehash, created FROM soundbites
	WHERE name = ?`

	row := m.DB.QueryRow(stmt, name)

	s := &Soundbite{}
	err := row.Scan(&s.ID, &s.Name, &s.Username, &s.UserID, &s.FilePath, &s.FileHash, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecords
		}

		return nil, err
	}

	return s, nil
}

// Check whether a soundbite exists based on the name of the command and the filehash
func (m *SoundbiteModel) Exists(name, hash string) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT true FROM soundbites where name = ? OR filehash = ?)`
	err := m.DB.QueryRow(stmt, name, hash).Scan(&exists)

	return exists, err
}

func (m *SoundbiteModel) Delete(name, uid string) error {
	if err := m.userCreatedCommand(name, uid); err != nil {
		return err
	}

	stmt := `DELETE FROM soundbites WHERE name = ? AND user_id = ? LIMIT 1`

	res, err := m.DB.Exec(stmt, name, uid)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return nil
	}

	if int(c) == 0 {
		return ErrDoesNotExist
	}

	return nil
}

// Checks that the user_id of the soundbite belongs to the user requesting the delete
func (m *SoundbiteModel) userCreatedCommand(name, uid string) error {
	var exists bool
	stmt := `SELECT EXISTS(SELECT true FROM soundbites where name = ? AND user_id = ?)`

	err := m.DB.QueryRow(stmt, name, uid).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrCommandOwnership
	}

	return nil
}