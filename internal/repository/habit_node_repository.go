package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/repository/colscanner"
	"habit-tracker/internal/repository/setmap"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type habitNodeRepository struct {
	hRepo *habitRepository
	cols  domain.HabitNodeColsType
	db    *sql.DB
	csf   colscanner.Factory[
		habitNodeCSHolder, domain.HabitNode,
	]
}

type habitNodeCSHolder struct {
	Id          int64
	MinPerDay   uint
	Unit        string
	RestDay     uint
	RestDayMode domain.RestDayModeEnum
	ArchivedAt  sql.NullString
	StartAt     string
	HabitId     int64
}

func CreateHabitNodeRepository(db *sql.DB) domain.HabitNodeRepository {
	r := &habitNodeRepository{
		cols: domain.HabitNodeCols.Str,
		db:   db,
	}

	r.csf = colscanner.CreateFactory[habitNodeCSHolder, domain.HabitNode](
		func(holder *habitNodeCSHolder, result *domain.HabitNode) colscanner.Binder {
			return colscanner.Binder{
				r.cols.Id: {
					PtrToHolder: &holder.Id,
					ConvToResult: func() error {
						result.Id = holder.Id
						return nil
					},
				},
				r.cols.MinPerDay: {
					PtrToHolder: &holder.MinPerDay,
					ConvToResult: func() error {
						result.MinPerDay = holder.MinPerDay
						return nil
					},
				},
				r.cols.Unit: {
					PtrToHolder: &holder.Unit,
					ConvToResult: func() error {
						result.Unit = holder.Unit
						return nil
					},
				},
				r.cols.RestDay: {
					PtrToHolder: &holder.RestDay,
					ConvToResult: func() error {
						result.RestDay = holder.RestDay
						return nil
					},
				},
				r.cols.RestDayMode: {
					PtrToHolder: &holder.RestDayMode,
					ConvToResult: func() error {
						result.RestDayMode = holder.RestDayMode
						return nil
					},
				},
				r.cols.StartAt: {
					PtrToHolder: &holder.StartAt,
					ConvToResult: func() error {
						var err error = nil
						result.StartAt, err = time.Parse("2006-01-02", holder.StartAt)
						if err != nil {
							return fmt.Errorf("HabitNode::ColsScannerFactory : %v", err)
						}
						return nil
					},
				},
				r.cols.ArchivedAt: {
					PtrToHolder: &holder.ArchivedAt,
					ConvToResult: func() error {
						if holder.ArchivedAt.Valid {
							t, err := time.Parse("2006-01-02", holder.ArchivedAt.String)
							if err != nil {
								return fmt.Errorf("Habit::scanOneHabit : %v", err)
							}
							result.ArchivedAt = &t
						}
						return nil
					},
				},
			}
		},
	)

	return r
}

func (r *habitNodeRepository) setHRepo(hRepo *habitRepository) {
	r.hRepo = hRepo
}

// allow nil tx.
// if tx != nil, then its assumed that the caller is habit repo
// if not, this function will archive previous habit (if any)
func (r *habitNodeRepository) Create(c context.Context, dto domain.CreateHabitNodeDTO, _tx *sql.Tx) (*domain.HabitNode, error) {
	var tx *sql.Tx
	var habit *domain.Habit = nil
	if _tx != nil {
		tx = _tx
	} else {
		_habit, err := r.hRepo.getOne(c, dto.HabitId)
		if err != nil {
			return nil, fmt.Errorf("HabitNode::Create : %v", err)
		}
		habit = &_habit
		habit.LastNode, err = r.getOne(c, habit.NodeIDs[habit.NodeLength-1])
		if err != nil {
			return nil, fmt.Errorf("Habit::getOne : %v", err)
		}
		tx, err = r.db.BeginTx(c, nil)
		if err != nil {
			return nil, fmt.Errorf("HabitNode::Create : %v", err)
		}
		defer func() {
			if err != nil {
				err = errors.Join(err, tx.Rollback())
				return
			}
			err = errors.Join(err, tx.Commit())
		}()
	}

	today := time.Now()

	valueMap := map[string]interface{}{
		r.cols.MinPerDay:   dto.MinPerDay,
		r.cols.Unit:        dto.Unit,
		r.cols.RestDay:     dto.RestDay,
		r.cols.RestDayMode: dto.RestDayMode,
		r.cols.StartAt:     today.String()[:10],
		r.cols.HabitId:     dto.HabitId,
	}

	sql, args := sq.
		Insert(domain.HabitNodeTableName).
		SetMap(valueMap).MustSql()

	res, err := tx.ExecContext(c, sql, args)
	if err != nil {
		return nil, fmt.Errorf("HabitNode::Create : %v", err)
	}
	resId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("HabitNode::Create : %v", err)
	}

	if habit != nil {
		// archive prev habit
		err := r.update(c, map[string]interface{}{
			r.cols.ArchivedAt: today.String()[:10],
		}, sq.Eq{
			r.cols.Id: habit.LastNode.Id,
		}, tx)
		if err != nil {
			return nil, fmt.Errorf("HabitNode::Create : %v", err)
		}

		// update habit's NodeIDs & lastNodeID
		newNodeIDs := append(habit.NodeIDs, resId)
		nodeIDsJson, err := json.Marshal(newNodeIDs)
		if err != nil {
			return nil, fmt.Errorf("HabitNode::Create : %v", err)
		}
		r.hRepo.update(c, map[string]interface{}{
			r.hRepo.cols.NodeIDs:    string(nodeIDsJson),
			r.hRepo.cols.LastNodeID: resId,
		}, sq.Eq{
			domain.HabitCols.Str.Id: habit.Id,
		}, tx)
	}

	return &domain.HabitNode{
		Id:          resId,
		MinPerDay:   uint(dto.MinPerDay),
		Unit:        dto.Unit,
		RestDay:     uint(dto.RestDay),
		RestDayMode: dto.RestDayMode,
		StartAt:     today,
		ArchivedAt:  nil,
	}, nil
}
func (r *habitNodeRepository) Update(c context.Context, id int64, dto domain.UpdateHabitNodeDTO) error {
	setMap := setmap.NewSetMap(nil)
	setMap.SetIfNotNilMap(map[string]interface{}{
		r.cols.Unit:        dto.Unit,
		r.cols.MinPerDay:   dto.MinPerDay,
		r.cols.RestDay:     dto.RestDay,
		r.cols.RestDayMode: dto.RestDayMode,
	})

	err := r.update(c, setMap.GetMap(), sq.Eq{
		r.cols.Id: id,
	}, nil)
	if err != nil {
		return fmt.Errorf("Habit::Update : %v", err)
	}
	return nil
}

// func (r *habitNodeRepository) Delete(c context.Context, id int64) error

func (r *habitNodeRepository) get(c context.Context, where interface{}) ([]domain.HabitNode, error) {
	result := make([]domain.HabitNode, 1)

	cs := r.csf.SelectCols(nil, domain.HabitNodeCols.AsArray...)
	sql, args := sq.Select(cs.Cols...).
		From(domain.HabitNodeTableName).
		Where(where).MustSql()

	rows, err := r.db.QueryContext(c, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("Habit::get : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		_cs := r.csf.From(cs)
		err := colscanner.ScanRow(rows, &_cs)
		if err != nil {
			return nil, fmt.Errorf("Habit::getOne : %v", err)
		}
		res, _ := _cs.GetResult()
		result = append(result, res)
	}

	return result, nil
}

// return err if not found
func (r *habitNodeRepository) getOne(c context.Context, id int64) (domain.HabitNode, error) {
	cs := r.csf.SelectCols(nil, domain.HabitNodeCols.AsArray...)
	sql, args := sq.Select(cs.Cols...).
		From(domain.HabitNodeTableName).
		Where(sq.Eq{
			r.cols.Id: id,
		}).MustSql()

	row := r.db.QueryRowContext(c, sql, args...)
	err := colscanner.ScanRow(row, &cs)
	if err != nil {
		return domain.HabitNode{}, fmt.Errorf("Habit::getOne : %v", err)
	}
	res, _ := cs.GetResult()
	return res, nil
}

// allow nil tx.
// pure updating, no other side effect
func (r *habitNodeRepository) update(c context.Context, setMap map[string]interface{}, where interface{}, tx *sql.Tx) error {
	var db execable
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	sql, args := sq.
		Update(domain.HabitNodeTableName).
		SetMap(setMap).
		Where(where).MustSql()

	res, err := db.ExecContext(c, sql, args...)
	if err != nil {
		return fmt.Errorf("HabitNode::update : %v", err)
	}
	nAffect, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("HabitNode::update : %v", err)
	}
	if nAffect == 0 {
		return fmt.Errorf("HabitNode::update : %v",
			errors.New("no row affected"),
		)
	}
	return nil
}
