package migrations

import (
	"habit-tracker/internal/domain"
	"strings"
)

func init_table() migration {
	sqls := make([]string, 4)
	{
		r := domain.HabitCols.Str
		sqls = append(sqls,
			"CREATE TABLE IF NOT EXISTS "+domain.HabitTableName+" (\n"+
				r.Id+" INTEGER PRIMARY KEY AUTOINCREMENT, \n"+
				r.Name+" VARCHAR(50) NOT NULL, \n"+
				r.NodeIDs+" TEXT \n"+
				");")
	}
	{
		r := domain.HabitNodeCols.Str
		sqls = append(sqls,
			"CREATE TABLE IF NOT EXISTS "+domain.HabitNodeTableName+" (\n"+
				r.Id+" INTEGER PRIMARY KEY AUTOINCREMENT, \n"+
				r.HabitId+" INTEGER NOT NULL, "+
				r.MinPerDay+" INTEGER NOT NULL, \n"+
				r.RestDay+" INTEGER DEFAULT 0, \n"+
				r.RestDayMode+" VARCHAR(50) NOT NULL, \n"+
				r.StartAt+" CHARACTER(10) NOT NULL, \n"+
				r.Unit+" VARCHAR(50) NOT NULL, \n"+
				r.ArchivedAt+" CHARACTER(10), \n"+
				"FOREIGN KEY("+r.HabitId+") REFERENCES "+
				domain.HabitTableName+"("+domain.HabitCols.Str.Id+")\n"+
				");")
	}
	{
		r := domain.HabitCols.Str
		sqls = append(sqls,
			"ALTER TABLE "+domain.HabitTableName+
				" ADD COLUMN "+r.LastNodeID+
				" INTEGER REFERENCES "+domain.HabitNodeTableName+
				"("+domain.HabitNodeCols.Str.Id+");")
	}
	{
		sqls = append(sqls,
			"CREATE TABLE IF NOT EXISTS "+domain.TrackerTableName+" (\n"+
				domain.TrackerCols.Str.HabitId+" INTEGER, \n"+
				domain.TrackerCols.Str.At+" CHARACTER(10), \n"+
				domain.TrackerCols.Str.Amount+" INTEGER, \n"+
				"PRIMARY KEY ("+domain.TrackerCols.Str.HabitId+", "+
				domain.TrackerCols.Str.At+"), \n"+
				"FOREIGN KEY("+domain.TrackerCols.Str.HabitId+") REFERENCES "+
				domain.HabitTableName+"("+domain.HabitCols.Str.Id+
				") ON DELETE CASCADE ON UPDATE CASCADE"+
				");")
	}
	return migration{
		"init_table",
		strings.Join(sqls, "\n"),
	}
}