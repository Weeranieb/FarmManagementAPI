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

type MerchantHandlerTestSuite struct {
	suite.Suite
	merchantService *mocks.MockMerchantService
	merchantHandler MerchantHandler
}

func (s *MerchantHandlerTestSuite) SetupTest() {
	s.merchantService = mocks.NewMockMerchantService(s.T())
	s.merchantHandler = NewMerchantHandler(s.merchantService)
}

func (s *MerchantHandlerTestSuite) TearDownTest() {
	s.merchantService.ExpectedCalls = nil
}

func TestMerchantHandlerSuite(t *testing.T) {
	suite.Run(t, new(MerchantHandlerTestSuite))
}

func (s *MerchantHandlerTestSuite) TestAddMerchant_Success() {
	createReq := &dto.CreateMerchantRequest{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}

	expectedResponse := &dto.MerchantResponse{
		Id:            1,
		Name:          createReq.Name,
		ContactNumber: createReq.ContactNumber,
		Location:      createReq.Location,
	}

	username := "admin"
	s.merchantService.On("Create", *createReq, username).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
	}))
	app.Post("/api/v1/merchant", s.merchantHandler.AddMerchant)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/merchant", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.merchantService.AssertExpectations(s.T())
}

func (s *MerchantHandlerTestSuite) TestGetMerchant_Success() {
	merchantId := 1
	expectedResponse := &dto.MerchantResponse{
		Id:            merchantId,
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}

	s.merchantService.On("Get", merchantId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/merchant/:id", s.merchantHandler.GetMerchant)

	req := httptest.NewRequest("GET", "/api/v1/merchant/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.merchantService.AssertExpectations(s.T())
}

func (s *MerchantHandlerTestSuite) TestGetMerchantList_Success() {
	expectedResponse := []*dto.MerchantResponse{
		{Id: 1, Name: "Merchant 1", ContactNumber: "111", Location: "Loc 1"},
		{Id: 2, Name: "Merchant 2", ContactNumber: "222", Location: "Loc 2"},
	}

	s.merchantService.On("GetList").Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/merchant", s.merchantHandler.GetMerchantList)

	req := httptest.NewRequest("GET", "/api/v1/merchant", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.merchantService.AssertExpectations(s.T())
}

