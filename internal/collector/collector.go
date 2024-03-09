package collector

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/samber/lo"
	"postgresHelper/internal/model"
)

type Collector interface {
	CollectKnobs(ctx context.Context) ([]model.Knob, error)
	CollectMetrics(ctx context.Context)
}

type Implementation struct {
	db *sql.DB
}

func NewCollector(db *sql.DB) *Implementation {
	return &Implementation{db: db}
}

func (i *Implementation) CollectKnobs(ctx context.Context) ([]model.Knob, error) {
	query := `
SELECT name, setting, vartype
From pg_settings
`
	knobs := make(map[string]interface{})

	rows, err := i.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("i.db.QueryContext: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			name, setting, vartype string
			value                  interface{}
		)
		err := rows.Scan(&name, &setting, &vartype)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		switch vartype {
		case "enum":
			//on current implementation step, does not collect enum settings.
			continue
		case "bool":
			if setting == "on" {
				value = true
			} else if setting == "off" {
				value = false
			} else {
				err = fmt.Errorf("unknown value=%s for name=%s", setting, name)
			}
		case "integer", "real":
			value, err = strconv.ParseFloat(setting, 64)
		case "string":
			value = setting
		default:
			err = fmt.Errorf("unknown type=%s for name=%s", vartype, name)
		}
		if err != nil {
			log.Println(err)
			continue
		}
		knobs[name] = value
	}

	res := lo.MapToSlice(knobs, func(key string, value interface{}) model.Knob {
		return model.Knob{Name: key, Value: value}
	})

	return res, nil
}

func (i *Implementation) CollectMetrics(ctx context.Context) {

}

func (i *Implementation) CollectQueryTypesDistribution(ctx context.Context) (map[QueryType]int64, error) {
	rows, err := i.db.QueryContext(ctx, SelectQueryTypesDistribution)
	if err != nil {
		return nil, fmt.Errorf("i.db.QueryContext: %w", err)
	}
	defer rows.Close()

	distribution := make(map[QueryType]int64)
	for rows.Next() {
		var (
			queryType string
			count     int64
		)
		err := rows.Scan(&queryType, &count)
		if err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		distribution[ToQueryType(queryType)] = count
	}
	return distribution, nil
}
