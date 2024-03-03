package psql_helper

import (
	"context"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	"postgresHelper/internal/model"
	desc "postgresHelper/internal/pkg/collector"
)

type Delivery struct {
	desc.CollectorServer
	storage Storager
}

func New(s Storager) *Delivery {
	return &Delivery{
		storage: s,
	}
}

type Storager interface {
	GetKnobs() []model.Knob
}

func (d *Delivery) CollectKnobs(ctx context.Context, _ *desc.CollectKnobsRequest) (*desc.CollectKnobsResponse, error) {
	return ConvertKnobsToDesc(d.storage.GetKnobs())
}

func ConvertKnobsToDesc(knobs []model.Knob) (*desc.CollectKnobsResponse, error) {
	response := &desc.CollectKnobsResponse{}
	var err error

	lo.Map(knobs, func(knob model.Knob, _ int) *desc.CollectKnobsResponse_Knob {
		protoKnob := &desc.CollectKnobsResponse_Knob{
			Name: knob.Name,
		}

		switch v := knob.Value.(type) {
		case string:
			protoKnob.Value = &desc.CollectKnobsResponse_Knob_StrValue{StrValue: v}
		case float64: // Assuming float values in Go are of type float64
			protoKnob.Value = &desc.CollectKnobsResponse_Knob_FloatValue{FloatValue: float32(v)} // Protobuf float is 32-bit
		case bool:
			protoKnob.Value = &desc.CollectKnobsResponse_Knob_BoolValue{BoolValue: v}
		default:
			err = fmt.Errorf("Unsupported type for knob %s: %v\n", knob.Name, reflect.TypeOf(knob.Value))
		}
		return protoKnob
	})
	if err != nil {
		return nil, err
	}

	return response, nil
}
