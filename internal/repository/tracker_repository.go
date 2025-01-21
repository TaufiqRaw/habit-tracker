package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"habit-tracker/internal/domain"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type trackerRepository struct {
	cols domain.TrackerColsType
	db   *sql.DB
}

func CreateTrackerRepository(db *sql.DB) domain.TrackerRepository {
	return &trackerRepository{
		cols: domain.TrackerCols.Str,
		db: db,
	}
}

func (r *trackerRepository) Set(c context.Context, dto domain.SetTrackerDto) error {
	var affected int64 
	{
		sql, args := sq.
        Update(domain.HabitTableName).
        SetMap(map[string]interface{}{
            r.cols.Amount : dto.Amount,
        }).Where(sq.Eq{
			r.cols.At : sq.Expr(time.Now().String()[:10]),
			r.cols.HabitId : dto.HabitId,
		}).MustSql()

		res, err := r.db.ExecContext(c, sql, args...)
		if err != nil {
			return fmt.Errorf("Tracker::Set : %v", err)
		}
		affected, err = res.RowsAffected()
		if err != nil {
			return fmt.Errorf("Tracker::Set : %v", err)
		}
	}
	
	if affected != 0 {
		return nil
	}

	// create tracker if affected rows is 0
	sql, args := sq.
        Insert(domain.TrackerTableName).
        SetMap(map[string]interface{}{
            r.cols.HabitId : dto.HabitId,
			r.cols.Amount : dto.Amount,
			r.cols.At : time.Now().String()[:10],
        }).MustSql()

    _, err := r.db.ExecContext(c, sql, args...)
	if err != nil {
		return fmt.Errorf("Tracker::Set : %v", err)
	}
	return nil
}

func (r *trackerRepository) Index(c context.Context, year int, month int) ([]domain.Tracker, error) {
    trackers := []domain.Tracker{}

	if month < 1 || month > 12 {
        return nil, fmt.Errorf("Tracker::Index : %v", errors.New("invalid month"))
	}
	targetDate := time.Date(year, time.Month(month), 0,0,0,0,0, time.Local)

    s, args := sq.
        Select(domain.TrackerCols.AsArray...).
        From(domain.TrackerTableName).
		Where(sq.Like{
			r.cols.At : targetDate.String()[:7],
		}).MustSql()

    rows, err := r.db.QueryContext(c, s, args...)
    if err != nil {
        return nil, fmt.Errorf("Tracker::Index : %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        tracker := domain.Tracker{}
        var atStr sql.NullString
        err := rows.Scan(&tracker.HabitId, &tracker.Amount, &atStr)

        if atStr.Valid {
            tracker.At, err = time.Parse("2006-01-02", atStr.String)
            if err != nil {
                return nil, fmt.Errorf("Tracker::Index : %v", err)
            }
        }

        if err != nil {
            return nil, fmt.Errorf("Tracker::Index : %v", err)
        }
        trackers = append(trackers, tracker)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Tracker::Index : %v", err)
    }
    return trackers, nil
}