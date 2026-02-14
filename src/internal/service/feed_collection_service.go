package service

import (
	"context"

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
	Update(ctx context.Context, request *model.FeedCollection, username string) error
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
		db:                 db,
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

	// Start transaction (ctx used so BaseModel hooks can set CreatedBy/UpdatedBy)
	tx := s.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create feed collection
	newFeedCollection := &model.FeedCollection{
		ClientId: clientId,
		Name:     request.Name,
		Unit:     request.Unit,
	}

	if err := tx.Create(newFeedCollection).Error; err != nil {
		tx.Rollback()
		return nil, errors.ErrGeneric.Wrap(err)
	}

	// Create feed price histories if provided
	var feedPriceHistories []interface{}
	if len(request.FeedPriceHistories) > 0 {
		priceHistories := make([]*model.FeedPriceHistory, 0, len(request.FeedPriceHistories))
		for _, priceHistoryReq := range request.FeedPriceHistories {
			priceHistory := &model.FeedPriceHistory{
				FeedCollectionId: newFeedCollection.Id,
				Price:            priceHistoryReq.Price,
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
			feedPriceHistories = append(feedPriceHistories, map[string]interface{}{
				"id":               ph.Id,
				"feedCollectionId": ph.FeedCollectionId,
				"price":            ph.Price,
				"priceUpdatedDate": ph.PriceUpdatedDate,
			})
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

func (s *feedCollectionService) Update(ctx context.Context, request *model.FeedCollection, username string) error {
	// Update feed collection (UpdatedBy set via BaseModel hook from ctx)
	if err := s.feedCollectionRepo.Update(ctx, request); err != nil {
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
		if fc.LatestPrice > 0 {
			response.LatestPrice = &fc.LatestPrice
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
	return &dto.FeedCollectionResponse{
		Id:        feedCollection.Id,
		ClientId:  feedCollection.ClientId,
		Name:      feedCollection.Name,
		Unit:      feedCollection.Unit,
		CreatedAt: feedCollection.CreatedAt,
		CreatedBy: feedCollection.CreatedBy,
		UpdatedAt: feedCollection.UpdatedAt,
		UpdatedBy: feedCollection.UpdatedBy,
	}
}

