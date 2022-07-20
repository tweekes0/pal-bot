package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

// Struct to present a record in the 'soundbites' table
type Soundbite struct {
	ID       int
	Name     string
	Username string
	UserID   string
	FilePath string
	FileHash string
	Created  time.Time
}

// Struct that holds the database connectivity
type SoundbiteModel struct {
	DB *sql.DB
}

// Initialize the 'soundbites' table in the sqlite db
func (m *SoundbiteModel) Initialize() error {
	stmt := `CREATE TABLE IF NOT EXISTS soundbites (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		user_id	TEXT NOT NULL,
		filepath TEXT NOT NULL,
		filehash TEXT NOT NULL,
		created TEXT NOT NULL,
		UNIQUE(name)
	);`

	if _, err := m.DB.Exec(stmt); err != nil {
		return err
	}

	return nil
}

// Insert Soundbites metadata into the 'soundbites' table
func (m *SoundbiteModel) Insert(name, username, uid, filepath, filehash string) (int, error) {
	stmt := `INSERT INTO soundbites (name, username, user_id, filepath, filehash, created)  
	VALUES(?, ?, ?, ?, ?, datetime('now'));`

	res, err := m.DB.Exec(stmt, name, username, uid, filepath, filehash)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint") {
			return 0, ErrUniqueConstraint
		}

		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// Gets a Soundbite based on the name command
func (m *SoundbiteModel) Get(name string) (*Soundbite, error) {
	stmt := `SELECT * FROM soundbites WHERE name = ?;`

	var date string
	s := &Soundbite{}

	err := m.DB.QueryRow(stmt, name).Scan(&s.ID, &s.Name, &s.Username, &s.UserID, &s.FilePath, &s.FileHash, &date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrDoesNotExist
		}

		return nil, err
	}

	t, err := time.Parse(TIME_LAYOUT, date)
	if err != nil {
		return nil, err
	}

	s.Created = t
	return s, nil
}

// Get all the soundbites in the 'soundbites' table
func (m *SoundbiteModel) GetAll() ([]*Soundbite, error) {
	stmt := `SELECT * FROM soundbites;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	soundbites := []*Soundbite{}

	for rows.Next() {
		var date string
		s := &Soundbite{}

		err = rows.Scan(&s.ID, &s.Name, &s.Username, &s.UserID, &s.FilePath, &s.FileHash, &date)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse(TIME_LAYOUT, date)
		if err != nil {
			return nil, err
		}

		s.Created = t
		soundbites = append(soundbites, s)
	}

	if len(soundbites) == 0 {
		return soundbites, ErrNoRecords
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return soundbites, nil
}

// Check whether a soundbite exists based on the name of the command and it's filehash
func (m *SoundbiteModel) Exists(name, hash string) (bool, error) {
	var exists bool

	stmt := `SELECT EXISTS(SELECT 1 FROM soundbites WHERE name = ? OR filehash = ?);`
	err := m.DB.QueryRow(stmt, name, hash).Scan(&exists)

	return exists, err
}

// Deletes the soundbite if the user_id and name belong to the same record
func (m *SoundbiteModel) Delete(name, uid string) error {
	if exists, _ := m.Exists(name, ""); !exists {
		return ErrDoesNotExist
	}

	if err := m.userCreatedSound(name, uid); err != nil {
		return err
	}

	stmt := `DELETE FROM soundbites WHERE name = ? AND user_id = ?;`

	res, err := m.DB.Exec(stmt, name, uid)
	if err != nil {
		return err
	}

	c, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if int(c) == 0 {
		return ErrDoesNotExist
	}

	return nil
}

func (m *SoundbiteModel) UpdateName(oldName, newName string) error {
	exists, err := m.Exists(newName, "")
	if err != nil {
		return err
	}

	if exists {
		return ErrUniqueConstraint
	}

	stmt := `UPDATE soundbites SET name = ? where name = ?`

	_, err = m.DB.Exec(stmt, newName, oldName)
	if err != nil {
		return err
	}

	return nil
}

// Checks that the user_id of the soundbite belongs to the user requesting the delete
func (m *SoundbiteModel) userCreatedSound(name, uid string) error {
	var exists bool
	stmt := `SELECT EXISTS(SELECT 1 FROM soundbites where name = ? AND user_id = ?);`

	err := m.DB.QueryRow(stmt, name, uid).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrCommandOwnership
	}

	return nil
}
