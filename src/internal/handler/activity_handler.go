package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
)

type ActivityHandler interface {
	GetActivityFeed(c *fiber.Ctx) error
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
// @Description  Newest discrete events (fill/move/sell) across all of the client's ponds, ordered by activity date desc. Daily-log saves are not included. Optional ?limit caps the row count (omit or 0 for the full history).
// @Tags         activity
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        limit query int false "Max rows to return; 0 or omitted = all"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /activity [get]
func (h *activityHandlerImpl) GetActivityFeed(c *fiber.Ctx) error {

	limit := 0
	if raw := c.Query("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil || parsed < 0 {
			return http.Error(c, errors.ErrValidationFailed.Code, "Invalid limit")
		}
		limit = parsed
	}

	feed, err := h.activityService.ListFeed(c.UserContext(), limit)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, feed)
}
