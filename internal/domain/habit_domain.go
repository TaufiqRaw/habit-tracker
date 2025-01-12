package domain

import (
	"time"
)

const HabitTableName = "habit"

type HabitColsType struct {
	Id string
	LastHabitID string
	Name string
	Amount string
	Unit string
	StartAt string
	ArchivedAt string
}

var HabitCols = NewColumnDataContainer(HabitColsType{
	Id: "id",
	LastHabitID : "last_habit_id",
	Name: "name",
	Amount: "amount",
	Unit: "unit",
	StartAt : "start_at",
	ArchivedAt: "archived_at",
})

type Habit struct {
	Id int64
	LastHabitId *int64
	Name string
	Amount uint
	Unit string
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
	Amount uint
	Unit string
	LastHabitID *int64
}

type UpdateHabitDTO struct {
	ID int64
	Name *string
	Amount *uint
	Unit *string
}

type HabitRepository interface {
	Index(page uint64, limit uint64, unarchived bool) ([]Habit, error)
	Create(dto CreateHabitDTO) (*Habit, error)
	GetNode(id int64) (*HabitNode, error)
	Update(dto UpdateHabitDTO) error
	ToggleArchived(id int64) error
	Delete(id int64) error
}