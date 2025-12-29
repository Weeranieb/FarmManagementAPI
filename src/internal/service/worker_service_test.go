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

type WorkerServiceTestSuite struct {
	suite.Suite
	workerRepo   *mocks.MockWorkerRepository
	workerService WorkerService
}

func (s *WorkerServiceTestSuite) SetupTest() {
	s.workerRepo = mocks.NewMockWorkerRepository(s.T())
	s.workerService = NewWorkerService(s.workerRepo)
}

func (s *WorkerServiceTestSuite) TearDownTest() {
	s.workerRepo.ExpectedCalls = nil
}

func TestWorkerServiceSuite(t *testing.T) {
	suite.Run(t, new(WorkerServiceTestSuite))
}

func (s *WorkerServiceTestSuite) TestCreate_Success() {
	lastName := "Doe"
	req := dto.CreateWorkerRequest{
		FarmGroupId:   1,
		FirstName:     "John",
		LastName:      &lastName,
		Nationality:   "Thai",
		Salary:        50000,
	}
	username := "admin"
	clientId := 1

	s.workerRepo.On("GetByFarmGroupId", req.FarmGroupId).Return(nil, nil)

	expectedTime := time.Now()
	expectedWorker := &model.Worker{
		Id:            1,
		ClientId:      clientId,
		FarmGroupId:   req.FarmGroupId,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Nationality:   req.Nationality,
		Salary:        req.Salary,
		IsActive:      true,
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	s.workerRepo.On("Create", mock.AnythingOfType("*model.Worker")).Return(nil).Run(func(args mock.Arguments) {
		worker := args.Get(0).(*model.Worker)
		worker.Id = expectedWorker.Id
		worker.CreatedAt = expectedWorker.CreatedAt
		worker.UpdatedAt = expectedWorker.UpdatedAt
	})

	result, err := s.workerService.Create(req, username, clientId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.FirstName, result.FirstName)
	assert.True(s.T(), result.IsActive)
	s.workerRepo.AssertExpectations(s.T())
}

func (s *WorkerServiceTestSuite) TestGetPage_Success() {
	clientId := 1
	page := 0
	pageSize := 10
	workers := []*model.Worker{
		{Id: 1, ClientId: clientId, FirstName: "John", Nationality: "Thai", Salary: 50000},
		{Id: 2, ClientId: clientId, FirstName: "Jane", Nationality: "Thai", Salary: 60000},
	}

	s.workerRepo.On("GetPage", clientId, page, pageSize, "", "").Return(workers, int64(2), nil)

	result, err := s.workerService.GetPage(clientId, page, pageSize, "", "")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), int64(2), result.Total)
	assert.Len(s.T(), result.Items.([]*dto.WorkerResponse), 2)
	s.workerRepo.AssertExpectations(s.T())
}

