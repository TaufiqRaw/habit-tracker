package domain

import (
	"time"
)

const HabitTableName = "habit"

type HabitColsType struct {
	ID string
	Name string
	Amount string
	Unit string
	ArchivedAt string
}

var HabitCols = NewColumnDataContainer(HabitColsType{
	ID: "id",
	Name: "name",
	Amount: "amount",
	Unit: "unit",
	ArchivedAt: "archived_at",
})

type Habit struct {
	ID int64
	Name string
	Amount uint
	Unit string
	ArchivedAt *time.Time
}

type CreateHabitDTO struct {
	Name string
	Amount uint
	Unit string
}