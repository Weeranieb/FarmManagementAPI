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

type FeedCollectionHandlerTestSuite struct {
	suite.Suite
	feedCollectionService *mocks.MockFeedCollectionService
	feedCollectionHandler FeedCollectionHandler
}

func (s *FeedCollectionHandlerTestSuite) SetupTest() {
	s.feedCollectionService = mocks.NewMockFeedCollectionService(s.T())
	s.feedCollectionHandler = NewFeedCollectionHandler(s.feedCollectionService)
}

func (s *FeedCollectionHandlerTestSuite) TearDownTest() {
	s.feedCollectionService.ExpectedCalls = nil
}

func TestFeedCollectionHandlerSuite(t *testing.T) {
	suite.Run(t, new(FeedCollectionHandlerTestSuite))
}

func (s *FeedCollectionHandlerTestSuite) TestAddFeedCollection_Success() {
	createReq := &dto.CreateFeedCollectionRequest{
		Code: "FEED001",
		Name: "Test Feed",
		Unit: "kg",
	}

	expectedResponse := &dto.CreateFeedCollectionResponse{
		FeedCollection: &dto.FeedCollectionResponse{
			Id:       1,
			ClientId: 1,
			Code:     createReq.Code,
			Name:     createReq.Name,
			Unit:     createReq.Unit,
		},
		FeedPriceHistory: []interface{}{},
	}

	username := "admin"
	clientId := 1
	s.feedCollectionService.On("Create", *createReq, username, clientId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"username": username,
		"clientId": clientId,
	}))
	app.Post("/api/v1/feedcollection", s.feedCollectionHandler.AddFeedCollection)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/feedcollection", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.feedCollectionService.AssertExpectations(s.T())
}

func (s *FeedCollectionHandlerTestSuite) TestGetFeedCollection_Success() {
	feedCollectionId := 1
	expectedResponse := &dto.FeedCollectionResponse{
		Id:       feedCollectionId,
		ClientId: 1,
		Code:     "FEED001",
		Name:     "Test Feed",
		Unit:     "kg",
	}

	s.feedCollectionService.On("Get", feedCollectionId).Return(expectedResponse, nil)

	app := fiber.New()
	app.Get("/api/v1/feedcollection/:id", s.feedCollectionHandler.GetFeedCollection)

	req := httptest.NewRequest("GET", "/api/v1/feedcollection/1", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.feedCollectionService.AssertExpectations(s.T())
}

func (s *FeedCollectionHandlerTestSuite) TestListFeedCollection_Success() {
	clientId := 1
	page := 0
	pageSize := 10
	expectedResponse := &dto.PageResponse{
		Items: []*dto.FeedCollectionPageResponse{
			{FeedCollectionResponse: dto.FeedCollectionResponse{Id: 1, ClientId: clientId, Code: "FEED001", Name: "Feed 1", Unit: "kg"}},
		},
		Total: 1,
	}

	s.feedCollectionService.On("GetPage", clientId, page, pageSize, "", "").Return(expectedResponse, nil)

	app := fiber.New()
	app.Use(setLocalsMiddleware(map[string]interface{}{
		"clientId": clientId,
	}))
	app.Get("/api/v1/feedcollection", s.feedCollectionHandler.ListFeedCollection)

	req := httptest.NewRequest("GET", "/api/v1/feedcollection?page=0&pageSize=10", nil)

	resp, err := app.Test(req)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), fiber.StatusOK, resp.StatusCode)
	s.feedCollectionService.AssertExpectations(s.T())
}

