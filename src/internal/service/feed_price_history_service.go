package service

import (
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=FeedPriceHistoryService --output=./mocks --outpkg=service --filename=feed_price_history_service.go --structname=MockFeedPriceHistoryService --with-expecter=false
type FeedPriceHistoryService interface {
	Create(request dto.CreateFeedPriceHistoryRequest, username string) (*dto.FeedPriceHistoryResponse, error)
	Get(id int) (*dto.FeedPriceHistoryResponse, error)
	Update(request *model.FeedPriceHistory, username string) error
	GetAll(feedCollectionId int) ([]*dto.FeedPriceHistoryResponse, error)
}

type feedPriceHistoryService struct {
	feedPriceHistoryRepo repository.FeedPriceHistoryRepository
}

func NewFeedPriceHistoryService(feedPriceHistoryRepo repository.FeedPriceHistoryRepository) FeedPriceHistoryService {
	return &feedPriceHistoryService{
		feedPriceHistoryRepo: feedPriceHistoryRepo,
	}
}

func (s *feedPriceHistoryService) Create(request dto.CreateFeedPriceHistoryRequest, username string) (*dto.FeedPriceHistoryResponse, error) {
	// Check if feed price history already exists
	checkFeedPriceHistory, err := s.feedPriceHistoryRepo.GetByFeedCollectionIdAndDate(request.FeedCollectionId, request.PriceUpdatedDate)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if checkFeedPriceHistory != nil {
		return nil, errors.ErrFeedPriceHistoryAlreadyExists
	}

	newFeedPriceHistory := &model.FeedPriceHistory{
		FeedCollectionId: request.FeedCollectionId,
		Price:            request.Price,
		PriceUpdatedDate: request.PriceUpdatedDate,
		BaseModel: model.BaseModel{
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	err = s.feedPriceHistoryRepo.Create(newFeedPriceHistory)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	return s.toFeedPriceHistoryResponse(newFeedPriceHistory), nil
}

func (s *feedPriceHistoryService) Get(id int) (*dto.FeedPriceHistoryResponse, error) {
	feedPriceHistory, err := s.feedPriceHistoryRepo.GetByID(id)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	if feedPriceHistory == nil {
		return nil, errors.ErrFeedPriceHistoryNotFound
	}

	return s.toFeedPriceHistoryResponse(feedPriceHistory), nil
}

func (s *feedPriceHistoryService) Update(request *model.FeedPriceHistory, username string) error {
	request.UpdatedBy = username
	if err := s.feedPriceHistoryRepo.Update(request); err != nil {
		return errors.ErrGeneric.Wrap(err)
	}
	return nil
}

func (s *feedPriceHistoryService) GetAll(feedCollectionId int) ([]*dto.FeedPriceHistoryResponse, error) {
	feedPriceHistories, err := s.feedPriceHistoryRepo.ListByFeedCollectionId(feedCollectionId)
	if err != nil {
		return nil, errors.ErrGeneric.Wrap(err)
	}

	responses := make([]*dto.FeedPriceHistoryResponse, 0, len(feedPriceHistories))
	for _, fph := range feedPriceHistories {
		responses = append(responses, s.toFeedPriceHistoryResponse(fph))
	}

	return responses, nil
}

func (s *feedPriceHistoryService) toFeedPriceHistoryResponse(fph *model.FeedPriceHistory) *dto.FeedPriceHistoryResponse {
	return &dto.FeedPriceHistoryResponse{
		Id:               fph.Id,
		FeedCollectionId: fph.FeedCollectionId,
		Price:            fph.Price,
		PriceUpdatedDate: fph.PriceUpdatedDate,
		CreatedAt:        fph.CreatedAt,
		CreatedBy:        fph.CreatedBy,
		UpdatedAt:        fph.UpdatedAt,
		UpdatedBy:        fph.UpdatedBy,
	}
}
