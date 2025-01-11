package repository

import (
	"database/sql"
	"fmt"
	"habit-tracker/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type habitRepository struct {
	cols domain.HabitColsType
    db *sql.DB
}

func CreateHabitRepository(db *sql.DB) *habitRepository {
	return &habitRepository{
		cols: domain.HabitCols.Str,
        db: db,
	}
}

func (r *habitRepository) Index(page uint64, limit uint64) ([]domain.Habit, error) {
    habits := []domain.Habit{}

    s, args := sq.
        Select(domain.HabitCols.AsArray...).
        From(domain.HabitTableName).
        OrderBy(r.cols.ID).
		Limit(limit).
		Offset(page * limit).MustSql()

    rows, err := r.db.Query(s, args...)
    if err != nil {
        return nil, fmt.Errorf("Habit::Index : %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        habit := domain.Habit{}
        var archivedAtStr sql.NullString
        err := rows.Scan(&habit.ID, &habit.Name, &habit.Amount, &habit.Unit, &archivedAtStr)

        if archivedAtStr.Valid {
            t, err := time.Parse("2006-01-02", archivedAtStr.String)
            if err != nil {
                return nil, fmt.Errorf("Habit::Index : %v", err)
            }
            habit.ArchivedAt = &t
        }

        if err != nil {
            return nil, fmt.Errorf("Habit::Index : %v", err)
        }
        habits = append(habits, habit)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Habit::Index : %v", err)
    }
    return habits, nil
}

func (r *habitRepository) Create(dto domain.CreateHabitDTO) (*domain.Habit, error) {
    sql, args := sq.
        Insert(domain.HabitTableName).
        SetMap(map[string]interface{}{
            r.cols.Name : dto.Name,
            r.cols.Amount : dto.Amount,
            r.cols.Unit : dto.Unit,
        }).MustSql()

    res, err := r.db.Exec(sql, args...)
    if err != nil {
        return nil, fmt.Errorf("Habit::Create : %v", err)
    }
    ID, err := res.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("Habit::Create : %v", err)
    }
    return &domain.Habit{
        ID: ID,
        Name: dto.Name,
        Amount: dto.Amount,
        Unit: dto.Unit,
        ArchivedAt: nil,
    }, nil
}

func (r *habitRepository) ToggleArchived(ID int64) error {
    isArchived := false
    {
		sql, args := sq.
			Select(fmt.Sprintf(` 
				case 
					when %[1]v is null then false
					else true
				end as %[1]v`, r.cols.ArchivedAt)).
			From(domain.HabitTableName).
			Where(sq.Eq{
				r.cols.ID: ID,
			}).MustSql()

		row := r.db.QueryRow(sql, args...)
		err := row.Scan(&isArchived)
		if err != nil {
            return fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
	}
    newArchivedAt := map[string]interface{}{
        r.cols.ArchivedAt : time.Now().String()[:10],
    }
    if isArchived {
		newArchivedAt = map[string]interface{}{
            r.cols.ArchivedAt : nil,
        }
	}
	sql, args := sq.
		Update(domain.HabitTableName).
		SetMap(newArchivedAt).
		MustSql()
	_, err := r.db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("Habit::ToggleArchived : %v", err)
	}
	return nil
}

func (r *habitRepository) Delete(ID int64) error {
	sql, args := sq.
        Delete(domain.HabitTableName).
        Where(sq.Eq{
            r.cols.ID : ID,
        }).MustSql()

	_, err := r.db.Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("Habit::Delete : %v", err)
	}
	return nil
}