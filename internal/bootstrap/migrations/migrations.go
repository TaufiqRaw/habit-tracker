package migrations

import (
	"database/sql"
)

type migration struct{
	Name string
	Migration string
}

func Run(db *sql.DB) error {
	//Update this when creating new migration.
	//Order matter!
	migrations := []migration{
		init_table(),
	}

	if !helper.isMigrationTblExist(db) {
		err := helper.createMigrationTbl(db)
		if err != nil {
			return err
		}
	}

	executedMigration, err := helper.getExecutedMigration(db)
	if err != nil {
		return err
	}
	
	for _, m := range migrations {
		if executedMigration[m.Name] {
			continue
		} 
		err := helper.migrate(db, m)
		if err != nil {
			return err
		}
	}
	return nil	
}
