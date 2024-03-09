package collector

import (
	"database/sql"
	"time"
)

type QueryType int

const (
	Insert QueryType = iota + 1
	Update
	Delete
	Select
	Other
)

func (q QueryType) String() string {
	switch q {
	case Insert:
		return "INSERT"
	case Update:
		return "UPDATE"
	case Delete:
		return "DELETE"
	case Select:
		return "SELECT"
	default:
		return "OTHER"
	}
}

func ToQueryType(q string) QueryType {
	switch q {
	case "INSERT":
		return Insert
	case "UPDATE":
		return Update
	case "DELETE":
		return Delete
	case "SELECT":
		return Select
	default:
		return Other
	}
}

// TableStat brief information about table
type TableStat struct {
	RelationID   int64
	RelationName string

	NumberOfLiveTuples int64
	NumberOfDeadTuples int64

	NumberOfSeqScans   int64
	NumberOfIndexScans int64

	NumberOfInserts int64
	NumberOfUpdates int64
	NumberOfDeletes int64

	LastVacuumTime     sql.Null[time.Time]
	LastAutoVacuumTime sql.Null[time.Time]
	LastAnalyzeTime    sql.Null[time.Time]
	LastAutoAnalyze    sql.Null[time.Time]
}
