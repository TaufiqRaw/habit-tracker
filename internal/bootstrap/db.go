package bootstrap

import "fmt"

import (
	"database/sql"
	"habit-tracker/internal/domain"
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
    // TODO: create versioning (like migration), if table exist update the table depend on the version
    {
        habitSql := 
            "CREATE TABLE IF NOT EXISTS " + domain.HabitTableName + " (\n" +
                domain.HabitCols.Str.Id + " INTEGER PRIMARY KEY AUTOINCREMENT, \n" + 
                domain.HabitCols.Str.LastHabitID + " INTEGER, \n" +
                domain.HabitCols.Str.Name + " VARCHAR(50) NOT NULL, \n" + 
                domain.HabitCols.Str.Amount + " INTEGER NOT NULL, \n" +
                domain.HabitCols.Str.RestDay + " INTEGER DEFAULT 0, \n" + 
                domain.HabitCols.Str.RestDayMode + " VARCHAR(50) NOT NULL, \n" + 
                domain.HabitCols.Str.StartAt + " CHARACTER(10) NOT NULL, \n" + 
                domain.HabitCols.Str.Unit + " VARCHAR(50) NOT NULL, \n" +
                domain.HabitCols.Str.ArchivedAt + " CHARACTER(10), \n" + 
                "FOREIGN KEY("+ domain.HabitCols.Str.LastHabitID +") REFERENCES "+ 
                    domain.HabitTableName +"("+ domain.HabitCols.Str.Id +")" +
            ");";
        
        trackerSql :=
            "CREATE TABLE IF NOT EXISTS " + domain.TrackerTableName + " (\n" +
                domain.TrackerCols.Str.HabitId + " INTEGER, \n" +
                domain.TrackerCols.Str.At + " CHARACTER(10), \n" +
                domain.TrackerCols.Str.Amount + " INTEGER, \n" +
                "PRIMARY KEY ("+ domain.TrackerCols.Str.HabitId + ", " + 
                    domain.TrackerCols.Str.At + "), \n" +
                "FOREIGN KEY("+ domain.TrackerCols.Str.HabitId +") REFERENCES "+ 
                    domain.HabitTableName +"("+ domain.HabitCols.Str.Id +
                    ") ON DELETE CASCADE ON UPDATE CASCADE" +
            ");";

        sqls := []string{
            habitSql, trackerSql,
        }

        for i, sql := range sqls {
            _, err := db.Exec(sql)
            if err != nil {
                log.Fatal("initDB:create table::sql index: " + fmt.Sprint(i) + "\n", err)
            }
        }
    }

	return db
}