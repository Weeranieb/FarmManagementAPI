package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/weeranieb/boonmafarm-backend/src/internal/config"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
	"go.uber.org/dig"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyFeedHandler --output=./mocks --outpkg=handler --filename=daily_feed_handler.go --structname=MockDailyFeedHandler --with-expecter=false
type DailyFeedHandler interface {
	GetMonth(c *fiber.Ctx) error
	BulkUpsert(c *fiber.Ctx) error
	DeleteTable(c *fiber.Ctx) error
	UploadExcel(c *fiber.Ctx) error
}

type DailyFeedHandlerParams struct {
	dig.In

	Config           *config.Config
	DailyFeedService service.DailyFeedService
}

type dailyFeedHandlerImpl struct {
	dailyFeedService service.DailyFeedService
	cfg              *config.Config
}

func NewDailyFeedHandler(p DailyFeedHandlerParams) DailyFeedHandler {
	return &dailyFeedHandlerImpl{
		dailyFeedService: p.DailyFeedService,
		cfg:              p.Config,
	}
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

func (h *dailyFeedHandlerImpl) UploadExcel(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	month := c.FormValue("month")
	if month == "" {
		return http.Error(c, errors.ErrValidationFailed.Code, "month is required")
	}

	feedCollectionId, err := strconv.Atoi(c.FormValue("feedCollectionId"))
	if err != nil || feedCollectionId <= 0 {
		return http.Error(c, errors.ErrValidationFailed.Code, "feedCollectionId is required")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "file is required")
	}

	filename := filepath.Base(fileHeader.Filename)
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		return http.Error(c, errors.ErrValidationFailed.Code, "only .xlsx files are allowed")
	}

	uploadRoot := h.cfg.App.DailyFeedUploadPath
	if uploadRoot == "" {
		uploadRoot = "./data/uploads/daily-feed"
	}

	pondDir := filepath.Join(uploadRoot, fmt.Sprintf("pond_%d", pondId))
	if err := os.MkdirAll(pondDir, 0755); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	storedName := fmt.Sprintf("%s_%s", uuid.New().String(), filename)
	destPath := filepath.Join(pondDir, storedName)

	if err := c.SaveFile(fileHeader, destPath); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	relPath := filepath.ToSlash(filepath.Join(fmt.Sprintf("pond_%d", pondId), storedName))

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	rows, err := h.dailyFeedService.ImportFromExcelFile(
		c.UserContext(),
		pondId,
		feedCollectionId,
		month,
		destPath,
		username,
	)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, dto.DailyFeedExcelUploadResponse{
		RowsImported: rows,
		SavedPath:    relPath,
	})
}
