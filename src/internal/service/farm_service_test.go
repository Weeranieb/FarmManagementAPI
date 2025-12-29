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
	farmRepo   *mocks.MockFarmRepository
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
		Code: "FARM001",
		Name: "Test Farm",
	}
	username := "admin"
	clientId := 1

	s.farmRepo.On("GetByCodeAndClientId", req.Code, clientId).Return(nil, nil)

	expectedTime := time.Now()
	expectedFarm := &model.Farm{
		Id:       1,
		ClientId: clientId,
		Code:     req.Code,
		Name:     req.Name,
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
	assert.Equal(s.T(), req.Code, result.Code)
	assert.Equal(s.T(), req.Name, result.Name)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestCreate_FarmExists() {
	req := dto.CreateFarmRequest{
		Code: "FARM001",
		Name: "Test Farm",
	}
	username := "admin"
	clientId := 1

	existingFarm := &model.Farm{
		Id:       1,
		Code:     req.Code,
		ClientId: clientId,
	}

	s.farmRepo.On("GetByCodeAndClientId", req.Code, clientId).Return(existingFarm, nil)

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
		Code:     "FARM001",
		Name:     "Test Farm",
	}

	s.farmRepo.On("GetByID", farmId).Return(expectedFarm, nil)

	result, err := s.farmService.Get(farmId, clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), farmId, result.Id)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestGet_NotFound() {
	farmId := 999
	clientId := 1

	s.farmRepo.On("GetByID", farmId).Return(nil, nil)

	result, err := s.farmService.Get(farmId, clientId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.farmRepo.AssertExpectations(s.T())
}

func (s *FarmServiceTestSuite) TestUpdate_Success() {
	username := "admin"
	farm := &model.Farm{
		Id:       1,
		ClientId: 1,
		Code:     "FARM001",
		Name:     "Updated Farm",
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
		{Id: 1, ClientId: clientId, Code: "FARM001", Name: "Farm 1"},
		{Id: 2, ClientId: clientId, Code: "FARM002", Name: "Farm 2"},
	}

	s.farmRepo.On("ListByClientId", clientId).Return(farms, nil)

	result, err := s.farmService.GetList(clientId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.farmRepo.AssertExpectations(s.T())
}

