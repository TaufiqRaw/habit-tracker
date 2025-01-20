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

type TrackerArrResult struct {
	Data []domain.Tracker
	Err *string
}

func (s *TrackerService) Index(year int, month int) TrackerArrResult {
	if year < 1 || month < 1 || month > 12 {
		e := "TrackerService::Index : invalid args"
		return TrackerArrResult{
			Data: nil,
			Err: &e,
		}
	} 
	t, err := s.tRepo.Index(year, month)
	if err != nil {
		e := err.Error()
		return TrackerArrResult{
			Data: nil,
			Err: &e,
		}
	}
	return TrackerArrResult{
		Data: t,
		Err: nil,
	}
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

