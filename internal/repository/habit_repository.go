package repository

import (
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

func CreateHabitRepository(db *sql.DB) domain.HabitRepository {
	return &habitRepository{
		cols: domain.HabitCols.Str,
        db: db,
	}
}

func (r *habitRepository) Index(page uint64, limit uint64, unarchived bool) ([]domain.Habit, error) {
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

    rows, err := r.db.Query(s, args...)
    if err != nil {
        return nil, fmt.Errorf("Habit::Index : %v", err)
    }
    defer rows.Close()

    for rows.Next() {
		rows.Scan()
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

func (r *habitRepository) Create(dto domain.CreateHabitDTO) (*domain.Habit, error) {
	createdAt := time.Now()

	valueMap := map[string]interface{}{
		r.cols.Name : dto.Name,
		r.cols.Amount : dto.Amount,
		r.cols.Unit : dto.Unit,
		r.cols.StartAt : createdAt.String()[:10],
	}

	if dto.LastHabitID != nil {
		valueMap[r.cols.LastHabitID] = *dto.LastHabitID
	}

    sql, args := sq.
        Insert(domain.HabitTableName).
        SetMap(valueMap).MustSql()

    res, err := r.db.Exec(sql, args...)
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
        Amount: dto.Amount,
        Unit: dto.Unit,
        ArchivedAt: nil,
		StartAt: createdAt,
		LastHabitId: nil,
    }, nil
}

// Can return nil habit
func (r *habitRepository) getOne(id int64) (*domain.Habit, error) {
	s, args := sq.
        Select(domain.HabitCols.AsArray...).
        From(domain.HabitTableName).
		Where(sq.Eq{
			r.cols.Id : id,
		}).MustSql()
	
	row := r.db.QueryRow(s, args...)
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
func (r *habitRepository) GetNode(id int64) (*domain.HabitNode, error) {
	habitNode := &domain.HabitNode{}

	habit, err := r.getOne(id)
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
		prevHabit, err = r.getOne(*habit.LastHabitId)
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
	nhRow := r.db.QueryRow(nhSql, nhArgs...)
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

func (r *habitRepository) Update(dto domain.UpdateHabitDTO) error {
	oldHabit, err := r.getOne(dto.ID)
	if err != nil {
		return fmt.Errorf("Habit::Update : %v", err)
	}
	if oldHabit == nil {
		return nil
	}
	//if start_at is now() just change it, otherwise create new habit with dto's id as last_habit_id
	if oldHabit.StartAt.String()[:10] == time.Now().String()[:10] {
		updateMap := make(map[string]interface{})
		if dto.Name != nil {
			updateMap[r.cols.Name] = *dto.Name
		}
		if dto.Amount != nil {
			updateMap[r.cols.Amount] = *dto.Amount
		}
		if dto.Unit != nil {
			updateMap[r.cols.Unit] = *dto.Unit
		} 
		s, args := sq.
			Update(domain.HabitTableName).
			SetMap(updateMap).
			Where(sq.Eq{
				r.cols.Id : dto.ID,
			}).MustSql()
		_, err := r.db.Exec(s, args...)
		if err != nil {
			return fmt.Errorf("Habit::Update : %v", err)
		}
		return nil
	}
	createDto := domain.CreateHabitDTO{
		LastHabitID: &dto.ID,
	}
	if dto.Name != nil {
		createDto.Name = *dto.Name
	} else {
		createDto.Name = oldHabit.Name
	}

	if dto.Amount != nil {
		createDto.Amount = *dto.Amount
	}else {
		createDto.Amount = oldHabit.Amount
	}

	if dto.Unit != nil {
		createDto.Unit = *dto.Unit
	} else {
		createDto.Unit = oldHabit.Unit
	}
	_, err = r.Create(createDto)
	if err != nil {
		return fmt.Errorf("Habit::Update : %v", err)
	}

	return nil
}

func (r *habitRepository) ToggleArchived(id int64) error {
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

		row := r.db.QueryRow(sql, args...)
		err := row.Scan(&isArchived, &startAt)
		if err != nil {
            return fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
	}

	//if already archived and not created now, create new habit but with same values
	if isArchived && startAt != time.Now().String()[:10] {
		habit, err := r.getOne(id)
		if err != nil {
			return fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
		if habit == nil {
			return nil
		}
		_, err = r.Create(domain.CreateHabitDTO{
			Name: habit.Name,
			Amount: habit.Amount,
			Unit: habit.Unit,
			LastHabitID: &id,
		})
		if err != nil {
			return fmt.Errorf("Habit::ToggleArchived : %v", err)
		} 
		return nil
	}

    newArchivedAt := map[string]interface{}{
        r.cols.ArchivedAt : time.Now().String()[:10],
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

func (r *habitRepository) Delete(id int64) error {
	//TODO: if there is habit (A) that have this habit as last habit and this habit has last habit (B), change habit (A) last habit to habit (B)
	habit, err := r.GetNode(id)
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
		_, err := r.db.Exec(s, args...)
		if err != nil {
			return fmt.Errorf("Habit::Delete : %v", err)
		}
	}

	sql, args := sq.
        Delete(domain.HabitTableName).
        Where(sq.Eq{
            r.cols.Id : id,
        }).MustSql()

	_, err = r.db.Exec(sql, args...)
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
    err := row.Scan(&habit.Id, &lastHabitID,&habit.Name, &habit.Amount, &habit.Unit, &startAtStr,&archivedAtStr)

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