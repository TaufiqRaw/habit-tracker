package domain

import (
	"context"
	"database/sql"
	"time"
)

const HabitNodeTableName = "habit_node"

type RestDayModeEnum string

const (
	Week  RestDayModeEnum = "WEEK"
	Month RestDayModeEnum = "MONTH"
)

var AllRestDayMode = []struct {
	Value  RestDayModeEnum
	TSName string
}{
	{Week, "Week"},
	{Month, "Month"},
}

type HabitNodeColsType struct {
	Id          string
	HabitId 	string
	MinPerDay      string
	Unit        string
	RestDay     string
	RestDayMode string
	StartAt     string
	ArchivedAt  string
}

var HabitNodeCols = NewColumnDataContainer(HabitNodeColsType{
	Id:          "id",
	HabitId: 	 "habit_id",
	MinPerDay:   "min_per_day",
	Unit:        "unit",
	RestDay:     "rest_day",
	RestDayMode: "rest_day_mode",
	StartAt:     "start_at",
	ArchivedAt:  "archived_at",
})

type HabitNode struct {
	Id          int64
	MinPerDay      uint
	Unit        string
	RestDay     uint
	RestDayMode RestDayModeEnum
	StartAt     time.Time
	ArchivedAt  *time.Time
}

type CreateHabitNodeDTO struct {
	MinPerDay int
	Unit string
	RestDay int
	RestDayMode RestDayModeEnum
	HabitId int64
}

type UpdateHabitNodeDTO struct {
	MinPerDay *int
	Unit *string
	RestDay *int
	RestDayMode *RestDayModeEnum
}

type HabitNodeRepository interface {
	Create(c context.Context, dto CreateHabitNodeDTO, tx *sql.Tx) (*HabitNode, error)
	Update(c context.Context, id int64, dto UpdateHabitNodeDTO) (error)
	// Delete(c context.Context, id int64) error
}