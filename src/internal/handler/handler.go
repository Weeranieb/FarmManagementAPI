package handler

import (
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

// validateClientAccess checks if the user can access the target clientId
// Super admin (clientId == nil) can access any clientId
// Regular users can only access their own clientId
func validateClientAccess(c *fiber.Ctx, targetClientId int) error {
	clientId, canAccess := utils.GetClientIdForAccess(c.UserContext())
	if !canAccess {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, "client id not found")
	}

	// Super admin can access any clientId
	if clientId == nil {
		return nil
	}

	// Regular users can only access their own clientId
	if *clientId != targetClientId {
		return http.Error(c, errors.ErrAuthPermissionDenied.Code, errors.ErrAuthPermissionDenied.Message)
	}

	return nil
}
