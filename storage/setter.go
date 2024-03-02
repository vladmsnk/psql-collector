package storage

import "postgresHelper/model"

func (s *Storage) SetKnobs(knobs []model.Knob) {
	if s != nil {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.knobs = knobs
	}
}
