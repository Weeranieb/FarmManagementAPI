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
)

// Test validateAndParse function
type ValidateAndParseTestSuite struct {
	suite.Suite
	app *fiber.App
}

func (s *ValidateAndParseTestSuite) SetupTest() {
	s.app = fiber.New()
}

func TestValidateAndParseSuite(t *testing.T) {
	suite.Run(t, new(ValidateAndParseTestSuite))
}

func (s *ValidateAndParseTestSuite) TestValidateAndParse_Success() {
	var result dto.CreateClientRequest

	s.app.Post("/test", func(c *fiber.Ctx) error {
		return validateAndParse(c, &result)
	})

	reqBody := dto.CreateClientRequest{
		Name:          "Test Client",
		OwnerName:     "Test Owner",
		ContactNumber: "1234567890",
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	assert.Equal(s.T(), reqBody.Name, result.Name)
	assert.Equal(s.T(), reqBody.OwnerName, result.OwnerName)
	assert.Equal(s.T(), reqBody.ContactNumber, result.ContactNumber)
}

func (s *ValidateAndParseTestSuite) TestValidateAndParse_InvalidJSON() {
	var result dto.CreateClientRequest

	s.app.Post("/test", func(c *fiber.Ctx) error {
		return validateAndParse(c, &result)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	// Error is wrapped in ResponseModel
	errorResp, ok := response["error"].(map[string]interface{})
	assert.True(s.T(), ok, "Error should be present in response")
	assert.Equal(s.T(), "500011", errorResp["code"]) // Code is returned as string
}

func (s *ValidateAndParseTestSuite) TestValidateAndParse_ValidationFailed() {
	var result dto.CreateClientRequest

	s.app.Post("/test", func(c *fiber.Ctx) error {
		return validateAndParse(c, &result)
	})

	// Missing required fields
	reqBody := map[string]interface{}{
		"name": "", // Empty name should fail validation
	}

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	// Error is wrapped in ResponseModel
	errorResp, ok := response["error"].(map[string]interface{})
	assert.True(s.T(), ok, "Error should be present in response")
	assert.Equal(s.T(), "500010", errorResp["code"]) // Code is returned as string
}

func (s *ValidateAndParseTestSuite) TestValidateAndParse_EmptyBody() {
	var result dto.CreateClientRequest

	s.app.Post("/test", func(c *fiber.Ctx) error {
		return validateAndParse(c, &result)
	})

	req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)

	// Error is wrapped in ResponseModel
	// Empty body causes BodyParser to fail, returning invalid request body error
	errorResp, ok := response["error"].(map[string]interface{})
	assert.True(s.T(), ok, "Error should be present in response")
	assert.Equal(s.T(), "500011", errorResp["code"]) // Invalid request body error code
}

// Test NewHandler function
type NewHandlerTestSuite struct {
	suite.Suite
}

func TestNewHandlerSuite(t *testing.T) {
	suite.Run(t, new(NewHandlerTestSuite))
}

func (s *NewHandlerTestSuite) TestNewHandler_WithNilHandlers() {
	// Test that NewHandler can be called with nil handlers (though not recommended)
	handler := NewHandler(HandlerParams{})

	assert.NotNil(s.T(), handler)
	assert.Nil(s.T(), handler.UserHandler)
	assert.Nil(s.T(), handler.AuthHandler)
	assert.Nil(s.T(), handler.ClientHandler)
	assert.Nil(s.T(), handler.FarmHandler)
	assert.Nil(s.T(), handler.MerchantHandler)
	assert.Nil(s.T(), handler.PondHandler)
	assert.Nil(s.T(), handler.WorkerHandler)
	assert.Nil(s.T(), handler.FeedCollectionHandler)
	assert.Nil(s.T(), handler.FeedPriceHistoryHandler)
}
