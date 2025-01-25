package domain

import (
	"context"
)

const HabitTableName = "habit"

type HabitColsType struct {
	Id string
	Name string
    NodeIDs string
    LastNodeID string
}

var HabitCols = NewColumnDataContainer(HabitColsType{
	Id: "id",
	Name: "name",
	NodeIDs: "node_ids",
    LastNodeID: "last_node_id",
})

type Habit struct {
	Id int64
	Name string
    NodeIDs []int64
    NodeLength int
	LastNode HabitNode
}

type HabitAllNodes struct {
    Habit
    Nodes []HabitNode
}

type CreateHabitDTO struct {
	Name string
	MinPerDay int
	Unit string
	RestDay int
	RestDayMode RestDayModeEnum
}

type HabitRepository interface {
	Index(c context.Context, page uint64, limit uint64, unarchived bool) ([]Habit, error)
    Get(c context.Context, id int64) (*HabitAllNodes, error)
	Create(c context.Context, dto CreateHabitDTO) (*Habit, error)
	Update(c context.Context, id int64, name string) (error)
	ToggleArchived(c context.Context, id int64) (*HabitNode, error)
	Delete(c context.Context, id int64) error
}