package handler

import (
	"fmt"
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedPriceHistoryHandler --output=./mocks --outpkg=handler --filename=feed_price_history_handler.go --structname=MockFeedPriceHistoryHandler --with-expecter=false
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

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed price history ID")
	}

	result, err := h.feedPriceHistoryService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *feedPriceHistoryHandlerImpl) UpdateFeedPriceHistory(c *fiber.Ctx) error {
	var updateFeedPriceHistory dto.UpdateFeedPriceHistoryRequest

	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			_ = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	feedCollectionIdStr := c.Query("feedCollectionId")
	feedCollectionId, err := strconv.Atoi(feedCollectionIdStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed collection ID")
	}

	feedCollection, err := h.feedCollectionService.Get(feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	canAccess, accessErr := utils.CanAccessClient(c.UserContext(), feedCollection.ClientId)
	if accessErr != nil || !canAccess {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	result, err := h.feedPriceHistoryService.GetAll(feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
