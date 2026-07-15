package handler

import (
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

type FeedPriceHistoryHandler interface {
	AddFeedPriceHistory(c *fiber.Ctx) error
	GetFeedPriceHistory(c *fiber.Ctx) error
	UpdateFeedPriceHistory(c *fiber.Ctx) error
	GetAllFeedPriceHistory(c *fiber.Ctx) error
}

type feedPriceHistoryHandlerImpl struct {
	feedPriceHistoryService service.FeedPriceHistoryService
	feedCollectionService   service.FeedCollectionService
}

func NewFeedPriceHistoryHandler(feedPriceHistoryService service.FeedPriceHistoryService, feedCollectionService service.FeedCollectionService) FeedPriceHistoryHandler {
	return &feedPriceHistoryHandlerImpl{
		feedPriceHistoryService: feedPriceHistoryService,
		feedCollectionService:   feedCollectionService,
	}
}

func (h *feedPriceHistoryHandlerImpl) AddFeedPriceHistory(c *fiber.Ctx) error {
	var createFeedPriceHistoryRequest dto.CreateFeedPriceHistoryRequest

	if err := validateAndParse(c, &createFeedPriceHistoryRequest); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.feedPriceHistoryService.Create(c.UserContext(), createFeedPriceHistoryRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *feedPriceHistoryHandlerImpl) GetFeedPriceHistory(c *fiber.Ctx) error {

	id, err := parseParamInt(c, "id", "Invalid feed price history ID")
	if err != nil {
		return err
	}

	result, err := h.feedPriceHistoryService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *feedPriceHistoryHandlerImpl) UpdateFeedPriceHistory(c *fiber.Ctx) error {
	var updateFeedPriceHistory dto.UpdateFeedPriceHistoryRequest

	if err := validateAndParse(c, &updateFeedPriceHistory); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.feedPriceHistoryService.Update(c.UserContext(), updateFeedPriceHistory, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *feedPriceHistoryHandlerImpl) GetAllFeedPriceHistory(c *fiber.Ctx) error {

	feedCollectionIdStr := c.Query("feedCollectionId")
	feedCollectionId, err := strconv.Atoi(feedCollectionIdStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed collection ID")
	}

	feedCollection, err := h.feedCollectionService.Get(feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	if err := requireClientAccess(c, feedCollection.ClientId); err != nil {
		return err
	}

	result, err := h.feedPriceHistoryService.GetAll(feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
