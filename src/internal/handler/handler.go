package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"
	"go.uber.org/dig"
)

type Handler struct {
	UserHandler             UserHandler
	AuthHandler             AuthHandler
	ClientHandler           ClientHandler
	FarmHandler             FarmHandler
	FarmGroupHandler        FarmGroupHandler
	MerchantHandler         MerchantHandler
	PondHandler             PondHandler
	WorkerHandler           WorkerHandler
	FeedCollectionHandler   FeedCollectionHandler
	FeedPriceHistoryHandler FeedPriceHistoryHandler
	FishSizeGradeHandler    FishSizeGradeHandler
	DailyLogHandler         DailyLogHandler
	ActivityHandler         ActivityHandler
}

type HandlerParams struct {
	dig.In

	UserHandler             UserHandler
	AuthHandler             AuthHandler
	ClientHandler           ClientHandler
	FarmHandler             FarmHandler
	FarmGroupHandler        FarmGroupHandler
	MerchantHandler         MerchantHandler
	PondHandler             PondHandler
	WorkerHandler           WorkerHandler
	FeedCollectionHandler   FeedCollectionHandler
	FeedPriceHistoryHandler FeedPriceHistoryHandler
	FishSizeGradeHandler    FishSizeGradeHandler
	DailyLogHandler         DailyLogHandler
	ActivityHandler         ActivityHandler
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		UserHandler:             params.UserHandler,
		AuthHandler:             params.AuthHandler,
		ClientHandler:           params.ClientHandler,
		FarmHandler:             params.FarmHandler,
		FarmGroupHandler:        params.FarmGroupHandler,
		MerchantHandler:         params.MerchantHandler,
		PondHandler:             params.PondHandler,
		WorkerHandler:           params.WorkerHandler,
		FeedCollectionHandler:   params.FeedCollectionHandler,
		FeedPriceHistoryHandler: params.FeedPriceHistoryHandler,
		FishSizeGradeHandler:    params.FishSizeGradeHandler,
		DailyLogHandler:         params.DailyLogHandler,
		ActivityHandler:         params.ActivityHandler,
	}
}

// validateAndParse parses the request body and validates the struct.
//
// Errors are wrapped (not flattened) so http.NewError can surface field-level
// details: bad JSON → 400 with the parser error as `details`; struct validation
// failure → 422 with a `fields` array listing each offending field's tag and
// human-readable cause.
func validateAndParse(c *fiber.Ctx, target any) error {
	if err := c.BodyParser(target); err != nil {
		return http.NewError(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Wrap(err))
	}

	if err := utils.ValidateStruct(target); err != nil {
		return http.NewError(c, errors.ErrValidationFailed.Code, errors.ErrValidationFailed.Wrap(err))
	}

	return nil
}

func parseParamInt(c *fiber.Ctx, name, errMsg string) (int, error) {
	id, err := strconv.Atoi(c.Params(name))
	if err != nil {
		return 0, http.Error(c, errors.ErrValidationFailed.Code, errMsg)
	}
	return id, nil
}

func requireSuperAdmin(c *fiber.Ctx) error {
	ok, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil || !ok {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}
	return nil
}

// resolveWriteClientId returns the client id for a client-scoped create: the
// JWT client id when present, otherwise a super admin must pass clientId in the
// body (and must have access to it). Shared by the feed-collection and merchant
// create handlers.
func resolveWriteClientId(c *fiber.Ctx, bodyClientId *int) (int, error) {
	clientIdPtr := utils.GetClientId(c.UserContext())
	if clientIdPtr != nil {
		return *clientIdPtr, nil
	}

	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		return 0, http.NewError(c, errors.ErrGeneric.Code, err)
	}
	if !isSuperAdmin {
		return 0, http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}
	if bodyClientId == nil || *bodyClientId <= 0 {
		return 0, http.Error(c, errors.ErrValidationFailed.Code, "clientId is required when your account has no client in token (select a client in the header)")
	}

	if err := requireClientAccess(c, *bodyClientId); err != nil {
		return 0, err
	}
	return *bodyClientId, nil
}

func requireClientAdmin(c *fiber.Ctx) error {
	ok, err := utils.IsClientAdminOrAbove(c.UserContext())
	if err != nil || !ok {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}
	return nil
}

func requireClientAccess(c *fiber.Ctx, targetClientId int) error {
	ok, err := utils.CanAccessClient(c.UserContext(), targetClientId)
	if err != nil || !ok {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}
	return nil
}

// resolveListClientId returns an optional client filter for list/dropdown endpoints.
// Super admin: optional ?clientId= (0 = no filter). Others: JWT client id.
func resolveListClientId(c *fiber.Ctx) (int, error) {
	isSuperAdmin, err := utils.IsSuperAdmin(c.UserContext())
	if err != nil {
		return 0, http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}
	if isSuperAdmin {
		if clientIdStr := c.Query("clientId"); clientIdStr != "" {
			clientIdVal, err := strconv.Atoi(clientIdStr)
			if err != nil {
				return 0, http.Error(c, errors.ErrValidationFailed.Code, "Invalid clientId parameter")
			}
			return clientIdVal, nil
		}
		return 0, nil
	}
	clientIdPtr := utils.GetClientId(c.UserContext())
	if clientIdPtr == nil {
		return 0, http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}
	return *clientIdPtr, nil
}
