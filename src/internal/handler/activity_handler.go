package handler

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
)

type ActivityHandler interface {
	GetActivityFeed(c fiber.Ctx) error
	GetActivitySellDetails(c fiber.Ctx) error
}

type activityHandlerImpl struct {
	activityService service.ActivityService
}

func NewActivityHandler(activityService service.ActivityService) ActivityHandler {
	return &activityHandlerImpl{
		activityService: activityService,
	}
}

// GET /activity
// Farm-wide activity feed — newest fill/move/sell events across every pond
// the caller's client owns. Daily-log saves are excluded by design.
// @Summary      Get the activity feed
// @Description  Newest discrete events (fill/move/sell) across all of the client's ponds, ordered by activity date desc. Daily-log saves are not included. Optional ?limit caps the row count (omit or 0 for the full history). To page, pass beforeDate + beforeId from the last row of the previous page; the response is a plain array, so a page shorter than ?limit means the end.
// @Tags         activity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        limit query int false "Max rows to return; 0 or omitted = all"
// @Param        beforeDate query string false "Cursor: activityDate of the last row already seen (YYYY-MM-DD or RFC3339). Requires beforeId."
// @Param        beforeId query int false "Cursor: id of the last row already seen. Requires beforeDate."
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      422  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /activity [get]
func (h *activityHandlerImpl) GetActivityFeed(c fiber.Ctx) error {

	limit := 0
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			return http.Error(c, errors.ErrValidationFailed.Code, "Invalid limit")
		}
		limit = parsed
	}

	before, err := parseFeedCursor(c)
	if err != nil {
		return err
	}

	feed, err := h.activityService.ListFeed(c.Context(), limit, before)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, feed)
}

// parseFeedCursor reads the (beforeDate, beforeId) pair. Both or neither: half a
// cursor is a client bug that would silently return page 1 again, so say so
// rather than paging in a loop over the same rows.
func parseFeedCursor(c fiber.Ctx) (*repository.ActivityFeedCursor, error) {
	rawDate := c.Query("beforeDate")
	rawId := c.Query("beforeId")
	if rawDate == "" && rawId == "" {
		return nil, nil
	}
	if rawDate == "" || rawId == "" {
		return nil, http.Error(c, errors.ErrValidationFailed.Code, "beforeDate and beforeId must be sent together")
	}

	id, err := strconv.Atoi(rawId)
	if err != nil || id <= 0 {
		return nil, http.Error(c, errors.ErrValidationFailed.Code, "Invalid beforeId")
	}

	// activityDate is date-only in the feed payload, but accept a full timestamp
	// too so a client can echo back whatever it was given.
	date, err := time.Parse(time.RFC3339, rawDate)
	if err != nil {
		date, err = time.Parse(time.DateOnly, rawDate)
		if err != nil {
			return nil, http.Error(c, errors.ErrValidationFailed.Code, "Invalid beforeDate")
		}
	}

	return &repository.ActivityFeedCursor{ActivityDate: date, Id: id}, nil
}

// GET /activity/:activityId/sell-details
// The per-size-grade breakdown of one sale. A sale is priced per grade, so the
// feed's summed total/weight cannot answer "how many baht for each size" — this
// is the only endpoint that can.
// @Summary      Get a sale's size-grade breakdown
// @Description  Size-grade lines for one sell activity (weight, ฿/kg, head count and line total per grade), smallest grade first. Returns an empty list when the activity is not a sell, does not exist, or belongs to another client.
// @Tags         activity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        activityId path int true "Sell activity id"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /activity/{activityId}/sell-details [get]
func (h *activityHandlerImpl) GetActivitySellDetails(c fiber.Ctx) error {

	id, err := parseParamInt(c, "activityId", "Invalid activity ID")
	if err != nil {
		return err
	}

	lines, err := h.activityService.ListSellDetails(c.Context(), id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, lines)
}
