package store

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/go-sqlite"
	"github.com/jmoiron/sqlx"
)

// register driver
var _ sqlite.Driver

// temp folder store the db file
var basedir = filepath.Join(os.TempDir(), "dot")

// temp db file path
var datafile = filepath.Join(basedir, "data.db")

// db stored
var db *sqlx.DB

func Init() {
	err := os.MkdirAll(basedir, os.ModeDir)
	if err != nil {
		log.Fatal(err)
	}
	db1, err := sqlx.Open("sqlite", datafile)
	if err != nil {
		log.Fatal(err)
	}
	err = initDB(db1)
	if err != nil {
		log.Fatal(err)
	}
	db = db1
}

func initDB(db1 *sqlx.DB) error {
	rows, err := db1.Query("SELECT name FROM sqlite_master WHERE type='table' and name in ('components', 'providers', 'provider_requirements')")
	if err != nil {
		return err
	}
	defer rows.Close()
	mp := map[string]string{
		"components":            TableComponents,
		"providers":             TableProviders,
		"provider_requirements": TableProviderRequirements,
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		mp[name] = ""
	}
	for _, ddl := range mp {
		if ddl == "" {
			continue
		}
		_, err := db1.Exec(ddl)
		if err != nil {
			return err
		}
	}
	return nil
}

func Clean() {
	stmt := []string{"drop table components", "drop table providers", "drop table provider_requirements"}
	for i := range stmt {
		_, err := db.Exec(stmt[i])
		if err != nil {
			log.Fatal(err)
		}
	}
}
