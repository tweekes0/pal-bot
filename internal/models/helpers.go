package models

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tweekes0/pal-bot/config"
)

func modelsTestSetup(t *testing.T) (SoundbiteModel, func()) {
	t.Parallel()

	f, err := ioutil.TempFile("", "*")
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}	

	db, err := sql.Open(config.DB_DRIVER, f.Name())
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("failed: %v", err)
	}

	m := SoundbiteModel{DB: db}
	m.Initialize()
	teardown := func() {
		os.Remove(f.Name())
	}

	return m, teardown
}