package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FishSizeGradeService --output=./mocks --outpkg=service --filename=fish_size_grade_service.go --structname=MockFishSizeGradeService --with-expecter=false
type FishSizeGradeService interface {
	GetDropdown() ([]*dto.DropdownItem, error)
}

type fishSizeGradeService struct {
	fishSizeGradeRepo repository.FishSizeGradeRepository
}

func NewFishSizeGradeService(fishSizeGradeRepo repository.FishSizeGradeRepository) FishSizeGradeService {
	return &fishSizeGradeService{
		fishSizeGradeRepo: fishSizeGradeRepo,
	}
}

func (s *fishSizeGradeService) GetDropdown() ([]*dto.DropdownItem, error) {
	grades, err := s.fishSizeGradeRepo.List()
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	items := make([]*dto.DropdownItem, 0, len(grades))
	for _, g := range grades {
		items = append(items, &dto.DropdownItem{
			Key:   g.Id,
			Value: g.Name,
		})
	}
	return items, nil
}
