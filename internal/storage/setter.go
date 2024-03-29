package storage

import (
	"postgresHelper/internal/model"
)

func (s *Storage) SetKnobs(knobs []model.Knob) {
	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.knobs = knobs
	}
}
