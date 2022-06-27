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