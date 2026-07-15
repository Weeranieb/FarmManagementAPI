package handler

import (
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

type PondHandler interface {
	AddPonds(c *fiber.Ctx) error
	GetPond(c *fiber.Ctx) error
	GetPondList(c *fiber.Ctx) error
	GetPondActivities(c *fiber.Ctx) error
	UpdatePond(c *fiber.Ctx) error
	DeletePond(c *fiber.Ctx) error
	FillPond(c *fiber.Ctx) error
	MovePond(c *fiber.Ctx) error
	SellPond(c *fiber.Ctx) error
	FillPondPreview(c *fiber.Ctx) error
	MovePondPreview(c *fiber.Ctx) error
	SellPondPreview(c *fiber.Ctx) error
	FillPondCalc(c *fiber.Ctx) error
	MovePondCalc(c *fiber.Ctx) error
	SellPondCalc(c *fiber.Ctx) error
	DownloadTemplate(c *fiber.Ctx) error
	BulkImportFarmPond(c *fiber.Ctx) error
}

type pondHandlerImpl struct {
	pondService service.PondService
}

func NewPondHandler(pondService service.PondService) PondHandler {
	return &pondHandlerImpl{
		pondService: pondService,
	}
}

// Create multiple ponds for a farm (farmId, names array). New ponds are created with status maintenance.
// Allowed for super admin (any farm) or client admin (farms in their own client only).
// @Summary      Create multiple ponds
// @Description  Create multiple ponds for a farm. Request: farmId, array of names. New ponds have status maintenance. Requires client-admin role or above.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.CreatePondsRequest true "farmId, names[]"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond [post]
func (h *pondHandlerImpl) AddPonds(c *fiber.Ctx) error {
	var createPondsRequest dto.CreatePondsRequest

	if err := validateAndParse(c, &createPondsRequest); err != nil {
		return err
	}

	// Must be client admin or super admin. Per-client scoping is enforced
	// in pondService.CreatePonds (which loads the farm to find its clientId).
	if err := requireClientAdmin(c); err != nil {
		return err
	}

	if err := h.pondService.CreatePonds(c.UserContext(), createPondsRequest); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, nil)
}

// Get a pond by ID.
// @Summary      Get a pond by ID
// @Description  Retrieve a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        id path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{id} [get]
func (h *pondHandlerImpl) GetPond(c *fiber.Ctx) error {

	id, err := parseParamInt(c, "id", "Invalid pond ID")
	if err != nil {
		return err
	}

	pond, err := h.pondService.Get(c.UserContext(), id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, pond)
}

// GET /pond/:pondId/activities
// List the activity history (fill/move/sell) for a pond, newest first.
// @Summary      List pond activity history
// @Description  Return the chronological fill/move/sell activity timeline for a pond, ordered by activity_date desc.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        pondId path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/activities [get]
func (h *pondHandlerImpl) GetPondActivities(c *fiber.Ctx) error {

	id, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	activities, err := h.pondService.ListActivities(c.UserContext(), id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, activities)
}

// Get a list of ponds by farm ID.
// @Summary      Get a list of ponds by farm ID
// @Description  Retrieve a list of ponds belonging to a specific farm
// @Tags         pond
// @Accept       json
// @Produce      json
// @Param        farmId query int true "Farm ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond [get]
func (h *pondHandlerImpl) GetPondList(c *fiber.Ctx) error {

	farmIdStr := c.Query("farmId")
	farmId, err := strconv.Atoi(farmIdStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid farm ID")
	}

	pondList, err := h.pondService.GetList(c.UserContext(), farmId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, pondList)
}

// Update a pond.
// @Summary      Update a pond
// @Description  Update an existing pond. Id in path; body contains optional farmId, name, status.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id   path int true "Pond ID"
// @Param        body body dto.UpdatePondBody true "Updated pond data (farmId, name, status optional)"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{id} [put]
func (h *pondHandlerImpl) UpdatePond(c *fiber.Ctx) error {

	id, err := parseParamInt(c, "id", "Invalid pond ID")
	if err != nil {
		return err
	}

	var body dto.UpdatePondBody
	if err := validateAndParse(c, &body); err != nil {
		return err
	}

	req := dto.UpdatePondRequest{Id: id, FarmId: body.FarmId, Name: body.Name, Status: body.Status, Area: body.Area}
	err = h.pondService.Update(c.UserContext(), req)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// Delete a pond.
// @Summary      Delete a pond
// @Description  Delete a pond by its ID
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        id path int true "Pond ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{id} [delete]
func (h *pondHandlerImpl) DeletePond(c *fiber.Ctx) error {
	id, err := parseParamInt(c, "id", "Invalid pond ID")
	if err != nil {
		return err
	}

	err = h.pondService.Delete(c.UserContext(), id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// Add fish to a pond (fill). Creates an active_pond if the pond is in maintenance.
// @Summary      Fill pond with fish
// @Description  Record a fill activity for a pond. If the pond has no active cycle, creates one. Request: fishType, amount, activityDate; optional fishWeight, fishUnit, pricePerUnit.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Pond ID"
// @Param        body   body dto.PondFillRequest true "fishType, amount, activityDate"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/fill [post]
func (h *pondHandlerImpl) FillPond(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondFillRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	response, err := h.pondService.FillPond(c.UserContext(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// Move fish from this pond (source) to another. Path = source pondId; body includes toPondId.
// @Summary      Move fish to another pond
// @Description  Transfer fish from this pond to another. If destination is in maintenance, backend creates active_pond for it.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Source pond ID"
// @Param        body   body dto.PondMoveRequest true "toPondId, fishType, amount, activityDate"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/move [post]
func (h *pondHandlerImpl) MovePond(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondMoveRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	response, err := h.pondService.MovePond(c.UserContext(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// Record a sell transaction from a pond. Optionally close the active cycle.
// @Summary      Sell fish from pond
// @Description  Record a sell activity. If markToClose is true, close the active cycle and set pond to maintenance.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Pond ID"
// @Param        body   body dto.PondSellRequest true "activityDate, details[], merchantId, markToClose"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/sell [post]
func (h *pondHandlerImpl) SellPond(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondSellRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	response, err := h.pondService.SellPond(c.UserContext(), pondId, request, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// POST /pond/:pondId/fill/preview
// Preview add-stock summary (Review & Confirm). Does not persist.
// @Summary      Preview fill pond
// @Description  Validate and compute fill summary (costs, stock impact) without persisting. For Review & Confirm flow.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Pond ID"
// @Param        body   body dto.PondFillRequest true "fishType, amount, activityDate, pricePerUnit, additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/fill/preview [post]
func (h *pondHandlerImpl) FillPondPreview(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondFillRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	response, err := h.pondService.PreviewFillPond(c.UserContext(), pondId, request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// POST /pond/:pondId/move/preview
// Preview transfer summary (Review & Confirm). Does not persist.
// @Summary      Preview move pond
// @Description  Validate and compute transfer summary (costs, stock impact) without persisting. For Review & Confirm flow.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Source pond ID"
// @Param        body   body dto.PondMoveRequest true "toPondId, fishType, amount, activityDate, pricePerUnit, additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/move/preview [post]
func (h *pondHandlerImpl) MovePondPreview(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondMoveRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	response, err := h.pondService.PreviewMovePond(c.UserContext(), pondId, request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// POST /pond/:pondId/sell/preview
// Preview sell summary (Review & Confirm). Does not persist.
// @Summary      Preview sell pond
// @Description  Validate and compute sell summary (revenue, items, stock impact) without persisting. For Review & Confirm flow.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        pondId path int true "Pond ID"
// @Param        body   body dto.PondSellRequest true "activityDate, details[], merchantId, markToClose, additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/{pondId}/sell/preview [post]
func (h *pondHandlerImpl) SellPondPreview(c *fiber.Ctx) error {

	pondId, err := parseParamInt(c, "pondId", "Invalid pond ID")
	if err != nil {
		return err
	}

	var request dto.PondSellRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	response, err := h.pondService.PreviewSellPond(c.UserContext(), pondId, request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, response)
}

// POST /pond/fill/calc
// Live totals for the fill form (no persistence, no DB). Pure math.
// @Summary      Calculate fill totals
// @Description  Compute live cost/weight totals for the fill (add stock) form. Pure math, no validation against pond state. Used for live form updates as the user types.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.PondFillCalcRequest true "amount, fishWeight, pricePerUnit, additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Router       /pond/fill/calc [post]
func (h *pondHandlerImpl) FillPondCalc(c *fiber.Ctx) error {

	var request dto.PondFillCalcRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}
	return http.Success(c, h.pondService.CalcFillPond(c.UserContext(), request))
}

// POST /pond/move/calc
// Live totals for the move form (no persistence, no DB). Pure math.
// @Summary      Calculate move totals
// @Description  Compute live cost/weight totals for the move (transfer) form. Pure math, no validation against pond state.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.PondMoveCalcRequest true "amount, fishWeight, pricePerUnit, additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Router       /pond/move/calc [post]
func (h *pondHandlerImpl) MovePondCalc(c *fiber.Ctx) error {

	var request dto.PondMoveCalcRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}
	return http.Success(c, h.pondService.CalcMovePond(c.UserContext(), request))
}

// POST /pond/sell/calc
// Live totals for the sell form (no persistence, no DB). Pure math.
// @Summary      Calculate sell totals
// @Description  Compute live revenue/weight totals for the sell form. Pure math, no fish-size-grade validation.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        body body dto.PondSellCalcRequest true "details[], additionalCosts"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Router       /pond/sell/calc [post]
func (h *pondHandlerImpl) SellPondCalc(c *fiber.Ctx) error {

	var request dto.PondSellCalcRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}
	return http.Success(c, h.pondService.CalcSellPond(c.UserContext(), request))
}

// GET /pond/template
// Download the bulk-import Excel template (1 farm + many ponds with area).
// @Summary      Download bulk-import pond template
// @Description  Returns the .xlsx template users fill in to bulk-create ponds for a single farm.
// @Tags         pond
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Security     BearerAuth
// @Security     CookieAuth
// @Success      200  {file}    file
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/template [get]
func (h *pondHandlerImpl) DownloadTemplate(c *fiber.Ctx) error {

	const templatePath = "./src/assets/templates/pond_template.xlsx"
	const downloadName = "pond_template.xlsx"

	c.Set(fiber.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if err := c.Download(templatePath, downloadName); err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return nil
}

// POST /pond/bulk-import/:clientId
// Bulk-import farms and ponds for a client. Idempotent: missing farms and ponds are created,
// existing ponds get their area updated when provided. Nothing is ever deleted.
// Allowed for super admin (any client) or client admin (their own client only).
// @Summary      Bulk import farms and ponds
// @Description  Upsert farms and ponds from a parsed template. Missing farms are created. Missing ponds are created with status maintenance. Existing ponds get their area updated when an area is provided. No deletes. Requires client-admin role or above.
// @Tags         pond
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Security     CookieAuth
// @Param        clientId path int true "Client ID"
// @Param        body body dto.BulkImportFarmPondRequest true "farms[] with ponds[]"
// @Success      200  {object}  http.ResponseModel{data=dto.BulkImportFarmPondResponse}
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      403  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /pond/bulk-import/{clientId} [post]
func (h *pondHandlerImpl) BulkImportFarmPond(c *fiber.Ctx) error {

	clientId, err := parseParamInt(c, "clientId", "Invalid client ID")
	if err != nil {
		return err
	}

	// Must be client admin or super admin.
	if err := requireClientAdmin(c); err != nil {
		return err
	}

	// Super admin can target any client; client admin can only target their own.
	if err := requireClientAccess(c, clientId); err != nil {
		return err
	}

	var request dto.BulkImportFarmPondRequest
	if err := validateAndParse(c, &request); err != nil {
		return err
	}

	resp, err := h.pondService.BulkImportFarmPond(c.UserContext(), clientId, request)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}
	return http.Success(c, resp)
}
