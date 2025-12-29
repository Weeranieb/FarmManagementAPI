package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type WorkerHandlerTestSuite struct {
	suite.Suite
	workerService *mocks.MockWorkerService
	workerHandler WorkerHandler
}

func (s *WorkerHandlerTestSuite) SetupTest() {
	s.workerService = mocks.NewMockWorkerService(s.T())
	s.workerHandler = NewWorkerHandler(s.workerService)
}

func (s *WorkerHandlerTestSuite) TearDownTest() {
	s.workerService.ExpectedCalls = nil
}

func TestWorkerHandlerSuite(t *testing.T) {
	suite.Run(t, new(WorkerHandlerTestSuite))
}

func (s *WorkerHandlerTestSuite) TestAddWorker_Success() {
	lastName := "Doe"
	createReq := &dto.CreateWorkerRequest{
		FarmGroupId:   1,
		FirstName:     "John",
		LastName:      &lastName,
		Nationality:   "Thai",
		Salary:        50000,
	}

	expectedResponse := &dto.WorkerResponse{
		Id:          1,
		ClientId:    1,
		FarmGroupId: createReq.FarmGroupId,
		FirstName:   createReq.FirstName,
		LastName:    createReq.LastName,
		Nationality: createReq.Nationality,
		Salary:      createReq.Salary,
		IsActive:    true,
	}

	username := "admin"
	clientId := 1
	s.workerService.On("Create", *createReq, username, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
		"clientId": clientId,
	}))
	app.Post("/api/v1/worker", s.workerHandler.AddWorker)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/worker", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.workerService.AssertExpectations(s.T())
}

func (s *WorkerHandlerTestSuite) TestGetWorker_Success() {
	workerId := 1
	expectedResponse := &dto.WorkerResponse{
		Id:           workerId,
		ClientId:     1,
		FarmGroupId:  1,
		FirstName:    "John",
		Nationality:  "Thai",
		Salary:       50000,
		IsActive:     true,
	}

	s.workerService.On("Get", workerId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/worker/:id", s.workerHandler.GetWorker)

	req := httptest.NewRequest("GET", "/api/v1/worker/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.workerService.AssertExpectations(s.T())
}

func (s *WorkerHandlerTestSuite) TestListWorker_Success() {
	clientId := 1
	page := 0
	pageSize := 10
	expectedResponse := &dto.PageResponse{
		Items: []*dto.WorkerResponse{
			{Id: 1, ClientId: clientId, FirstName: "John", Nationality: "Thai", Salary: 50000},
			{Id: 2, ClientId: clientId, FirstName: "Jane", Nationality: "Thai", Salary: 60000},
		},
		Total: 2,
	}

	s.workerService.On("GetPage", clientId, page, pageSize, "", "").Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/worker", s.workerHandler.ListWorker)

	req := httptest.NewRequest("GET", "/api/v1/worker?page=0&pageSize=10", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.workerService.AssertExpectations(s.T())
}

