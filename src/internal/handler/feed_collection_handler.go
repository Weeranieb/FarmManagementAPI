package handler

import (
	"fmt"
	"strconv"

	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/service"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils/http"

	"github.com/gofiber/fiber/v2"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedCollectionHandler --output=./mocks --outpkg=handler --filename=feed_collection_handler.go --structname=MockFeedCollectionHandler --with-expecter=false
type FeedCollectionHandler interface {
	AddFeedCollection(c *fiber.Ctx) error
	GetFeedCollection(c *fiber.Ctx) error
	UpdateFeedCollection(c *fiber.Ctx) error
	ListFeedCollection(c *fiber.Ctx) error
}

type feedCollectionHandlerImpl struct {
	feedCollectionService service.FeedCollectionService
}

func NewFeedCollectionHandler(feedCollectionService service.FeedCollectionService) FeedCollectionHandler {
	return &feedCollectionHandlerImpl{
		feedCollectionService: feedCollectionService,
	}
}

// POST /api/v1/feedcollection
// Add a new feed collection.
// @Summary      Add a new feed collection
// @Description  Add a new feed collection with the provided details
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body dto.CreateFeedCollectionRequest true "Feed collection data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/feedcollection [post]
func (h *feedCollectionHandlerImpl) AddFeedCollection(c *fiber.Ctx) error {
	var createFeedCollectionRequest dto.CreateFeedCollectionRequest

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := validateAndParse(c, &createFeedCollectionRequest, errors.ErrValidationFailed.Code); err != nil {
		return err
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	// Get client id
	clientId, err := utils.GetClientId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	result, err := h.feedCollectionService.Create(createFeedCollectionRequest, username, clientId)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, result)
}

// GET /api/v1/feedcollection/:id
// Get a feed collection by ID.
// @Summary      Get a feed collection by ID
// @Description  Retrieve a feed collection by its ID
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        id path int true "Feed collection ID"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      404  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/feedcollection/{id} [get]
func (h *feedCollectionHandlerImpl) GetFeedCollection(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get id from param
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid feed collection ID")
	}

	feedCollection, err := h.feedCollectionService.Get(id)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, feedCollection)
}

// PUT /api/v1/feedcollection
// Update a feed collection.
// @Summary      Update a feed collection
// @Description  Update an existing feed collection with new details
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        Authorization header string true "Bearer token"
// @Param        body body model.FeedCollection true "Updated feed collection data"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/feedcollection [put]
func (h *feedCollectionHandlerImpl) UpdateFeedCollection(c *fiber.Ctx) error {
	var updateFeedCollection *model.FeedCollection

	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	if err := c.BodyParser(&updateFeedCollection); err != nil {
		return http.Error(c, errors.ErrInvalidRequestBody.Code, errors.ErrInvalidRequestBody.Message)
	}

	// Get username
	username, err := utils.GetUsername(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	err = h.feedCollectionService.Update(updateFeedCollection, username)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.SuccessWithoutData(c)
}

// GET /api/v1/feedcollection
// Get a list of feed collections with pagination.
// @Summary      Get a list of feed collections with pagination
// @Description  Retrieve a paginated list of feed collections for the current client
// @Tags         feedcollection
// @Accept       json
// @Produce      json
// @Param        page query int true "Page number"
// @Param        pageSize query int true "Page size"
// @Param        orderBy query string false "Order by field"
// @Param        keyword query string false "Search keyword"
// @Success      200  {object}  http.ResponseModel
// @Failure      400  {object}  http.ErrorResponseModel
// @Failure      500  {object}  http.ErrorResponseModel
// @Router       /api/v1/feedcollection [get]
func (h *feedCollectionHandlerImpl) ListFeedCollection(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			http.Error(c, errors.ErrGeneric.Code, fmt.Sprintf("%s: %v", errors.ErrGeneric.Message, r))
		}
	}()

	// Get query parameters
	sPage := c.Query("page")
	sPageSize := c.Query("pageSize")
	orderBy := c.Query("orderBy")
	keyword := c.Query("keyword")

	page, err := strconv.Atoi(sPage)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid page number")
	}

	pageSize, err := strconv.Atoi(sPageSize)
	if err != nil {
		return http.Error(c, errors.ErrValidationFailed.Code, "Invalid page size")
	}

	// Get client id
	clientId, err := utils.GetClientId(c.UserContext())
	if err != nil {
		return http.Error(c, errors.ErrAuthTokenInvalid.Code, errors.ErrAuthTokenInvalid.Message)
	}

	feedCollectionList, err := h.feedCollectionService.GetPage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return http.NewError(c, errors.ErrGeneric.Code, err)
	}

	return http.Success(c, feedCollectionList)
}

