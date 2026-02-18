package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type FarmServiceTestSuite struct {
	suite.Suite
	farmRepo    *mocks.MockFarmRepository
	pondRepo    *mocks.MockPondRepository
	farmService FarmService
}

func (s *FarmServiceTestSuite) SetupTest() {
	s.farmRepo = mocks.NewMockFarmRepository(s.T())
	s.pondRepo = mocks.NewMockPondRepository(s.T())
	s.farmService = NewFarmService(s.farmRepo, s.pondRepo)
}

func (s *FarmServiceTestSuite) TearDownTest() {
	s.farmRepo.ExpectedCalls = nil
	s.pondRepo.ExpectedCalls = nil
}

func TestFarmServiceSuite(t *testing.T) {
	suite.Run(t, new(FarmServiceTestSuite))
}

func (s *FarmServiceTestSuite) TestCreate_Success() {
	req := dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	clientId := 1

	s.farmRepo.On("GetByNameAndClientId", req.Name, clientId).Return(nil, nil)

	expectedTime := time.Now()
	expectedFarm := &model.Farm{
		Id:       1,
		ClientId: clientId,
		Name:     req.Name,
		Status:   "maintenance",
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
		},
	}

	s.farmRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Farm")).Return(nil).Run(func(args mock.Arguments) {
		farm := args.Get(1).(*model.Farm)
		farm.Id = expectedFarm.Id
		farm.CreatedAt = expectedFarm.CreatedAt
		farm.UpdatedAt = expectedFarm.UpdatedAt
	})

	result, err := s.farmService.Create(context.Background(), req, clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Name, result.Name)
	assert.Equal(s.T(), "maintenance", result.Status)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestCreate_FarmExists() {
	req := dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	clientId := 1

	existingFarm := &model.Farm{
		Id:       1,
		Name:     req.Name,
		ClientId: clientId,
	}

	s.farmRepo.On("GetByNameAndClientId", req.Name, clientId).Return(existingFarm, nil)

	result, err := s.farmService.Create(context.Background(), req, clientId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGet_Success() {
	farmId := 1
	clientId := 1
	expectedFarm := &model.Farm{
		Id:       farmId,
		ClientId: clientId,
		Name:     "Test Farm",
		Status:   "active",
	}
	ponds := []*model.Pond{
		{Id: 1, FarmId: farmId, Name: "Pond A1", Status: "active"},
		{Id: 2, FarmId: farmId, Name: "Pond A2", Status: "active"},
	}

	s.farmRepo.On("GetByID", farmId).Return(expectedFarm, nil)
	s.pondRepo.On("ListByFarmId", farmId).Return(ponds, nil)

	result, err := s.farmService.Get(farmId, &clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), farmId, result.Id)
	assert.Equal(s.T(), "Test Farm", result.Name)
	assert.Equal(s.T(), 2, result.Summary.TotalPonds)
	assert.Equal(s.T(), 2, result.Summary.ActivePonds)
	assert.Len(s.T(), result.Ponds, 2)
	s.farmRepo.AssertExpectations(s.T())
	s.pondRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGet_NotFound() {
	farmId := 999
	clientId := 1

	s.farmRepo.On("GetByID", farmId).Return(nil, nil)

	result, err := s.farmService.Get(farmId, &clientId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.farmRepo.AssertExpectations(s.T())
	// ListByFarmId not called when farm not found
}

func (s *FarmServiceTestSuite) TestUpdate_Success() {
	updateReq := dto.UpdateFarmRequest{
		Id:   1,
		Name: "Updated Farm",
	}
	existingFarm := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Old Farm",
		Status:   "active",
	}
	expectedUpdateFarm := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Updated Farm",
		Status:   "active",
	}

	s.farmRepo.On("GetByID", 1).Return(existingFarm, nil)
	s.farmRepo.On("GetByNameAndClientId", "Updated Farm", 1).Return(nil, nil)
	s.farmRepo.On("Update", mock.Anything, expectedUpdateFarm).Return(nil)

	err := s.farmService.Update(context.Background(), updateReq)

	assert.NoError(s.T(), err)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGetList_Success() {
	clientId := 1
	list := []*model.FarmWithPonds{
		{
			Farm:  model.Farm{Id: 1, ClientId: clientId, Name: "Farm 1", Status: "active"},
			Ponds: []*model.Pond{{Id: 1}, {Id: 2}, {Id: 3}}, // 3 ponds
		},
		{
			Farm:  model.Farm{Id: 2, ClientId: clientId, Name: "Farm 2", Status: "active"},
			Ponds: []*model.Pond{},
		},
	}
	counts := &model.FarmCountByClientId{Total: 2, ActiveCount: 2}

	s.farmRepo.On("ListByClientIdWithPonds", clientId).Return(list, nil)
	s.farmRepo.On("CountByClientId", clientId).Return(counts, nil)

	result, err := s.farmService.GetList(clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Farms, 2)
	assert.Equal(s.T(), 2, result.Total)
	assert.Equal(s.T(), 2, result.TotalActive)
	assert.Equal(s.T(), 3, result.Farms[0].PondCount, "first farm has 3 ponds")
	assert.Equal(s.T(), 0, result.Farms[1].PondCount, "second farm has 0 ponds")
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGetHierarchy_Success() {
	clientId := 1
	list := []*model.FarmWithPonds{
		{
			Farm:  model.Farm{Id: 1, ClientId: clientId, Name: "River Farm", Status: "active"},
			Ponds: []*model.Pond{{Id: 1, FarmId: 1, Name: "Pond A1", Status: "active"}, {Id: 2, FarmId: 1, Name: "Pond A2", Status: "maintenance"}},
		},
		{
			Farm:  model.Farm{Id: 2, ClientId: clientId, Name: "Delta Farm", Status: "active"},
			Ponds: []*model.Pond{},
		},
	}

	s.farmRepo.On("ListByClientIdWithPonds", clientId).Return(list, nil)

	result, err := s.farmService.GetHierarchy(clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 2)
	assert.Equal(s.T(), 1, result[0].Id)
	assert.Equal(s.T(), "River Farm", result[0].Name)
	assert.Len(s.T(), result[0].Ponds, 2)
	assert.Equal(s.T(), "Pond A1", result[0].Ponds[0].Name)
	assert.Equal(s.T(), "maintenance", result[0].Ponds[1].Status)
	assert.Equal(s.T(), 2, result[1].Id)
	assert.Equal(s.T(), "Delta Farm", result[1].Name)
	assert.Len(s.T(), result[1].Ponds, 0)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGetHierarchy_Empty() {
	clientId := 1
	s.farmRepo.On("ListByClientIdWithPonds", clientId).Return([]*model.FarmWithPonds{}, nil)

	result, err := s.farmService.GetHierarchy(clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result, 0)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGetHierarchy_RepoError() {
	clientId := 1
	s.farmRepo.On("ListByClientIdWithPonds", clientId).Return(([]*model.FarmWithPonds)(nil), errors.New("db error"))

	result, err := s.farmService.GetHierarchy(clientId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.farmRepo.AssertExpectations(s.T())
}
