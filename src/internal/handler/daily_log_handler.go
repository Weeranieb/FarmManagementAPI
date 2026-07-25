package handler

import (
	"io"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
	"go.uber.org/dig"
)

type DailyLogHandler interface {
	GetMonth(c fiber.Ctx) error
	BulkUpsert(c fiber.Ctx) error
	UploadTemplate(c fiber.Ctx) error
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
func (h *dailyLogHandlerImpl) GetMonth(c fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	month := c.Query("month")
	if month == "" {
		return http.Error(c, errors.ErrValidationFailed.Code, "month query parameter is required (YYYY-MM)")
	}

	result, err := h.dailyLogService.GetMonth(c.Context(), pondId, month)
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
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogMonthResponse}
// @Router       /pond/{pondId}/daily-logs [put]
func (h *dailyLogHandlerImpl) BulkUpsert(c fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.DailyLogBulkUpsertRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.dailyLogService.BulkUpsert(c.Context(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	result, err := h.dailyLogService.GetMonth(c.Context(), pondId, request.Month)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

// POST /farm/:farmId/daily-logs/import-template
// @Summary      Upload multi-pond Excel template and import daily logs
// @Tags         farm
// @Param        farmId path int true "Farm ID"
// @Param        selectedPondIds formData []int true "Pond IDs to import"
// @Param        file formData file true "xlsx file"
// @Success      200  {object}  http.ResponseModel{data=dto.DailyLogTemplateImportResponse}
// @Router       /farm/{farmId}/daily-logs/import-template [post]
func (h *dailyLogHandlerImpl) UploadTemplate(c fiber.Ctx) error {

	farmId, err := parseParamInt(c, "farmId", "Invalid farm ID")
	if err != nil {
		return err
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

	username, err := utils.GetUsername(c.Context())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.dailyLogService.ImportFromTemplate(c.Context(), farmId, selectedPondIds, fileBytes, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
