package domain

import (
	"time"
)

const HabitTableName = "habit"

type RestDayModeEnum string

const (
	Weekly RestDayModeEnum = "Weekly"
	Monthly RestDayModeEnum = "Monthly"
)

var AllRestDayMode = []struct {
    Value  RestDayModeEnum
    TSName string
}{
    {Weekly, "WEEKLY"},
	{Monthly, "MONTHLY"},
}

type HabitColsType struct {
	Id string
	LastHabitID string
	Name string
	Amount string
	Unit string
	RestDay string
	RestDayMode string
	StartAt string
	ArchivedAt string
}

var HabitCols = NewColumnDataContainer(HabitColsType{
	Id: "id",
	LastHabitID : "last_habit_id",
	Name: "name",
	Amount: "amount",
	Unit: "unit",
	RestDay: "rest_day",
	RestDayMode: "rest_day_mode",
	StartAt : "start_at",
	ArchivedAt: "archived_at",
})

type Habit struct {
	Id int64
	LastHabitId *int64
	Name string
	Amount uint
	Unit string
	RestDay uint
	RestDayMode RestDayModeEnum
	StartAt time.Time
	ArchivedAt *time.Time
}

type HabitNode struct {
	Habit
	PreviousHabit *Habit
	NextHabit *Habit
}

type CreateHabitDTO struct {
	Name string
	Amount int
	Unit string
	RestDay int
	RestDayMode RestDayModeEnum
	LastHabitID *int64
}

type UpdateHabitDTO struct {
	ID int64
	Amount *int
	Unit *string
	RestDay *int
	RestDayMode *RestDayModeEnum
}

type HabitRepository interface {
	Index(page uint64, limit uint64, unarchived bool) ([]Habit, error)
	Create(dto CreateHabitDTO) (*Habit, error)
	GetNode(id int64) (*HabitNode, error)
	Update(dto UpdateHabitDTO) (*Habit, error)
	UpdateName(id int64, name string) (error)
	ToggleArchived(id int64) (*Habit, error)
	Delete(id int64) error
}