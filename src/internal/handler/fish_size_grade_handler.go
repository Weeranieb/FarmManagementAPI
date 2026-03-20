package handler

import (
	"fmt"

	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FishSizeGradeHandler --output=./mocks --outpkg=handler --filename=fish_size_grade_handler.go --structname=MockFishSizeGradeHandler --with-expecter=false
type FishSizeGradeHandler interface {
	GetDropdown(c *fiber.Ctx) error
}

type fishSizeGradeHandlerImpl struct {
	fishSizeGradeService service.FishSizeGradeService
}

func NewFishSizeGradeHandler(fishSizeGradeService service.FishSizeGradeService) FishSizeGradeHandler {
	return &fishSizeGradeHandlerImpl{
		fishSizeGradeService: fishSizeGradeService,
	}
}

// GET /fish-size-grade/dropdown
// @Summary      Get fish size grade dropdown
// @Description  Returns dropdown items for fish size grades sorted by sort_index
// @Tags         fish-size-grade
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {object}  http.ResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /fish-size-grade/dropdown [get]
func (h *fishSizeGradeHandlerImpl) GetDropdown(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	items, err := h.fishSizeGradeService.GetDropdown()
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, items)
}
