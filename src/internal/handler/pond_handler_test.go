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

type PondHandlerTestSuite struct {
	suite.Suite
	pondService *mocks.MockPondService
	pondHandler PondHandler
}

func (s *PondHandlerTestSuite) SetupTest() {
	s.pondService = mocks.NewMockPondService(s.T())
	s.pondHandler = NewPondHandler(s.pondService)
}

func (s *PondHandlerTestSuite) TearDownTest() {
	s.pondService.ExpectedCalls = nil
}

func TestPondHandlerSuite(t *testing.T) {
	suite.Run(t, new(PondHandlerTestSuite))
}

func (s *PondHandlerTestSuite) TestAddPond_Success() {
	createReq := &dto.CreatePondRequest{
		FarmId: 1,
		Name:   "Test Pond",
	}

	expectedResponse := &dto.PondResponse{
		Id:     1,
		FarmId: createReq.FarmId,
		Name:   createReq.Name,
		Status: "active",
	}

	username := "admin"
	s.pondService.On("Create", *createReq, username).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Post("/api/v1/pond", s.pondHandler.AddPond)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestAddPonds_Success() {
	createReqs := []dto.CreatePondRequest{
		{FarmId: 1, Name: "Pond 1"},
		{FarmId: 1, Name: "Pond 2"},
	}

	expectedResponse := []*dto.PondResponse{
		{Id: 1, FarmId: 1, Name: "Pond 1", Status: "active"},
		{Id: 2, FarmId: 1, Name: "Pond 2", Status: "active"},
	}

	username := "admin"
	s.pondService.On("CreateBatch", createReqs, username).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Post("/api/v1/pond/batch", s.pondHandler.AddPonds)

	body, _ := json.Marshal(createReqs)
	req := httptest.NewRequest("POST", "/api/v1/pond/batch", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestGetPond_Success() {
	pondId := 1
	expectedResponse := &dto.PondResponse{
		Id:     pondId,
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}

	s.pondService.On("Get", pondId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/pond/:id", s.pondHandler.GetPond)

	req := httptest.NewRequest("GET", "/api/v1/pond/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestGetPondList_Success() {
	farmId := 1
	expectedResponse := []*dto.PondResponse{
		{Id: 1, FarmId: farmId, Name: "Pond 1", Status: "active"},
		{Id: 2, FarmId: farmId, Name: "Pond 2", Status: "active"},
	}

	s.pondService.On("GetList", farmId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/pond", s.pondHandler.GetPondList)

	req := httptest.NewRequest("GET", "/api/v1/pond?farmId=1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

