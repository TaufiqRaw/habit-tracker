package service

import (
	"context"
	"habit-tracker/internal/domain"
)

type HabitNodeService struct {
	hnRepo domain.HabitNodeRepository
	ctx context.Context
}

func CreateHabitNodeService(hnRepo domain.HabitNodeRepository) *HabitNodeService {
	return &HabitNodeService{
		hnRepo: hnRepo,
	}
}

type habitNodeResult struct {
	Data *domain.HabitNode
	Err  *string
}

type habitNodeArrResult struct {
	Data []domain.HabitNode
	Err *string
}

func (s *HabitNodeService) startup(c context.Context){
    s.ctx = c
}

func (s *HabitNodeService) Create(dto domain.CreateHabitNodeDTO) habitNodeResult {
	if dto.MinPerDay < 0 || dto.RestDay < 0{
		e:="HabitService::Create : Invalid dto"
		return habitNodeResult{
			Data: nil,
			Err: &e,
		}
	}
	h, err := s.hnRepo.Create(s.ctx, dto, nil)
	if err != nil {
		e := err.Error()
		return habitNodeResult{
			Data: nil,
			Err: &e,
		}
	}
	return habitNodeResult{
		Data: h,
		Err: nil,
	}
}

func (s *HabitNodeService) Update(id int64, dto domain.UpdateHabitNodeDTO) *string {
	err := s.hnRepo.Update(s.ctx, id, dto)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}
