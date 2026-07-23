package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/shopspring/decimal"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"

	"gorm.io/gorm"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedCollectionService --output=./mocks --outpkg=service --filename=feed_collection_service.go --structname=MockFeedCollectionService --with-expecter=false
type FeedCollectionService interface {
	Create(ctx context.Context, request dto.CreateFeedCollectionRequest, username string, clientId int) (*dto.CreateFeedCollectionResponse, error)
	Get(id int) (*dto.FeedCollectionResponse, error)
	Update(ctx context.Context, request dto.UpdateFeedCollectionRequest, username string) error
	GetPage(clientId, page, pageSize int, orderBy, keyword string) (*dto.PageResponse, error)
}

type feedCollectionService struct {
	feedCollectionRepo   repository.FeedCollectionRepository
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
	db                   *gorm.DB
}

func NewFeedCollectionService(
	feedCollectionRepo repository.FeedCollectionRepository,
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository,
	db *gorm.DB,
) FeedCollectionService {
	return &feedCollectionService{
		feedCollectionRepo:   feedCollectionRepo,
		feedPriceHistoryRepo: feedPriceHistoryRepo,
		db:                   db,
	}
}

func (s *feedCollectionService) Create(ctx context.Context, request dto.CreateFeedCollectionRequest, username string, clientId int) (*dto.CreateFeedCollectionResponse, error) {
	// Check if feed collection already exists
	checkFeedCollection, err := s.feedCollectionRepo.GetByClientIdAndName(clientId, request.Name)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkFeedCollection != nil {
		return nil, errors.ErrFeedCollectionAlreadyExists
	}

	feedType := request.FeedType
	if feedType == "" {
		feedType = constants.FeedTypePellet
	}
	if !constants.IsValidFeedType(feedType) {
		return nil, errors.ErrValidationFailed.Wrap(fmt.Errorf("invalid feedType"))
	}

	var fcr decimal.NullDecimal
	if request.Fcr != nil {
		fcr = decimal.NullDecimal{Decimal: decimal.NewFromFloat(*request.Fcr), Valid: true}
	}

	// pack_size_kg is NOT NULL — when the caller omits it, fall back to the
	// type-based default (fresh 30 / pellet 20 กก.) so the insert always resolves.
	packSize := constants.DefaultPackSizeKg(feedType)
	if request.PackSizeKg != nil {
		packSize = *request.PackSizeKg
	}
	packSizeKg := decimal.NullDecimal{Decimal: decimal.NewFromFloat(packSize), Valid: true}

	// Start transaction (ctx used so BaseModel hooks can set CreatedBy/UpdatedBy)
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create feed collection
	newFeedCollection := &model.FeedCollection{
		ClientId:   clientId,
		Name:       request.Name,
		Unit:       request.Unit,
		FeedType:   feedType,
		Fcr:        fcr,
		PackSizeKg: packSizeKg,
		Supplier:   normalizeSupplier(request.Supplier),
	}

	if err := tx.Create(newFeedCollection).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Create feed price histories if provided
	var feedPriceHistories []any
	if len(request.FeedPriceHistories) > 0 {
		priceHistories := make([]*model.FeedPriceHistory, 0, len(request.FeedPriceHistories))
		for _, priceHistoryReq := range request.FeedPriceHistories {
			var pricePerKg decimal.NullDecimal
			if priceHistoryReq.PricePerKg != nil {
				pricePerKg = decimal.NullDecimal{Decimal: decimal.NewFromFloat(*priceHistoryReq.PricePerKg), Valid: true}
			}
			priceHistory := &model.FeedPriceHistory{
				FeedCollectionId: newFeedCollection.Id,
				Price:            decimal.NewFromFloat(priceHistoryReq.Price),
				PricePerKg:       pricePerKg,
				PriceUpdatedDate: priceHistoryReq.PriceUpdatedDate,
			}
			priceHistories = append(priceHistories, priceHistory)
		}

		if err := tx.Create(priceHistories).Error; err != nil {
			tx.Rollback()
			return nil, errors.ErrGeneric.Wrap(err)
		}

		// Convert to response format
		for _, ph := range priceHistories {
			entry := map[string]any{
				"id":               ph.Id,
				"feedCollectionId": ph.FeedCollectionId,
				"price":            ph.Price.InexactFloat64(),
				"priceUpdatedDate": ph.PriceUpdatedDate,
			}
			if ph.PricePerKg.Valid {
				entry["pricePerKg"] = ph.PricePerKg.Decimal.InexactFloat64()
			}
			feedPriceHistories = append(feedPriceHistories, entry)
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return &dto.CreateFeedCollectionResponse{
		FeedCollection:   s.toFeedCollectionResponse(newFeedCollection),
		FeedPriceHistory: feedPriceHistories,
	}, nil
}

func (s *feedCollectionService) Get(id int) (*dto.FeedCollectionResponse, error) {
	feedCollection, err := s.feedCollectionRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if feedCollection == nil {
		return nil, errors.ErrFeedCollectionNotFound
	}

	return s.toFeedCollectionResponse(feedCollection), nil
}

func (s *feedCollectionService) Update(ctx context.Context, request dto.UpdateFeedCollectionRequest, username string) error {
	existingFeedCollection, err := s.feedCollectionRepo.GetByID(request.Id)
	if err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	if existingFeedCollection == nil {
		return errors.ErrFeedCollectionNotFound
	}

	if request.Name != "" {
		existingFeedCollection.Name = request.Name
	}
	if request.Unit != "" {
		existingFeedCollection.Unit = request.Unit
	}
	if request.FeedType != "" {
		if !constants.IsValidFeedType(request.FeedType) {
			return errors.ErrValidationFailed.Wrap(fmt.Errorf("invalid feedType"))
		}
		existingFeedCollection.FeedType = request.FeedType
	}
	if request.Fcr != nil {
		existingFeedCollection.Fcr = decimal.NullDecimal{Decimal: decimal.NewFromFloat(*request.Fcr), Valid: true}
	}
	if request.PackSizeKg != nil {
		existingFeedCollection.PackSizeKg = decimal.NullDecimal{Decimal: decimal.NewFromFloat(*request.PackSizeKg), Valid: true}
	}
	if request.Supplier != nil {
		existingFeedCollection.Supplier = normalizeSupplier(request.Supplier)
	}

	// Update feed collection (UpdatedBy set via BaseModel hook from ctx)
	if err := s.feedCollectionRepo.Update(ctx, existingFeedCollection); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *feedCollectionService) GetPage(clientId, page, pageSize int, orderBy, keyword string) (*dto.PageResponse, error) {
	feedCollections, total, err := s.feedCollectionRepo.GetPage(clientId, page, pageSize, orderBy, keyword)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.FeedCollectionPageResponse, 0, len(feedCollections))
	for _, fc := range feedCollections {
		response := &dto.FeedCollectionPageResponse{
			FeedCollectionResponse: *s.toFeedCollectionResponse(&fc.FeedCollection),
		}
		if fc.LatestPrice.GreaterThan(decimal.Zero) {
			v := fc.LatestPrice.InexactFloat64()
			response.LatestPrice = &v
		}
		if fc.LatestPriceUpdatedDate != nil {
			response.LatestPriceUpdatedDate = fc.LatestPriceUpdatedDate
		}
		responses = append(responses, response)
	}

	return &dto.PageResponse{
		Items: responses,
		Total: total,
	}, nil
}

func (s *feedCollectionService) toFeedCollectionResponse(feedCollection *model.FeedCollection) *dto.FeedCollectionResponse {
	resp := &dto.FeedCollectionResponse{
		Id:        feedCollection.Id,
		ClientId:  feedCollection.ClientId,
		Name:      feedCollection.Name,
		Unit:      feedCollection.Unit,
		FeedType:  feedCollection.FeedType,
		CreatedAt: feedCollection.CreatedAt,
		CreatedBy: feedCollection.CreatedBy,
		UpdatedAt: feedCollection.UpdatedAt,
		UpdatedBy: feedCollection.UpdatedBy,
	}
	if feedCollection.Fcr.Valid {
		v := feedCollection.Fcr.Decimal.InexactFloat64()
		resp.Fcr = &v
	}
	if feedCollection.PackSizeKg.Valid {
		v := feedCollection.PackSizeKg.Decimal.InexactFloat64()
		resp.PackSizeKg = &v
	}
	resp.Supplier = feedCollection.Supplier
	return resp
}

// normalizeSupplier trims the optional supplier string, returning nil for a
// nil or blank value so an empty field is stored as NULL rather than "".
func normalizeSupplier(s *string) *string {
	if s == nil {
		return nil
	}
	return lo.EmptyableToPtr(strings.TrimSpace(*s))
}
