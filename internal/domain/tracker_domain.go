package domain

import (
	"context"
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
	Amount int
}

type TrackerRepository interface {
	Set(c context.Context, dto SetTrackerDto) error
	Index(c context.Context, year int, month int) ([]Tracker, error)
}