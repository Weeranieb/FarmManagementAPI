package handler

import (
	"fmt"
	"log"
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

//go:generate go run github.com/vektra/mockery/v2@latest --name=DailyLogHandler --output=./mocks --outpkg=handler --filename=daily_log_handler.go --structname=MockDailyLogHandler --with-expecter=false
type DailyLogHandler interface {
	GetMonth(c *fiber.Ctx) error
	BulkUpsert(c *fiber.Ctx) error
	UploadExcel(c *fiber.Ctx) error
}

type DailyLogHandlerParams struct {
	dig.In

	Config          *config.Config
	DailyLogService service.DailyLogService
}

type dailyLogHandlerImpl struct {
	dailyLogService service.DailyLogService
	cfg             *config.Config
}

func NewDailyLogHandler(p DailyLogHandlerParams) DailyLogHandler {
	return &dailyLogHandlerImpl{
		dailyLogService: p.DailyLogService,
		cfg:             p.Config,
	}
}

func (h *dailyLogHandlerImpl) dailyLogUploadPath() string {
	p := h.cfg.App.DailyLogUploadPath
	if p != "" {
		return p
	}
	if h.cfg.App.DailyFeedUploadPath != "" {
		return h.cfg.App.DailyFeedUploadPath
	}
	return "./data/uploads/daily-log"
}

// GET /pond/:pondId/daily-logs
// @Summary      Daily logs for a pond month
// @Description  Returns one sheet per month (fresh + pellet columns, deaths, tourist catch).
// @Tags         pond
// @Param        pondId path int true "Pond ID"
// @Param        month query string true "YYYY-MM"
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogMonthResponse}
// @Router       /pond/{pondId}/daily-logs [get]
func (h *dailyLogHandlerImpl) GetMonth(c *fiber.Ctx) error {
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

	result, err := h.dailyLogService.GetMonth(c.UserContext(), pondId, month)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

// PUT /pond/:pondId/daily-logs
// @Summary      Upsert daily logs for a month
// @Tags         pond
// @Param        pondId path int true "Pond ID"
// @Param        body body dto.DailyLogBulkUpsertRequest true "Month + optional collection IDs + entries"
// @Success      200  {object}  http.ResponseModel
// @Router       /pond/{pondId}/daily-logs [put]
func (h *dailyLogHandlerImpl) BulkUpsert(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	var request dto.DailyLogBulkUpsertRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.dailyLogService.BulkUpsert(c.UserContext(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

func parseOptionalPositiveIntPtr(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil || v <= 0 {
		return nil
	}
	return &v
}

// POST /pond/:pondId/daily-logs/upload
// @Summary      Upload Excel and import daily logs
// @Tags         pond
// @Param        pondId path int true "Pond ID"
// @Param        month formData string true "YYYY-MM"
// @Param        freshFeedCollectionId formData int false "Fresh feed collection id"
// @Param        pelletFeedCollectionId formData int false "Pellet feed collection id"
// @Param        file formData file true "xlsx file"
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogExcelUploadResponse}
// @Router       /pond/{pondId}/daily-logs/upload [post]
func (h *dailyLogHandlerImpl) UploadExcel(c *fiber.Ctx) error {
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

	freshID := parseOptionalPositiveIntPtr(c.FormValue("freshFeedCollectionId"))
	pelletID := parseOptionalPositiveIntPtr(c.FormValue("pelletFeedCollectionId"))

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "file is required")
	}

	filename := filepath.Base(fileHeader.Filename)
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		return http.Error(c, errors.ErrValidationFailed.Code, "only .xlsx files are allowed")
	}

	uploadRoot := h.dailyLogUploadPath()
	pondDir := filepath.Join(uploadRoot, fmt.Sprintf("pond_%d", pondId))
	if err := os.MkdirAll(pondDir, 0755); err != nil {
		log.Printf("daily log upload: mkdir %q: %v", pondDir, err)
		return http.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
	}

	storedName := fmt.Sprintf("%s_%s", uuid.New().String(), filename)
	destPath := filepath.Join(pondDir, storedName)

	if err := c.SaveFile(fileHeader, destPath); err != nil {
		log.Printf("daily log upload: save file %q: %v", destPath, err)
		_ = os.Remove(destPath)
		return http.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
	}

	relPath := filepath.ToSlash(filepath.Join(fmt.Sprintf("pond_%d", pondId), storedName))

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	rows, err := h.dailyLogService.ImportFromExcelFile(
		c.UserContext(),
		pondId,
		freshID,
		pelletID,
		month,
		destPath,
		username,
	)
	if err != nil {
		if remErr := os.Remove(destPath); remErr != nil {
			log.Printf("daily log upload: cleanup %q after import error: %v", destPath, remErr)
		}
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, dto.DailyLogExcelUploadResponse{
		RowsImported: rows,
		SavedPath:    relPath,
	})
}
