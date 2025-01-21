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

type habitRepository struct {
	cols domain.HabitColsType
    db *sql.DB
}

func CreateHabitRepository(db *sql.DB) domain.HabitRepository  {
	return &habitRepository{
		cols: domain.HabitCols.Str,
        db: db,
	}
}

func (r *habitRepository) Index(c context.Context, page uint64, limit uint64, unarchived bool) ([]domain.Habit, error) {
    habits := []domain.Habit{}

	preSql := sq.
		Select(domain.HabitCols.AsArray...).
		From(domain.HabitTableName)

	if unarchived {
		preSql = preSql.
			Where(sq.Expr(r.cols.ArchivedAt + " IS NULL"))
	}

    s, args := preSql.
        OrderBy(r.cols.Id).
		Limit(limit).
		Offset(page * limit).MustSql()

    rows, err := r.db.QueryContext(c, s, args...)
    if err != nil {
        return nil, fmt.Errorf("Habit::Index : %v", err)
    }
    defer rows.Close()

    for rows.Next() {
        habit, err := r.scanOneHabit(rows)
		if err != nil {
			return nil, fmt.Errorf("Habit::scanOneHabit : %v", err)
		}
        habits = append(habits, habit)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("Habit::Index : %v", err)
    }
    return habits, nil
}

func (r *habitRepository) Create(c context.Context,dto domain.CreateHabitDTO) (*domain.Habit, error) {
	createdAt := time.Now()

	valueMap := map[string]interface{}{
		r.cols.Name : dto.Name,
		r.cols.Amount : dto.Amount,
		r.cols.Unit : dto.Unit,
		r.cols.RestDay : dto.RestDay,
		r.cols.RestDayMode : dto.RestDayMode,
		r.cols.StartAt : createdAt.String()[:10],
	}

	if dto.LastHabitID != nil {
		valueMap[r.cols.LastHabitID] = *dto.LastHabitID
	}

    sql, args := sq.
        Insert(domain.HabitTableName).
        SetMap(valueMap).MustSql()

    res, err := r.db.ExecContext(c, sql, args...)
    if err != nil {
        return nil, fmt.Errorf("Habit::Create : %v", err)
    }
    ID, err := res.LastInsertId()
    if err != nil {
        return nil, fmt.Errorf("Habit::Create : %v", err)
    }
    return &domain.Habit{
        Id: ID,
        Name: dto.Name,
        Amount: uint(dto.Amount),
        Unit: dto.Unit,
        ArchivedAt: nil,
		RestDay: uint(dto.RestDay),
		RestDayMode: dto.RestDayMode,
		StartAt: createdAt,
		LastHabitId: nil,
    }, nil
}

// Can return nil habit
func (r *habitRepository) getOne(c context.Context,id int64) (*domain.Habit, error) {
	s, args := sq.
        Select(domain.HabitCols.AsArray...).
        From(domain.HabitTableName).
		Where(sq.Eq{
			r.cols.Id : id,
		}).MustSql()
	
	row := r.db.QueryRowContext(c, s, args...)
	habit, err := r.scanOneHabit(row)
	if err != nil {
		e := errors.Unwrap(err)
		if errors.Is(e, sql.ErrNoRows){
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

// will throw error if habit doesnt exist
func (r *habitRepository) getNode(c context.Context, id int64) (*domain.HabitNode, error) {
	habitNode := &domain.HabitNode{}

	habit, err := r.getOne(c, id)
	if err != nil {
		return nil, fmt.Errorf("Habit::GetNode : %v", err)
	}
	if habit == nil {
		//TODO: move the habit not found error to its own global var
		return nil, fmt.Errorf("Habit::GetNode : %v", errors.New("habit not found"))
	}
	habitNode.Habit = *habit

	var prevHabit *domain.Habit 
	if habit.LastHabitId != nil {
		prevHabit, err = r.getOne(c, *habit.LastHabitId)
		if err != nil {
			return nil, fmt.Errorf("Habit::GetNode : %v", err)
		}
	}
	habitNode.PreviousHabit = prevHabit 

	var nextHabit domain.Habit
	nhSql, nhArgs := sq.
		Select(domain.HabitCols.AsArray...).
		Where(sq.Eq{
			r.cols.LastHabitID : id,
		}).MustSql()
	nhRow := r.db.QueryRowContext(c, nhSql, nhArgs...)
	nextHabit, err = r.scanOneHabit(nhRow)
	if err != nil {
		e := errors.Unwrap(err)
		if errors.Is(e, sql.ErrNoRows){
			habitNode.NextHabit = nil
		} else{
			return nil, fmt.Errorf("Habit::GetNode : %v", err)
		}
	}
	habitNode.NextHabit = &nextHabit
	return habitNode, nil
}

func (r *habitRepository) Update(c context.Context, dto domain.UpdateHabitDTO) (*domain.Habit, error) {
	oldHabit, err := r.getOne(c, dto.ID)
	if err != nil {
		return nil, fmt.Errorf("Habit::Update : %v", err)
	}
	if oldHabit == nil {
		return nil, nil
	}
	//if start_at is now() just change it, otherwise create new habit with dto's id as last_habit_id
	if oldHabit.StartAt.String()[:10] == time.Now().String()[:10] {
		updateMap := make(map[string]interface{})
		if dto.Amount != nil {
			updateMap[r.cols.Amount] = *dto.Amount
		}
		if dto.Unit != nil {
			updateMap[r.cols.Unit] = *dto.Unit
		} 
        if dto.RestDay != nil {
            updateMap[r.cols.RestDay] = *dto.RestDay
        }
        if dto.RestDayMode != nil {
            updateMap[r.cols.RestDayMode] = *dto.RestDayMode
        }
		//TODO: Add rest day and rest mode

		s, args := sq.
			Update(domain.HabitTableName).
			SetMap(updateMap).
			Where(sq.Eq{
				r.cols.Id : dto.ID,
			}).MustSql()
		_, err := r.db.ExecContext(c, s, args...)
		if err != nil {
			return nil, fmt.Errorf("Habit::Update : %v", err)
		}
		return nil, nil
	}else {
        createDto := domain.CreateHabitDTO{
            LastHabitID: &dto.ID,
            Name : oldHabit.Name,
        }

        //Amount
        if dto.Amount != nil {
            createDto.Amount = *dto.Amount
        }else {
            createDto.Amount = int(oldHabit.Amount)
        }
        //Unit
        if dto.Unit != nil {
            createDto.Unit = *dto.Unit
        } else {
            createDto.Unit = oldHabit.Unit
        }
        //Rest Day
        if dto.RestDay != nil {
            createDto.RestDay = *dto.RestDay
        } else {
            createDto.RestDay = int(oldHabit.RestDay)
        }
        //Rest Day Mode
        if dto.RestDayMode != nil {
            createDto.RestDayMode = *dto.RestDayMode
        } else {
            createDto.RestDayMode = oldHabit.RestDayMode
        }

        habit, err := r.Create(c, createDto)
        if err != nil {
            return nil, fmt.Errorf("Habit::Update : %v", err)
        }

        return habit, nil
    }
}

func (r *habitRepository) UpdateName(c context.Context, id int64, name string) error {
    s, args := sq.
		Update(domain.HabitTableName).
		SetMap(map[string]interface{}{
            r.cols.Name : name,
        }).
		Where(sq.Eq{
			r.cols.Id : id,
		}).MustSql()
	_, err := r.db.ExecContext(c, s, args...)
	if err != nil {
		return fmt.Errorf("Habit::UpdateName : %v", err)
	}
	return nil
}

func (r *habitRepository) ToggleArchived(c context.Context, id int64) (*domain.Habit, error) {
    isArchived := false
	startAt := ""
    {
		sql, args := sq.
			Select(fmt.Sprintf(` 
				case 
					when %[1]v is null then false
					else true
				end as %[1]v`, r.cols.ArchivedAt), 
				r.cols.StartAt).
			From(domain.HabitTableName).
			Where(sq.Eq{
				r.cols.Id: id,
			}).MustSql()

		row := r.db.QueryRowContext(c, sql, args...)
		err := row.Scan(&isArchived, &startAt)
		if err != nil {
            return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
	}

	//if already archived and not created now, create new habit but with same values
	if isArchived && startAt != time.Now().String()[:10] {
		habit, err := r.getOne(c, id)
		if err != nil {
			return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
		if habit == nil {
			return nil, nil
		}
		_, err = r.Create(c, domain.CreateHabitDTO{
			Name: habit.Name,
			Amount: int(habit.Amount),
			Unit: habit.Unit,
			LastHabitID: &id,
		})
		if err != nil {
			return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		} 
		return habit, nil
	}

    newArchivedAt := map[string]interface{}{
        r.cols.ArchivedAt : time.Now().String()[:10],
    }
    
	sql, args := sq.
		Update(domain.HabitTableName).
		SetMap(newArchivedAt).
		MustSql()
	_, err := r.db.ExecContext(c, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
	}
	return nil, nil
}

func (r *habitRepository) Delete(c context.Context, id int64) error {
	//TODO: if there is habit (A) that have this habit as last habit and this habit has last habit (B), change habit (A) last habit to habit (B)
	habit, err := r.getNode(c, id)
	if err != nil {
		return fmt.Errorf("Habit::Delete : %v", err)
	}
	if habit.PreviousHabit != nil && habit.NextHabit != nil {
		s, args := sq.
			Update(domain.HabitTableName).
			SetMap(map[string]interface{}{
				r.cols.LastHabitID : habit.PreviousHabit.Id,
			}).Where(sq.Eq{
				r.cols.Id : habit.NextHabit.Id,
			}).MustSql()
		_, err := r.db.ExecContext(c, s, args...)
		if err != nil {
			return fmt.Errorf("Habit::Delete : %v", err)
		}
	}

	sql, args := sq.
        Delete(domain.HabitTableName).
        Where(sq.Eq{
            r.cols.Id : id,
        }).MustSql()

	_, err = r.db.ExecContext(c, sql, args...)
	if err != nil {
		return fmt.Errorf("Habit::Delete : %v", err)
	}
	return nil
}

func (r *habitRepository) scanOneHabit(row domain.Scannable) (domain.Habit, error) {
	habit := domain.Habit{}
    var archivedAtStr sql.NullString
	var startAtStr string
	var lastHabitID sql.NullInt64
    err := row.Scan(&habit.Id, &lastHabitID, &habit.Name, 
			&habit.Amount, &habit.Unit, 
			&habit.RestDay, &habit.RestDayMode, 
			&startAtStr, &archivedAtStr)

	if err != nil {
        return domain.Habit{}, fmt.Errorf("Habit::Index : %v", err)
    }
	
	//last_habit_id nullble check
	if lastHabitID.Valid {
		habit.LastHabitId = &lastHabitID.Int64
	}
	//start_at nullble check
	habit.StartAt, err = time.Parse("2006-01-02", startAtStr)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("Habit::Index : %v", err)
	}
	//archived_at nullble check
    if archivedAtStr.Valid {
        t, err := time.Parse("2006-01-02", archivedAtStr.String)
        if err != nil {
            return domain.Habit{}, fmt.Errorf("Habit::scanOneHabit : %v", err)
        }
        habit.ArchivedAt = &t
    }
	return habit, nil
}