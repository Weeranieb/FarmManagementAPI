package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
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

func (s *PondServiceTestSuite) TestCreatePonds_Success() {
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}
	username := "admin"

	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(nil, nil)

	s.pondRepo.On("CreateBatch", mock.AnythingOfType("[]*model.Pond")).Return(nil).Run(func(args mock.Arguments) {
		ponds := args.Get(0).([]*model.Pond)
		for i := range ponds {
			ponds[i].Id = i + 1
			ponds[i].CreatedAt = time.Now()
			ponds[i].UpdatedAt = time.Now()
		}
	})

	err := s.pondService.CreatePonds(req, username)

	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestCreatePonds_PondAlreadyExists() {
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}
	username := "admin"

	// First name is free, second name already exists for this farm
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	existingPond := &model.Pond{Id: 99, FarmId: 1, Name: "Pond 2", Status: "active"}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(existingPond, nil)

	err := s.pondService.CreatePonds(req, username)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
	// CreateBatch must not be called
	s.pondRepo.AssertNotCalled(s.T(), "CreateBatch")
}

func (s *PondServiceTestSuite) TestGet_Success() {
	pondId := 1
	expectedPond := &model.Pond{
		Id:     pondId,
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
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
		{Id: 1, FarmId: farmId, Name: "Pond 1", Status: "active"},
		{Id: 2, FarmId: farmId, Name: "Pond 2", Status: "active"},
	}

	s.pondRepo.On("ListByFarmId", farmId).Return(ponds, nil)

	result, err := s.pondService.GetList(farmId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.pondRepo.AssertExpectations(s.T())
}

