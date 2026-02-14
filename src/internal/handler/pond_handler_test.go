package handler

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (s *PondHandlerTestSuite) TestAddPonds_Success() {
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}

	username := "admin"
	s.pondService.On("CreatePonds", mock.Anything, *createReq, username).Return(nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username":  username,
		"userLevel": 3, // super admin only
	}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.pondService.AssertExpectations(s.T())
}

func (s *PondHandlerTestSuite) TestAddPonds_NonSuperAdmin_ReturnsPermissionDenied() {
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1"},
	}

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username":  "user",
		"userLevel": 1, // normal user, not super admin
	}))
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
	if errObj, ok := result["error"].(map[string]any); ok && errObj["code"] != nil {
		assert.Equal(s.T(), "500024", errObj["code"]) // ErrAuthPermissionDenied
	}
}

func (s *PondHandlerTestSuite) TestAddPonds_IsSuperAdminError() {
	createReq := &dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1"},
	}

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]any{})) // no userLevel
	app.Post("/api/v1/pond", s.pondHandler.AddPonds)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/pond", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	var result map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&result)
	assert.NotNil(s.T(), result["error"])
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

