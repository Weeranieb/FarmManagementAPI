package handler

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
	"go.uber.org/dig"
)

type ScanHandler interface {
	ScanDailyLog(c *fiber.Ctx) error
}

type ScanHandlerParams struct {
	dig.In

	ScanService service.ScanService
}

type scanHandlerImpl struct {
	scanService service.ScanService
}

func NewScanHandler(p ScanHandlerParams) ScanHandler {
	return &scanHandlerImpl{
		scanService: p.ScanService,
	}
}

// POST /pond/:pondId/daily-logs/scan
// @Summary      Scan handwritten paper and extract daily log data using AI
// @Description  Upload images of handwritten daily feeding records. AI will extract the data for review.
// @Tags         pond
// @Accept       multipart/form-data
// @Param        pondId path int true "Pond ID"
// @Param        month formData string true "YYYY-MM"
// @Param        images formData file true "Image files of paper records"
// @Success      200  {object}  http.ResponseModel{data=dto.ScanDailyLogResponse}
// @Router       /pond/{pondId}/daily-logs/scan [post]
func (h *scanHandlerImpl) ScanDailyLog(c *fiber.Ctx) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	pondId, err := strconv.Atoi(c.Params("pondId"))
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid pond ID")
	}

	month := c.FormValue("month")
	if month == "" {
		return http.Error(c, errors.ErrValidationFailed.Code, "month is required (YYYY-MM)")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "invalid multipart form")
	}

	files := form.File["images"]
	if len(files) == 0 {
		return http.Error(c, errors.ErrValidationFailed.Code, "at least one image is required")
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.scanService.ScanDailyLog(c.UserContext(), pondId, month, files, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}
