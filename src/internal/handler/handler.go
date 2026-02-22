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
	MerchantHandler         MerchantHandler
	PondHandler             PondHandler
	WorkerHandler           WorkerHandler
	FeedCollectionHandler   FeedCollectionHandler
	FeedPriceHistoryHandler FeedPriceHistoryHandler
}

type HandlerParams struct {
	dig.In

	UserHandler             UserHandler
	AuthHandler             AuthHandler
	ClientHandler           ClientHandler
	FarmHandler             FarmHandler
	MerchantHandler         MerchantHandler
	PondHandler             PondHandler
	WorkerHandler           WorkerHandler
	FeedCollectionHandler   FeedCollectionHandler
	FeedPriceHistoryHandler FeedPriceHistoryHandler
}

func NewHandler(params HandlerParams) *Handler {
	return &Handler{
		UserHandler:             params.UserHandler,
		AuthHandler:             params.AuthHandler,
		ClientHandler:           params.ClientHandler,
		FarmHandler:             params.FarmHandler,
		MerchantHandler:         params.MerchantHandler,
		PondHandler:             params.PondHandler,
		WorkerHandler:           params.WorkerHandler,
		FeedCollectionHandler:   params.FeedCollectionHandler,
		FeedPriceHistoryHandler: params.FeedPriceHistoryHandler,
	}
}

// validateAndParse parses the request body and validates the struct
func validateAndParse(c *fiber.Ctx, target any) error {
	if err := c.BodyParser(target); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	if err := utils.ValidateStruct(target); err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, errors.ErrValidationFailed.Message)
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
