package service

import (
	"context"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

type FarmGroupService interface {
	Create(ctx context.Context, request dto.CreateFarmGroupRequest) (*dto.FarmGroupResponse, error)
	Get(id int) (*dto.FarmGroupResponse, error)
	Update(ctx context.Context, request dto.UpdateFarmGroupRequest) error
	List(clientId int) ([]*dto.FarmGroupResponse, error)
	GetClientIdByFarmGroupId(farmGroupId int) (int, error)
	GetDropdown(clientId int) ([]*dto.DropdownItem, error)
}

type farmGroupService struct {
	farmGroupRepo repository.FarmGroupRepository
}

func NewFarmGroupService(farmGroupRepo repository.FarmGroupRepository) FarmGroupService {
	return &farmGroupService{
		farmGroupRepo: farmGroupRepo,
	}
}

func hasDuplicateFarmIds(ids []int) bool {
	seen := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			return true
		}
		seen[id] = struct{}{}
	}
	return false
}

func (s *farmGroupService) Create(ctx context.Context, request dto.CreateFarmGroupRequest) (*dto.FarmGroupResponse, error) {
	if hasDuplicateFarmIds(request.FarmIds) {
		return nil, errors.ErrFarmGroupInvalidInput
	}

	fg := &model.FarmGroup{
		ClientId: request.ClientId,
		Name:     request.Name,
	}

	if err := s.farmGroupRepo.Create(ctx, fg); err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	joins := make([]*model.FarmOnFarmGroup, 0, len(request.FarmIds))
	for _, farmId := range request.FarmIds {
		joins = append(joins, &model.FarmOnFarmGroup{
			FarmId:      farmId,
			FarmGroupId: fg.Id,
		})
	}
	if err := s.farmGroupRepo.CreateJoins(ctx, joins); err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.Get(fg.Id)
}

func (s *farmGroupService) Get(id int) (*dto.FarmGroupResponse, error) {
	fg, err := s.farmGroupRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}
	if fg == nil {
		return nil, errors.ErrFarmGroupNotFound
	}

	farms, err := s.farmGroupRepo.GetFarmsByFarmGroupId(fg.Id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toResponse(fg, farms), nil
}

func (s *farmGroupService) Update(ctx context.Context, request dto.UpdateFarmGroupRequest) error {
	if hasDuplicateFarmIds(request.FarmIds) {
		return errors.ErrFarmGroupInvalidInput
	}

	fg, err := s.farmGroupRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if fg == nil {
		return errors.ErrFarmGroupNotFound
	}

	fg.Name = request.Name
	if err := s.farmGroupRepo.Update(ctx, fg); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	if err := s.farmGroupRepo.DeleteJoinsByFarmGroupId(ctx, fg.Id); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	joins := make([]*model.FarmOnFarmGroup, 0, len(request.FarmIds))
	for _, farmId := range request.FarmIds {
		joins = append(joins, &model.FarmOnFarmGroup{
			FarmId:      farmId,
			FarmGroupId: fg.Id,
		})
	}
	if err := s.farmGroupRepo.CreateJoins(ctx, joins); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}

	return nil
}

func (s *farmGroupService) List(clientId int) ([]*dto.FarmGroupResponse, error) {
	groups, err := s.farmGroupRepo.ListByClientId(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.FarmGroupResponse, 0, len(groups))
	for _, fg := range groups {
		farms, err := s.farmGroupRepo.GetFarmsByFarmGroupId(fg.Id)
		if err != nil {
			return nil, errors.ErrGeneric.Wrap(err)
		}
		responses = append(responses, s.toResponse(fg, farms))
	}
	return responses, nil
}

func (s *farmGroupService) GetClientIdByFarmGroupId(farmGroupId int) (int, error) {
	fg, err := s.farmGroupRepo.GetByID(farmGroupId)
	if err != nil {
		return 0, errors.ErrGeneric.Wrap(err)
	}
	if fg == nil {
		return 0, errors.ErrFarmGroupNotFound
	}
	return fg.ClientId, nil
}

func (s *farmGroupService) GetDropdown(clientId int) ([]*dto.DropdownItem, error) {
	groups, err := s.farmGroupRepo.ListByClientId(clientId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	items := make([]*dto.DropdownItem, 0, len(groups))
	for _, fg := range groups {
		items = append(items, &dto.DropdownItem{
			Key:   fg.Id,
			Value: fg.Name,
		})
	}
	return items, nil
}

func (s *farmGroupService) toResponse(fg *model.FarmGroup, farms []*model.Farm) *dto.FarmGroupResponse {
	farmItems := make([]dto.FarmGroupFarmItem, 0, len(farms))
	for _, f := range farms {
		farmItems = append(farmItems, dto.FarmGroupFarmItem{
			Id:   f.Id,
			Name: f.Name,
		})
	}
	return &dto.FarmGroupResponse{
		Id:        fg.Id,
		ClientId:  fg.ClientId,
		Name:      fg.Name,
		Farms:     farmItems,
		CreatedAt: fg.CreatedAt,
		CreatedBy: fg.CreatedBy,
		UpdatedAt: fg.UpdatedAt,
		UpdatedBy: fg.UpdatedBy,
	}
}
