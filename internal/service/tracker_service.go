package service

import (
	"context"
	"habit-tracker/internal/domain"
)

type TrackerService struct {
	tRepo domain.TrackerRepository
    ctx context.Context
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

func (s *TrackerService) startup(c context.Context){
    s.ctx = c
}

func (s *TrackerService) Index(year int, month int) TrackerArrResult {
	if year < 1 || month < 1 || month > 12 {
		e := "TrackerService::Index : invalid args"
		return TrackerArrResult{
			Data: nil,
			Err: &e,
		}
	} 
	t, err := s.tRepo.Index(s.ctx, year, month)
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
	err := s.tRepo.Set(s.ctx, dto)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}

