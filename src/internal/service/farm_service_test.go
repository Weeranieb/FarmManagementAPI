package service

import (
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
	farmService FarmService
}

func (s *FarmServiceTestSuite) SetupTest() {
	s.farmRepo = mocks.NewMockFarmRepository(s.T())
	s.farmService = NewFarmService(s.farmRepo)
}

func (s *FarmServiceTestSuite) TearDownTest() {
	s.farmRepo.ExpectedCalls = nil
}

func TestFarmServiceSuite(t *testing.T) {
	suite.Run(t, new(FarmServiceTestSuite))
}

func (s *FarmServiceTestSuite) TestCreate_Success() {
	req := dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	username := "admin"
	clientId := 1

	s.farmRepo.On("GetByNameAndClientId", req.Name, clientId).Return(nil, nil)

	expectedTime := time.Now()
	expectedFarm := &model.Farm{
		Id:       1,
		ClientId: clientId,
		Name:     req.Name,
		Status:   "active",
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	s.farmRepo.On("Create", mock.AnythingOfType("*model.Farm")).Return(nil).Run(func(args mock.Arguments) {
		farm := args.Get(0).(*model.Farm)
		farm.Id = expectedFarm.Id
		farm.CreatedAt = expectedFarm.CreatedAt
		farm.UpdatedAt = expectedFarm.UpdatedAt
	})

	result, err := s.farmService.Create(req, username, clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Name, result.Name)
	assert.Equal(s.T(), "active", result.Status)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestCreate_FarmExists() {
	req := dto.CreateFarmRequest{
		ClientId: 1,
		Name:     "Test Farm",
	}
	username := "admin"
	clientId := 1

	existingFarm := &model.Farm{
		Id:       1,
		Name:     req.Name,
		ClientId: clientId,
	}

	s.farmRepo.On("GetByNameAndClientId", req.Name, clientId).Return(existingFarm, nil)

	result, err := s.farmService.Create(req, username, clientId)

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

	s.farmRepo.On("GetByID", farmId).Return(expectedFarm, nil)

	result, err := s.farmService.Get(farmId, &clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), farmId, result.Id)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGet_NotFound() {
	farmId := 999
	clientId := 1

	s.farmRepo.On("GetByID", farmId).Return(nil, nil)

	result, err := s.farmService.Get(farmId, &clientId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestUpdate_Success() {
	username := "admin"
	farm := &model.Farm{
		Id:       1,
		ClientId: 1,
		Name:     "Updated Farm",
		Status:   "active",
	}

	s.farmRepo.On("Update", farm).Return(nil)

	err := s.farmService.Update(farm, username)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), username, farm.UpdatedBy)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGetList_Success() {
	clientId := 1
	farms := []*model.Farm{
		{Id: 1, ClientId: clientId, Name: "Farm 1", Status: "active"},
		{Id: 2, ClientId: clientId, Name: "Farm 2", Status: "active"},
	}
	counts := &model.FarmCountByClientId{Total: 2, ActiveCount: 2}

	s.farmRepo.On("ListByClientId", clientId).Return(farms, nil)
	s.farmRepo.On("CountByClientId", clientId).Return(counts, nil)

	result, err := s.farmService.GetList(clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Len(s.T(), result.Farms, 2)
	assert.Equal(s.T(), 2, result.Total)
	assert.Equal(s.T(), 2, result.TotalActive)
	s.farmRepo.AssertExpectations(s.T())
}
