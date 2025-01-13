package service

import (
	"habit-tracker/internal/domain"
)

type HabitService struct {
	hRepo domain.HabitRepository
}

func CreateHabitService(hRepo domain.HabitRepository) *HabitService {
	return &HabitService{
		hRepo: hRepo,
	}
}

func (s *HabitService) Create(dto domain.CreateHabitDTO) (*domain.Habit, *string) {
	if dto.Amount < 0 || dto.RestDay < 0{
		e:="HabitService::Create : Invalid dto"
		return nil, &e
	}
	h, err := s.hRepo.Create(dto)
	if err != nil {
		e := err.Error()
		return nil, &e
	}
	return h, nil
}

func (s *HabitService) Index(page int64, limit int64, unarchived bool) ([]domain.Habit, *string) {
	if page < 1 || limit < 1 {
		e := "HabitService::Index : pagination invalid"
		return nil, &e
	}
	habits, err := s.hRepo.Index(uint64(page), uint64(limit), unarchived)
	if err != nil {
		e := err.Error()
		return nil,&e
	}
	return habits, nil
}

func (s *HabitService) GetNode(id int64) (*domain.HabitNode, *string) {
	n, err := s.hRepo.GetNode(id)
	if err != nil {
		e := err.Error()
		return nil, &e
	}
	return n, nil
}

func (s *HabitService) Update(dto domain.UpdateHabitDTO) (*domain.Habit, *string) {
	habit, err := s.hRepo.Update(dto)
	if err != nil {
		e := err.Error()
		return nil,&e
	}
	return habit, nil
}

func (s *HabitService) ToggleArchived(id int64) (*domain.Habit, *string) {
	habit, err := s.hRepo.ToggleArchived(id)
	if err != nil {
		e := err.Error()
		return nil,&e
	}
	return habit, nil
}

func (s *HabitService) Delete(id int64) *string {
	err := s.hRepo.Delete(id)
	if err != nil {
		e := err.Error()
		return &e
	}
	return nil
}