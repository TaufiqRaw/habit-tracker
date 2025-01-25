package migrations

import (
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type h struct{}

var helper = h{} //DONT MUTATE THIS VARIABLE

func(h) createMigrationTbl(db *sql.DB) error {
	sql := "CREATE TABLE IF NOT EXISTS " + "migration" + " (\n" +
		"name" + " TEXT PRIMARY KEY \n" +
		");"

	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("createMigrationTbl : %v", err)
	}
	return nil

}

func(h) isMigrationTblExist(db *sql.DB) bool {
	sql, args := sq.Select("name").
		From("sqlite_master").
		Where(sq.And{
			sq.Eq{
				"type": "table",
			},
			sq.Eq{
				"name": "{migration}",
			},
		}).MustSql()

	row := db.QueryRow(sql, args...)

	var name string
	err := row.Scan(&name)
	return err == nil
}

func(h) getExecutedMigration(db *sql.DB) (map[string]bool, error) {
	res := make(map[string]bool)
	sql, args := sq.Select("name").
		From("migration").MustSql()

	rows, err := db.Query(sql, args...)
	if err != nil {
		return nil, fmt.Errorf("migrations::Run : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		name := ""
		err := rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("migrations::Run : %v", err)
		}
		res[name] = true
	}
	return res, nil
}

func(h) migrate(db *sql.DB, m migration) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("migrate : name:%v : %v", m.Name,err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
			return
		}
		err = errors.Join(err, tx.Commit())
	}()

	_, err = tx.Exec(m.Migration)
	if err != nil {
		return fmt.Errorf("migrate : name:%v : %v", m.Name,err)
	}

	//add migration to migrated
	sql, args := sq.Insert("migration").
		SetMap(map[string]interface{}{
			"name": m.Name,
		}).MustSql()

	_, err = tx.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("migrate : name:%v : %v", m.Name,err)
	}
	return nil
}
