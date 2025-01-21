package service

import (
	"context"
	"habit-tracker/internal/domain"
)

type HabitService struct {
	hRepo domain.HabitRepository
	ctx context.Context
}

func CreateHabitService(hRepo domain.HabitRepository) *HabitService {
	return &HabitService{
		hRepo: hRepo,
	}
}

type habitResult struct{
	Data *domain.Habit
	Err *string
}

type habitArrResult struct {
	Data []domain.Habit
	Err *string
}

func (s *HabitService) startup(c context.Context){
    s.ctx = c
}

func (s *HabitService) Create(dto domain.CreateHabitDTO) habitResult {
	if dto.Amount < 0 || dto.RestDay < 0{
		e:="HabitService::Create : Invalid dto"
		return habitResult{
			Data: nil,
			Err: &e,
		}
	}
	h, err := s.hRepo.Create(s.ctx, dto)
	if err != nil {
		e := err.Error()
		return habitResult{
			Data: nil,
			Err: &e,
		}
	}
	return habitResult{
		Data: h,
		Err: nil,
	}
}

func (s *HabitService) Index(page int64, limit int64, unarchived bool) habitArrResult {
	if page < 1 || limit < 1 {
		e := "HabitService::Index : pagination invalid"
		return habitArrResult{
			Data: nil,
			Err: &e,
		}
	}
	habits, err := s.hRepo.Index(s.ctx, uint64(page), uint64(limit), unarchived)
	if err != nil {
		e := err.Error()
		return habitArrResult{
			Data: nil,
			Err: &e,
		}
	}
	return habitArrResult{
		Data: habits,
		Err: nil,
	}
}

// func (s *HabitService) GetNode(id int64) habitNodeResult {
// 	n, err := s.hRepo.GetNode(id)
// 	if err != nil {
// 		e := err.Error()
// 		return habitNodeResult{
// 			Data: nil,
// 			Err: &e,
// 		}
// 	}
// 	return habitNodeResult{
// 		Data: n,
// 		Err: nil,
// 	}
// }

func (s *HabitService) Update(dto domain.UpdateHabitDTO) habitResult {
	h, err := s.hRepo.Update(s.ctx, dto)
	if err != nil {
		e := err.Error()
		return habitResult{
			Data: nil,
			Err: &e,
		}
	}
	return habitResult{
		Data: h,
		Err: nil,
	}
}

func (s *HabitService) UpdateName(id int64, name string) *string {
	err := s.hRepo.UpdateName(s.ctx, id, name)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}

func (s *HabitService) ToggleArchived(id int64) habitResult {
	h, err := s.hRepo.ToggleArchived(s.ctx, id)
	if err != nil {
		e := err.Error()
		return habitResult{
			Data: nil,
			Err: &e,
		}
	}
	return habitResult{
		Data: h,
		Err: nil,
	}
}

func (s *HabitService) Delete(id int64) *string {
	err := s.hRepo.Delete(s.ctx, id)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}