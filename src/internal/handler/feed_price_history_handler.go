package handler

import (
	"fmt"
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
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
}

func NewFeedPriceHistoryHandler(feedPriceHistoryService service.FeedPriceHistoryService) FeedPriceHistoryHandler {
	return &feedPriceHistoryHandlerImpl{
		feedPriceHistoryService: feedPriceHistoryService,
	}
}

func (h *feedPriceHistoryHandlerImpl) AddFeedPriceHistory(c *fiber.Ctx) error {
	var createFeedPriceHistoryRequest dto.CreateFeedPriceHistoryRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createFeedPriceHistoryRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.feedPriceHistoryService.Create(createFeedPriceHistoryRequest, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *feedPriceHistoryHandlerImpl) GetFeedPriceHistory(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
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
	var updateFeedPriceHistory *model.FeedPriceHistory

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updateFeedPriceHistory); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	username, err := utils.GetUsername(c)
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.feedPriceHistoryService.Update(updateFeedPriceHistory, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *feedPriceHistoryHandlerImpl) GetAllFeedPriceHistory(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	feedCollectionIdStr := c.Query("feedCollectionId")
	feedCollectionId, err := strconv.Atoi(feedCollectionIdStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed collection ID")
	}

	result, err := h.feedPriceHistoryService.GetAll(feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
