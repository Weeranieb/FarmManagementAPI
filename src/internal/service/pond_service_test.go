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

type PondServiceTestSuite struct {
	suite.Suite
	pondRepo   *mocks.MockPondRepository
	pondService PondService
}

func (s *PondServiceTestSuite) SetupTest() {
	s.pondRepo = mocks.NewMockPondRepository(s.T())
	s.pondService = NewPondService(s.pondRepo)
}

func (s *PondServiceTestSuite) TearDownTest() {
	s.pondRepo.ExpectedCalls = nil
}

func TestPondServiceSuite(t *testing.T) {
	suite.Run(t, new(PondServiceTestSuite))
}

func (s *PondServiceTestSuite) TestCreate_Success() {
	req := dto.CreatePondRequest{
		FarmId: 1,
		Code:   "POND001",
		Name:   "Test Pond",
	}
	username := "admin"

	s.pondRepo.On("GetByFarmIdAndCode", req.FarmId, req.Code).Return(nil, nil)

	expectedTime := time.Now()
	expectedPond := &model.Pond{
		Id:     1,
		FarmId: req.FarmId,
		Code:   req.Code,
		Name:   req.Name,
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	s.pondRepo.On("Create", mock.AnythingOfType("*model.Pond")).Return(nil).Run(func(args mock.Arguments) {
		pond := args.Get(0).(*model.Pond)
		pond.Id = expectedPond.Id
		pond.CreatedAt = expectedPond.CreatedAt
		pond.UpdatedAt = expectedPond.UpdatedAt
	})

	result, err := s.pondService.Create(req, username)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Code, result.Code)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestCreateBatch_Success() {
	reqs := []dto.CreatePondRequest{
		{FarmId: 1, Code: "POND001", Name: "Pond 1"},
		{FarmId: 1, Code: "POND002", Name: "Pond 2"},
	}
	username := "admin"

	s.pondRepo.On("GetByFarmIdAndCode", 1, "POND001").Return(nil, nil)
	s.pondRepo.On("GetByFarmIdAndCode", 1, "POND002").Return(nil, nil)

	s.pondRepo.On("CreateBatch", mock.AnythingOfType("[]*model.Pond")).Return(nil).Run(func(args mock.Arguments) {
		ponds := args.Get(0).([]*model.Pond)
		for i := range ponds {
			ponds[i].Id = i + 1
			ponds[i].CreatedAt = time.Now()
			ponds[i].UpdatedAt = time.Now()
		}
	})

	result, err := s.pondService.CreateBatch(reqs, username)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGet_Success() {
	pondId := 1
	expectedPond := &model.Pond{
		Id:     pondId,
		FarmId: 1,
		Code:   "POND001",
		Name:   "Test Pond",
	}

	s.pondRepo.On("GetByID", pondId).Return(expectedPond, nil)

	result, err := s.pondService.Get(pondId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), pondId, result.Id)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGetList_Success() {
	farmId := 1
	ponds := []*model.Pond{
		{Id: 1, FarmId: farmId, Code: "POND001", Name: "Pond 1"},
		{Id: 2, FarmId: farmId, Code: "POND002", Name: "Pond 2"},
	}

	s.pondRepo.On("ListByFarmId", farmId).Return(ponds, nil)

	result, err := s.pondService.GetList(farmId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.pondRepo.AssertExpectations(s.T())
}

