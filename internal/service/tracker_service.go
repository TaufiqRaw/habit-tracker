package service

import "habit-tracker/internal/domain"

type TrackerService struct {
	tRepo domain.TrackerRepository
}

func CreateTrackerService(tRepo domain.TrackerRepository) *TrackerService {
	return &TrackerService{
		tRepo: tRepo,
	}
}

func (s *TrackerService) Index(year int, month int) ([]domain.Tracker, *string) {
	if year < 1 || month < 1 || month > 12 {
		e := "TrackerService::Index : invalid args"
		return nil, &e
	} 
	d, err := s.tRepo.Index(year, month)
	if err != nil {
		e := err.Error()
		return nil, &e
	}
	return d, nil
}

func (s *TrackerService) Set(dto domain.SetTrackerDto) (*string) {
	if dto.Amount < 0 {
		e := "TrackerService::Set : invalid args"
		return &e
	}
	err := s.tRepo.Set(dto)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}

