package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"habit-tracker/internal/domain"
	"habit-tracker/internal/repository/colscanner"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type habitRepository struct {
	cols   domain.HabitColsType
	hnRepo *habitNodeRepository
	db     *sql.DB
	csf    colscanner.Factory[habitCSHolder, domain.Habit]
}

type habitCSHolder struct {
	Id         int64
	Name       string
	NodeIDs    string
	LastNodeID sql.NullInt64
}

func CreateHabitRepository(db *sql.DB, hnRepo domain.HabitNodeRepository) domain.HabitRepository {
	actualHnRepo, ok := hnRepo.(*habitNodeRepository)
	if !ok {
		panic("Invalid Habit Node Repo for SQLite habit repo")
	}
	r := &habitRepository{
		cols:   domain.HabitCols.Str,
		hnRepo: actualHnRepo,
		db:     db,
	}
	actualHnRepo.setHRepo(r)

	r.csf = colscanner.CreateFactory[habitCSHolder, domain.Habit](
		func(holder *habitCSHolder, resHolder *domain.Habit) colscanner.Binder {
			return colscanner.Binder{
				r.cols.Id: {
					PtrToHolder: &holder.Id,
					ConvToResult: func() error {
						resHolder.Id = holder.Id
						return nil
					},
				},
				r.cols.Name: {
					PtrToHolder: &holder.Name,
					ConvToResult: func() error {
						resHolder.Name = holder.Name
						return nil
					},
				},
				r.cols.NodeIDs: {
					PtrToHolder: &holder.NodeIDs,
					ConvToResult: func() error {
						err := json.Unmarshal([]byte(holder.NodeIDs), &resHolder.NodeIDs)
						if err != nil {
							return fmt.Errorf("Habit::Scanner : %v", err)
						}
						resHolder.NodeLength = len(resHolder.NodeIDs)
						return nil
					},
				},
			}
		},
	)

	return r
}

func (r *habitRepository) Create(c context.Context, dto domain.CreateHabitDTO) (_ *domain.Habit, err error) {

	hValueMap := map[string]interface{}{
		r.cols.Name: dto.Name,
	}

	hnDto := domain.CreateHabitNodeDTO{
		MinPerDay:   dto.MinPerDay,
		Unit:        dto.Unit,
		RestDay:     dto.RestDay,
		RestDayMode: dto.RestDayMode,
	}

	tx, err := r.db.BeginTx(c, nil)
	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
			return
		}
		err = errors.Join(err, tx.Commit())
	}()

	// Create Habit
	var hID int64
	{
		sql, args := sq.
			Insert(domain.HabitTableName).
			SetMap(hValueMap).MustSql()

		res, err := tx.ExecContext(c, sql, args...)
		if err != nil {
			return nil, fmt.Errorf("Habit::Create : %v", err)
		}
		hID, err = res.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("Habit::Create : %v", err)
		}

		hnDto.HabitId = hID
	}

	// Create Habit Node
	var hNode *domain.HabitNode
	{
		hNode, err = r.hnRepo.Create(c, hnDto, tx)
		if err != nil {
			return nil, fmt.Errorf("Habit::Create : %v", err)
		}
	}

	// Update Habit's nodes & last_node_id
	{
		hnIDJson, err := json.Marshal([]int64{hNode.Id})
		if err != nil {
			return nil, fmt.Errorf("Habit::Create : %v", err)
		}
		err = r.update(c, map[string]interface{}{
			r.cols.NodeIDs:    string(hnIDJson),
			r.cols.LastNodeID: hNode.Id,
		}, sq.Eq{
			r.cols.Id: hID,
		}, tx)
		if err != nil {
			return nil, fmt.Errorf("Habit::Create : %v", err)
		}
	}

	return &domain.Habit{
		Id:       hID,
		Name:     dto.Name,
		NodeIDs:  []int64{hNode.Id},
		LastNode: *hNode,
	}, nil
}

func (r *habitRepository) Index(c context.Context, page uint64, limit uint64, unarchived bool) ([]domain.Habit, error) {
	habits := []domain.Habit{}

	//table alias
	hAs := "h"
	hnAs := "hn"

	hCs := r.csf.SelectCols(&hAs, domain.HabitCols.AsArray...)
	hnCs := r.hnRepo.csf.SelectCols(&hnAs, domain.HabitNodeCols.AsArray...)

	preSqlSelect := colscanner.CreateSelect(&hCs, &hnCs)
	preSql := sq.
		Select(preSqlSelect...).
		From((domain.HabitTableName + " AS " + hAs) +
			" JOIN " + (domain.HabitNodeTableName + " AS " + hnAs) + " ON " +
			(hAs + "." + r.cols.Id) + " = " + (hnAs + "." + r.hnRepo.cols.HabitId))

	if unarchived {
		col := domain.HabitNodeCols.Str.ArchivedAt
		preSql = preSql.
			Where(sq.Expr(hnAs + "." + col + " IS NULL"))
	}

	s, args := preSql.
		OrderBy(hnAs + "." + r.cols.Id).
		Limit(limit).
		Offset(page * limit).MustSql()

	rows, err := r.db.QueryContext(c, s, args...)
	if err != nil {
		return nil, fmt.Errorf("Habit::Index : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		_hCs := r.csf.From(hCs)
		_hnCs := r.hnRepo.csf.From(hnCs)
		err := colscanner.ScanRow(rows, &_hCs, &_hnCs)
		habit, _ := _hCs.GetResult()
		habitNode, _ := _hnCs.GetResult()
		habit.LastNode = habitNode
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

func (r *habitRepository) Get(c context.Context, id int64) (_ *domain.HabitAllNodes, err error) {
	var h domain.HabitAllNodes
	h.Habit, err = r.getOne(c, id)
	if err != nil {
		return nil, fmt.Errorf("Habit::Get : %v", err)
	}
	h.Nodes, err = r.hnRepo.get(c, sq.Eq{
		r.hnRepo.cols.HabitId: id,
	})
	if err != nil {
		return nil, fmt.Errorf("Habit::Get : %v", err)
	}
	return &h, nil
}

func (r *habitRepository) Update(c context.Context, id int64, name string) error {
	err := r.update(c, map[string]interface{}{
		r.cols.Name: name,
	}, sq.Eq{
		r.cols.Id: id,
	}, nil)

	return err
}


// return new habitNode if unarchiving old habit
func (r *habitRepository) ToggleArchived(c context.Context, id int64) (_ *domain.HabitNode,err error) {
	var habit domain.Habit
	today := time.Now().String()[:10]
	{
		habit, err = r.getOne(c, id)
		if err != nil {
			return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
		habit.LastNode, err = r.hnRepo.getOne(c, habit.NodeIDs[habit.NodeLength-1])
		if err != nil {
			return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
	}
	ln := &habit.LastNode

	//if already archived and not created now, create new habitNode but with same values
	lnStartAt := habit.LastNode.StartAt.String()[:10]
	if (ln.ArchivedAt != nil) && (lnStartAt != today) {
		hn, err := r.hnRepo.Create(c, domain.CreateHabitNodeDTO{
			MinPerDay:   int(ln.MinPerDay),
			Unit:        ln.Unit,
			RestDay:     int(ln.RestDay),
			RestDayMode: ln.RestDayMode,
			HabitId:     habit.Id,
		}, nil)
		if err != nil {
			return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
		}
		return hn,nil
	}

	err = r.hnRepo.update(c, map[string]interface{}{
		r.hnRepo.cols.ArchivedAt: today,
	}, sq.Eq{
		r.hnRepo.cols.Id: ln.Id,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("Habit::ToggleArchived : %v", err)
	}
	return nil, nil
}

func (r *habitRepository) Delete(c context.Context, id int64) error {
	sql, args := sq.
		Delete(domain.HabitTableName).
		Where(sq.Eq{
			r.cols.Id: id,
		}).MustSql()

	_, err := r.db.ExecContext(c, sql, args...)
	if err != nil {
		return fmt.Errorf("Habit::Delete : %v", err)
	}
	return nil
}

// return err if no habit found.
// Doesn't populate lastNode
func (r *habitRepository) getOne(c context.Context, id int64) (domain.Habit, error) {
	cs := r.csf.SelectCols(nil, domain.HabitCols.AsArray...)
	sql, args := sq.Select(cs.Cols...).
		From(domain.HabitTableName).
		Where(sq.Eq{
			r.cols.Id: id,
		}).MustSql()

	row := r.db.QueryRowContext(c, sql, args...)
	err := colscanner.ScanRow(row)
	if err != nil {
		return domain.Habit{}, fmt.Errorf("Habit::getOne : %v", err)
	}
	h, _ := cs.GetResult()
	return h, nil
}

// allow nil tx
func (r *habitRepository) update(c context.Context, setMap map[string]interface{}, where interface{}, tx *sql.Tx) error {
	var db execable
	if tx != nil {
		db = tx
	} else {
		db = r.db
	}

	sql, args := sq.
		Update(domain.HabitTableName).
		SetMap(setMap).
		Where(where).MustSql()

	res, err := db.ExecContext(c, sql, args...)
	if err != nil {
		return fmt.Errorf("Habit::update : %v", err)
	}
	nAffect, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("Habit::update : %v", err)
	}
	if nAffect == 0 {
		return fmt.Errorf("Habit::update : %v",
			errors.New("no row affected"),
		)
	}
	return nil
}
