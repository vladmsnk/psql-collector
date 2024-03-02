package runner

import (
	"context"
	"fmt"
	"log"
	"postgresHelper/collector"
	"time"
)

const (
	defaultCollectInterval = time.Second * 10
)

type Runner interface {
	Run(ctx context.Context)
}

type Implementation struct {
	collect         Collector
	collectInterval time.Duration
}

type Collector interface {
	CollectKnobs(ctx context.Context) ([]collector.Knob, error)
	SetKnobs([]collector.Knob)
}

func New(collect Collector) *Implementation {
	return &Implementation{
		collect:         collect,
		collectInterval: defaultCollectInterval,
	}
}

func (i *Implementation) Run(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(i.collectInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				knobs, err := i.collect.CollectKnobs(ctx)
				if err != nil {
					log.Println(err)
					continue
				}
				i.collect.SetKnobs(knobs)
				fmt.Println("Collected knobs")
			}
		}
	}()
}
