package bootstrap

import (
	"database/sql"
	"habit-tracker/internal/bootstrap/migrations"
	"log"
	"os"
	"path/filepath"
)

func InitDB() *sql.DB {
	var db *sql.DB
	{
		//TODO: get better exePath
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}

		pathToApp := filepath.Join(filepath.Dir(exePath), "data.db")
		//TODO: cross platform file (right now only windows)
		db, err = sql.Open("sqlite", "file:///"+pathToApp)
		if err != nil {
			log.Fatal(err)
		}

		pingErr := db.Ping()
		if pingErr != nil {
			log.Fatal(pingErr)
		}
	}

    // enable foreign key support
    {
        _, err := db.Exec("PRAGMA foreign_keys = ON")
        if err != nil {
			log.Fatal(err)
		}
    }

    // create table
	err := migrations.Run(db)
	if err != nil {
		log.Fatalf("migration : %v",err)
	}

	return db
}