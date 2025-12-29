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

type FarmHandlerTestSuite struct {
	suite.Suite
	farmService *mocks.MockFarmService
	farmHandler FarmHandler
}

func (s *FarmHandlerTestSuite) SetupTest() {
	s.farmService = mocks.NewMockFarmService(s.T())
	s.farmHandler = NewFarmHandler(s.farmService)
}

func (s *FarmHandlerTestSuite) TearDownTest() {
	s.farmService.ExpectedCalls = nil
}

func TestFarmHandlerSuite(t *testing.T) {
	suite.Run(t, new(FarmHandlerTestSuite))
}

func (s *FarmHandlerTestSuite) TestAddFarm_Success() {
	createReq := &dto.CreateFarmRequest{
		Code: "FARM001",
		Name: "Test Farm",
	}

	expectedResponse := &dto.FarmResponse{
		Id:       1,
		ClientId: 1,
		Code:     createReq.Code,
		Name:     createReq.Name,
	}

	username := "admin"
	clientId := 1
	s.farmService.On("Create", *createReq, username, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
		"clientId": clientId,
	}))
	app.Post("/api/v1/farm", s.farmHandler.AddFarm)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/farm", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarm_Success() {
	farmId := 1
	clientId := 1
	expectedResponse := &dto.FarmResponse{
		Id:       farmId,
		ClientId: clientId,
		Code:     "FARM001",
		Name:     "Test Farm",
	}

	s.farmService.On("Get", farmId, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/farm/:id", s.farmHandler.GetFarm)

	req := httptest.NewRequest("GET", "/api/v1/farm/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

func (s *FarmHandlerTestSuite) TestGetFarmList_Success() {
	clientId := 1
	expectedResponse := []*dto.FarmResponse{
		{Id: 1, ClientId: clientId, Code: "FARM001", Name: "Farm 1"},
		{Id: 2, ClientId: clientId, Code: "FARM002", Name: "Farm 2"},
	}

	s.farmService.On("GetList", clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/farm", s.farmHandler.GetFarmList)

	req := httptest.NewRequest("GET", "/api/v1/farm", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.farmService.AssertExpectations(s.T())
}

