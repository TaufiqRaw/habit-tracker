package domain

import (
	"time"
)

const TrackerTableName = "tracker"

type TrackerColsType struct {
	HabitId string
	Amount string
	At string
}

var TrackerCols = NewColumnDataContainer(TrackerColsType{
	HabitId: "habit_id",
	Amount: "amount",
	At: "at",
})

type Tracker struct {
	HabitId int64
	Amount uint
	At time.Time
}

type SetTrackerDto struct {
	HabitId int64
	Amount uint
}

type TrackerRepository interface {
	Set(dto SetTrackerDto) error
	Index(year int, month int) ([]Tracker, error)
}