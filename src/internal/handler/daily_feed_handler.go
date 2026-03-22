package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyFeedHandler --output=./mocks --outpkg=handler --filename=daily_feed_handler.go --structname=MockDailyFeedHandler --with-expecter=false
type DailyFeedHandler interface {
	GetMonth(c *fiber.Ctx) error
	BulkUpsert(c *fiber.Ctx) error
	DeleteTable(c *fiber.Ctx) error
}

type dailyFeedHandlerImpl struct {
	dailyFeedService service.DailyFeedService
}

func NewDailyFeedHandler(dailyFeedService service.DailyFeedService) DailyFeedHandler {
	return &dailyFeedHandlerImpl{dailyFeedService: dailyFeedService}
}

func (h *dailyFeedHandlerImpl) GetMonth(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	month := c.Query("month")
	if month == "" {
		return http.Error(c, errors.ErrValidationFailed.Code, "month query parameter is required (YYYY-MM)")
	}

	result, err := h.dailyFeedService.GetMonth(c.UserContext(), pondId, month)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

func (h *dailyFeedHandlerImpl) BulkUpsert(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	var request dto.DailyFeedBulkUpsertRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.dailyFeedService.BulkUpsert(c.UserContext(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func (h *dailyFeedHandlerImpl) DeleteTable(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	feedCollectionId, err := strconv.Atoi(c.Params("feedCollectionId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed collection ID")
	}

	err = h.dailyFeedService.DeleteTable(c.UserContext(), pondId, feedCollectionId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}
