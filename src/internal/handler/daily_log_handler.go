package handler

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
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
	UploadTemplate(c *fiber.Ctx) error
}

type DailyLogHandlerParams struct {
	dig.In

	DailyLogService service.DailyLogService
}

type dailyLogHandlerImpl struct {
	dailyLogService service.DailyLogService
}

func NewDailyLogHandler(p DailyLogHandlerParams) DailyLogHandler {
	return &dailyLogHandlerImpl{
		dailyLogService: p.DailyLogService,
	}
}

// GET /pond/:pondId/daily-logs
// @Summary      Daily logs for a pond month
// @Description  Returns one sheet per month (fresh + pellet columns, deaths, tourist catch).
// @Tags         pond
// @Param        pondId path int true "Pond ID"
// @Param        month query string true "YYYY-MM"
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogMonthResponse}
// @Router       /pond/{pondId}/daily-logs [get]
func (h *dailyLogHandlerImpl) GetMonth(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
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
func (h *dailyLogHandlerImpl) BulkUpsert(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
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

// POST /farm/:farmId/daily-logs/import-template
// @Summary      Upload multi-pond Excel template and import daily logs
// @Tags         farm
// @Param        farmId path int true "Farm ID"
// @Param        selectedPondIds formData []int true "Pond IDs to import"
// @Param        file formData file true "xlsx file"
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogTemplateImportResponse}
// @Router       /farm/{farmId}/daily-logs/import-template [post]
func (h *dailyLogHandlerImpl) UploadTemplate(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	farmId, err := strconv.Atoi(c.Params("farmId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid farm ID")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "invalid multipart form")
	}
	defer func() { _ = form.RemoveAll() }()

	selectedPondIds, err := utils.ConvertRepeatedFormInts("selectedPondIds", form.Value["selectedPondIds"])
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, err.Error())
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "file is required")
	}
	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".xlsx") {
		return http.Error(c, errors.ErrValidationFailed.Code, "only .xlsx files are allowed")
	}

	f, err := fileHeader.Open()
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
	}
	defer func() { _ = f.Close() }()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, errors.ErrGeneric.Wrap(err))
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.dailyLogService.ImportFromTemplate(c.UserContext(), farmId, selectedPondIds, fileBytes, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
