package domain

import (
	"time"
)

const TrackerTableName = "tracker"

type TrackerColsType struct {
	HabitID string
	Amount string
	At string
}

var TrackerCols = NewColumnDataContainer(TrackerColsType{
	HabitID: "habit_id",
	Amount: "amount",
	At: "at",
})

type Tracker struct {
	HabitID int64
	Amount uint
	At time.Time
}

type SetTrackerDto struct {
	HabitID int64
	Amount uint
}