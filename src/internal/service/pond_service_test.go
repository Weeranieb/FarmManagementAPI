//go:build cgo

package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
	svcmocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"
	"github.com/weeranieb/boonmafarm-backend/src/internal/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PondServiceTestSuite struct {
	suite.Suite
	pondRepo           *mocks.MockPondRepository
	farmRepo           *mocks.MockFarmRepository
	activePondRepo     *mocks.MockActivePondRepository
	activityRepo       *mocks.MockActivityRepository
	additionalCostRepo *mocks.MockAdditionalCostRepository
	sellDetailRepo     *mocks.MockSellDetailRepository
	merchantRepo       *mocks.MockMerchantRepository
	fishSizeGradeRepo  *mocks.MockFishSizeGradeRepository
	feedCostCalc       *svcmocks.MockFeedCostCalculator
	feedCostReturn     decimal.Decimal
	db                 *gorm.DB
	pondService        PondService
}

func (s *PondServiceTestSuite) SetupTest() {
	s.pondRepo = mocks.NewMockPondRepository(s.T())
	s.farmRepo = mocks.NewMockFarmRepository(s.T())
	s.activePondRepo = mocks.NewMockActivePondRepository(s.T())
	s.activityRepo = mocks.NewMockActivityRepository(s.T())
	s.additionalCostRepo = mocks.NewMockAdditionalCostRepository(s.T())
	s.sellDetailRepo = mocks.NewMockSellDetailRepository(s.T())
	s.merchantRepo = mocks.NewMockMerchantRepository(s.T())
	s.fishSizeGradeRepo = mocks.NewMockFishSizeGradeRepository(s.T())
	s.feedCostCalc = svcmocks.NewMockFeedCostCalculator(s.T())
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	err = s.db.AutoMigrate(&model.Pond{}, &model.ActivePond{}, &model.Activity{}, &model.AdditionalCost{})
	s.Require().NoError(err)
	s.pondService = NewPondService(PondServiceParams{
		PondRepo:           s.pondRepo,
		FarmRepo:           s.farmRepo,
		ActivePondRepo:     s.activePondRepo,
		ActivityRepo:       s.activityRepo,
		AdditionalCostRepo: s.additionalCostRepo,
		SellDetailRepo:     s.sellDetailRepo,
		MerchantRepo:       s.merchantRepo,
		FishSizeGradeRepo:  s.fishSizeGradeRepo,
		FeedCostCalc:       s.feedCostCalc,
		TxManager:          transaction.NewManager(s.db),
	})
	s.pondRepo.On("WithTx", mock.Anything).Maybe().Return(s.pondRepo)
	s.farmRepo.On("WithTx", mock.Anything).Maybe().Return(s.farmRepo)
	// Default feed cost is zero so existing net-result assertions (revenue −
	// cost) hold. Feed-cost tests set s.feedCostReturn before acting; the
	// function-based returns below read it at call time for both single and
	// batch paths.
	s.feedCostReturn = decimal.Zero
	s.feedCostCalc.On("CalcCycleFeedCost", mock.Anything, mock.Anything).Maybe().Return(
		func(_ context.Context, _ *model.ActivePond) decimal.Decimal { return s.feedCostReturn },
		func(_ context.Context, _ *model.ActivePond) error { return nil },
	)
	s.feedCostCalc.On("CalcCycleFeedCostBatch", mock.Anything, mock.Anything).Maybe().Return(
		func(_ context.Context, aps []*model.ActivePond) map[int]decimal.Decimal {
			m := make(map[int]decimal.Decimal, len(aps))
			for _, ap := range aps {
				if ap != nil {
					m[ap.Id] = s.feedCostReturn
				}
			}
			return m
		},
		func(_ context.Context, _ []*model.ActivePond) error { return nil },
	)
}

// fillPondCtx returns a context with super admin (userLevel 3) so CanAccessClient allows any client.
func fillPondCtx() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "testuser")
	ctx = context.WithValue(ctx, constants.UserLevelKey, 3)
	return ctx
}

// fillPondCtxNoAccess returns a context with normal user (clientId 1, userLevel 1) so CanAccessClient(clientId 2) is false.
func fillPondCtxNoAccess() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "user")
	ctx = context.WithValue(ctx, constants.ClientIDKey, 1)
	ctx = context.WithValue(ctx, constants.UserLevelKey, 1)
	return ctx
}

func (s *PondServiceTestSuite) TearDownTest() {
	s.pondRepo.ExpectedCalls = nil
	s.farmRepo.ExpectedCalls = nil
	s.activePondRepo.ExpectedCalls = nil
	s.activityRepo.ExpectedCalls = nil
	s.additionalCostRepo.ExpectedCalls = nil
	s.sellDetailRepo.ExpectedCalls = nil
	s.merchantRepo.ExpectedCalls = nil
	s.fishSizeGradeRepo.ExpectedCalls = nil
}

// mockFishSizeGradesForValidRequest mocks FishSizeGradeRepo.GetByIDs for the grade ID(s) used in validPondSellRequest (e.g. 1).
func (s *PondServiceTestSuite) mockFishSizeGradesForValidRequest() {
	s.fishSizeGradeRepo.On("GetByIDs", []int{1}).Return([]*model.FishSizeGrade{
		{Id: 1, Name: "6โล", SortIndex: 1},
	}, nil)
}

// setupReposWithTxForTransaction mocks WithTx to return the same mock; Create/Update assign IDs and return nil. Use Maybe() so tests that only Create or only Update still pass.
func (s *PondServiceTestSuite) setupReposWithTxForTransaction() {
	s.pondRepo.On("WithTx", mock.Anything).Maybe().Return(s.pondRepo)
	s.pondRepo.On("Update", mock.Anything, mock.Anything).Maybe().Return(nil)
	s.activePondRepo.On("WithTx", mock.Anything).Return(s.activePondRepo)
	s.activePondRepo.On("Create", mock.Anything, mock.Anything).Maybe().Return(nil).Run(func(args mock.Arguments) {
		ap := args.Get(1).(*model.ActivePond)
		if ap.Id == 0 {
			ap.Id = 99
		}
	})
	s.activePondRepo.On("Update", mock.Anything, mock.Anything).Maybe().Return(nil)
	s.activityRepo.On("WithTx", mock.Anything).Return(s.activityRepo)
	s.activityRepo.On("Create", mock.Anything, mock.Anything).Maybe().Return(nil).Run(func(args mock.Arguments) {
		a := args.Get(1).(*model.Activity)
		if a.Id == 0 {
			a.Id = 88
		}
	})
	s.additionalCostRepo.On("WithTx", mock.Anything).Return(s.additionalCostRepo)
	s.additionalCostRepo.On("CreateBatch", mock.Anything, mock.Anything).Maybe().Return(nil)
	s.sellDetailRepo.On("WithTx", mock.Anything).Return(s.sellDetailRepo)
	s.sellDetailRepo.On("CreateBatch", mock.Anything, mock.Anything).Maybe().Return(nil)
}

// expectFarmStatusSyncAfterMutation mocks pondRepo.ListByFarmId and farmRepo for syncFarmStatusFromPonds.
// pondsAfter is what ListByFarmId returns after the mutation; farmStoredStatus is the current farms.status row.
func (s *PondServiceTestSuite) expectFarmStatusSyncAfterMutation(farmId int, pondsAfter []*model.Pond, farmStoredStatus string) {
	s.pondRepo.On("ListByFarmId", farmId).Return(pondsAfter, nil)
	farm := &model.Farm{Id: farmId, ClientId: 1, Name: "Farm", Status: farmStoredStatus}
	s.farmRepo.On("GetByID", farmId).Return(farm, nil)
	want := utils.DeriveFarmStatusFromPonds(pondsAfter)
	if farmStoredStatus != want {
		s.farmRepo.On("Update", mock.Anything, mock.MatchedBy(func(f *model.Farm) bool {
			return f.Id == farmId && f.Status == want
		})).Return(nil)
	}
}

func TestPondServiceSuite(t *testing.T) {
	suite.Run(t, new(PondServiceTestSuite))
}

func (s *PondServiceTestSuite) TestCreatePonds_Success() {
	// GIVEN — request with farm and names; repo returns no duplicate names.
	// The pre-check GetByID(1) is satisfied by expectFarmStatusSyncAfterMutation
	// below (which registers the same expectation for the in-transaction sync).
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "Pond 1"}, {Name: "Pond 2"}},
	}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(nil, nil)
	s.pondRepo.On("CreateBatch", mock.Anything, mock.AnythingOfType("[]*model.Pond")).Return(nil).Run(func(args mock.Arguments) {
		ponds := args.Get(1).([]*model.Pond)
		for i := range ponds {
			ponds[i].Id = i + 1
			ponds[i].CreatedAt = time.Now()
			ponds[i].UpdatedAt = time.Now()
		}
	})
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{
		{Id: 1, FarmId: 1, Status: constants.FarmStatusMaintenance},
		{Id: 2, FarmId: 1, Status: constants.FarmStatusMaintenance},
	}, constants.FarmStatusMaintenance)

	// WHEN — CreatePonds is called (super-admin context)
	err := s.pondService.CreatePonds(fillPondCtx(), req)

	// THEN — no error; CreateBatch was used
	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestCreatePonds_PondAlreadyExists() {
	// GIVEN — request; second name already exists for this farm
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Ponds:  []dto.CreatePondItem{{Name: "Pond 1"}, {Name: "Pond 2"}},
	}
	// Pre-check: farm exists in client 1; super-admin context can access.
	s.farmRepo.On("GetByID", 1).Return(&model.Farm{Id: 1, ClientId: 1, Name: "F", Status: constants.FarmStatusMaintenance}, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	existingPond := &model.Pond{Id: 99, FarmId: 1, Name: "Pond 2", Status: "active"}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(existingPond, nil)

	// WHEN — CreatePonds is called (super-admin context)
	err := s.pondService.CreatePonds(fillPondCtx(), req)

	// THEN — ErrPondAlreadyExists; CreateBatch not called
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
	s.pondRepo.AssertNotCalled(s.T(), "CreateBatch")
}

func (s *PondServiceTestSuite) TestCreatePonds_FarmNotFound() {
	// GIVEN — farm does not exist
	s.farmRepo.On("GetByID", 99).Return(nil, nil)

	// WHEN
	err := s.pondService.CreatePonds(fillPondCtx(), dto.CreatePondsRequest{
		FarmId: 99,
		Ponds:  []dto.CreatePondItem{{Name: "P1"}},
	})

	// THEN — ErrFarmNotFound; no pond lookup
	assert.ErrorIs(s.T(), err, errors.ErrFarmNotFound)
	s.pondRepo.AssertNotCalled(s.T(), "GetByFarmIdAndName")
	s.pondRepo.AssertNotCalled(s.T(), "CreateBatch")
}

func (s *PondServiceTestSuite) TestCreatePonds_ClientAdminWrongClientDenied() {
	// GIVEN — client admin for client 1 trying to add ponds to farm in client 2
	s.farmRepo.On("GetByID", 5).Return(&model.Farm{Id: 5, ClientId: 2, Name: "F"}, nil)

	// Context: client admin (level 2) tied to client 1.
	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.UsernameKey, "admin")
	ctx = context.WithValue(ctx, constants.ClientIDKey, 1)
	ctx = context.WithValue(ctx, constants.UserLevelKey, 2)

	// WHEN
	err := s.pondService.CreatePonds(ctx, dto.CreatePondsRequest{
		FarmId: 5,
		Ponds:  []dto.CreatePondItem{{Name: "P1"}},
	})

	// THEN — permission denied; no pond operations attempted
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
	s.pondRepo.AssertNotCalled(s.T(), "GetByFarmIdAndName")
	s.pondRepo.AssertNotCalled(s.T(), "CreateBatch")
}

func (s *PondServiceTestSuite) TestGet_Success() {
	// GIVEN — pond exists with farm and client
	pondId := 1
	pa := &repository.PondWithFarmAndActivePond{
		Pond: &model.Pond{
			Id:     pondId,
			FarmId: 1,
			Name:   "Test Pond",
			Status: "active",
		},
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(pa, nil)

	// WHEN — Get is called
	result, err := s.pondService.Get(context.Background(), pondId)

	// THEN — result returned with same id
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), pondId, result.Id)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGet_NotFound() {
	// GIVEN — pond id does not exist (repo returns nil)
	pondId := 999
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, nil)

	// WHEN — Get is called
	result, err := s.pondService.Get(context.Background(), pondId)

	// THEN — ErrPondNotFound; no result
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGet_RepoError() {
	// GIVEN — repo returns an error
	pondId := 1
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, assert.AnError)

	// WHEN — Get is called
	result, err := s.pondService.Get(context.Background(), pondId)

	// THEN — error propagated; no result
	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGetList_Success() {
	// GIVEN — farm has two ponds; repo returns them
	farmId := 1
	list := []*repository.PondWithFarmAndActivePond{
		{Pond: &model.Pond{Id: 1, FarmId: farmId, Name: "Pond 1", Status: "active"}, ClientId: 1, ActivePond: nil},
		{Pond: &model.Pond{Id: 2, FarmId: farmId, Name: "Pond 2", Status: "active"}, ClientId: 1, ActivePond: nil},
	}
	s.pondRepo.On("ListByFarmIdWithActivePond", mock.Anything, farmId).Return(list, nil)

	// WHEN — GetList is called
	result, err := s.pondService.GetList(context.Background(), farmId)

	// THEN — two ponds returned
	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestUpdate_Success() {
	// GIVEN — existing pond; new name not taken
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Old Name", Status: "maintenance"}
	req := dto.UpdatePondRequest{Id: 1, Name: "New Name", Status: "active"}
	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "New Name").Return(nil, nil)
	s.pondRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pond")).Return(nil)
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{
		{Id: 1, FarmId: 1, Name: "New Name", Status: constants.FarmStatusActive},
	}, constants.FarmStatusMaintenance)

	// WHEN — Update is called
	err := s.pondService.Update(context.Background(), req)

	// THEN — no error
	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestUpdate_PondNotFound() {
	// GIVEN — pond id does not exist
	req := dto.UpdatePondRequest{Id: 999, Name: "Pond"}
	s.pondRepo.On("GetByID", 999).Return(nil, nil)

	// WHEN — Update is called
	err := s.pondService.Update(context.Background(), req)

	// THEN — ErrPondNotFound; Update not called
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
	s.pondRepo.AssertNotCalled(s.T(), "Update")
}

func (s *PondServiceTestSuite) TestUpdate_DuplicateName() {
	// GIVEN — existing pond; new name already taken by another pond
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Old", Status: "active"}
	otherPond := &model.Pond{Id: 2, FarmId: 1, Name: "New Name", Status: "active"}
	req := dto.UpdatePondRequest{Id: 1, Name: "New Name"}
	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "New Name").Return(otherPond, nil)

	// WHEN — Update is called
	err := s.pondService.Update(context.Background(), req)

	// THEN — ErrPondAlreadyExists; Update not called
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
	s.pondRepo.AssertNotCalled(s.T(), "Update")
}

func (s *PondServiceTestSuite) TestUpdate_RepoError() {
	// GIVEN — existing pond; Update will return error
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Pond", Status: "active"}
	req := dto.UpdatePondRequest{Id: 1, Status: "maintenance"}
	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pond")).Return(assert.AnError)

	// WHEN — Update is called
	err := s.pondService.Update(context.Background(), req)

	// THEN — error propagated
	assert.Error(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_PondNotFound() {
	// GIVEN — pond id does not exist
	pondId := 999
	req := validPondFillRequest()
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, nil)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrPondNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_RepoError() {
	// GIVEN — repo returns error
	pondId := 1
	req := validPondFillRequest()
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, assert.AnError)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — error propagated; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_FarmNotFound() {
	// GIVEN — pond data has no client (ClientId 0)
	pondId := 1
	req := validPondFillRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   0,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrFarmNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrFarmNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_PermissionDenied() {
	// GIVEN — pond belongs to client 2; user has access only to client 1
	pondId := 1
	req := validPondFillRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   2,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — FillPond is called with no-access context
	resp, err := s.pondService.FillPond(fillPondCtxNoAccess(), pondId, req, "user")

	// THEN — ErrAuthPermissionDenied; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_InvalidFishType() {
	// GIVEN — valid pond data; request has invalid fish type
	pondId := 1
	req := validPondFillRequest()
	req.FishType = "invalid"
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrInvalidFishType; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrInvalidFishType)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_InvalidActivityDate() {
	// GIVEN — valid pond data; request has invalid activity date
	pondId := 1
	req := validPondFillRequest()
	req.ActivityDate = "not-a-date"
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — validation error; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), errors.ErrValidationFailed.Message)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_Success_NewActivePond() {
	// GIVEN — pond in maintenance (no active cycle); tx mocks set up
	pondId := 1
	req := validPondFillRequest()
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "Pond", Status: constants.FarmStatusMaintenance}
	data := &repository.PondWithFarmAndActivePond{
		Pond:       pond,
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.setupReposWithTxForTransaction()
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{
		{Id: pondId, FarmId: 1, Name: "Pond", Status: constants.FarmStatusActive},
	}, constants.FarmStatusMaintenance)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — success; new active pond and activity ids returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivePondId, int64(0))
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_Success_ExistingActivePond() {
	// GIVEN — pond already has active cycle; tx mocks set up
	pondId := 1
	req := validPondFillRequest()
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "Pond", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id:          10,
		PondId:      pondId,
		IsActive:    true,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
	}
	data := &repository.PondWithFarmAndActivePond{
		Pond:       pond,
		ClientId:   1,
		ActivePond: activePond,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.setupReposWithTxForTransaction()
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{
		{Id: pondId, FarmId: 1, Name: "Pond", Status: constants.FarmStatusActive},
	}, constants.FarmStatusActive)

	// WHEN — FillPond is called
	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	// THEN — success; existing active pond id and new activity id returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func validPondFillRequest() dto.PondFillRequest {
	return dto.PondFillRequest{
		FishType:     constants.FishTypeNil,
		Amount:       100,
		FishWeight:   decimal.RequireFromString("0.5"),
		PricePerUnit: decimal.RequireFromString("10"),
		ActivityDate: "2025-01-15",
		AdditionalCosts: []dto.AdditionalCostItem{
			{Title: "Transport", Cost: decimal.RequireFromString("50")},
		},
	}
}

func validPondMoveRequest() dto.PondMoveRequest {
	return dto.PondMoveRequest{
		ToPondId:     2,
		FishType:     constants.FishTypeNil,
		Amount:       50,
		ActivityDate: "2025-06-01",
	}
}

func (s *PondServiceTestSuite) TestMovePond_SourceNotFound() {
	// GIVEN — source pond id does not exist
	req := validPondMoveRequest()
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, 1).Return(nil, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), 1, req, "user")

	// THEN — ErrPondNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_SourceNotActive() {
	// GIVEN — source pond has no active cycle (in maintenance)
	sourcePondId := 1
	req := validPondMoveRequest()
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusMaintenance},
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrPondInMaintenance (maintenance ponds cannot be source for move)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondInMaintenance)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_DestNotFound() {
	// GIVEN — source valid; destination pond id does not exist
	sourcePondId := 1
	req := validPondMoveRequest()
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(nil, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrPondNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_DestDifferentClient_ReturnsPermissionDenied() {
	// GIVEN — source client 1; destination belongs to client 2
	sourcePondId := 1
	req := validPondMoveRequest()
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	destData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: req.ToPondId, FarmId: 2, Name: "P2", Status: constants.FarmStatusActive},
		ClientId:   2,
		ActivePond: &model.ActivePond{Id: 20, PondId: req.ToPondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(destData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrAuthPermissionDenied; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_SamePond_ReturnsInvalidInput() {
	// GIVEN — source and destination are the same pond
	sourcePondId := 1
	req := validPondMoveRequest()
	req.ToPondId = sourcePondId
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrPondInvalidInput; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondInvalidInput)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_AmountExceedsSourceStock_ReturnsInsufficientFish() {
	// GIVEN — source has 100 fish, request asks to move 150
	sourcePondId := 1
	req := validPondMoveRequest()
	req.Amount = 150
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	destData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 20, PondId: req.ToPondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(destData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrPondInsufficientFish; nothing persisted
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondInsufficientFish)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_NonPositiveAmount_ReturnsInvalidInput() {
	// GIVEN — request asks to move zero fish
	sourcePondId := 1
	req := validPondMoveRequest()
	req.Amount = 0
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrPondInvalidInput; no destination lookup performed
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondInvalidInput)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_InvalidFishType() {
	// GIVEN — source and dest valid; request has invalid fish type
	sourcePondId := 1
	req := validPondMoveRequest()
	req.FishType = "invalid"
	sourceData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: sourcePondId, IsActive: true, TotalFish: 100},
	}
	destData := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 20, PondId: req.ToPondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(sourceData, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(destData, nil)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — ErrInvalidFishType; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrInvalidFishType)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_Success_BothActive() {
	// GIVEN — source and dest both have active cycles; same client; tx mocks set up
	sourcePondId := 1
	req := validPondMoveRequest()
	sourcePond := &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	sourceActive := &model.ActivePond{
		Id:          10,
		PondId:      sourcePondId,
		IsActive:    true,
		TotalFish:   100,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
		FishTypes:   []string{constants.FishTypeNil},
	}
	destPond := &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusActive}
	destActive := &model.ActivePond{
		Id:          20,
		PondId:      req.ToPondId,
		IsActive:    true,
		TotalFish:   0,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
		FishTypes:   []string{},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: sourcePond, ClientId: 1, ActivePond: sourceActive,
	}, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: destPond, ClientId: 1, ActivePond: destActive,
	}, nil)
	s.setupReposWithTxForTransaction()
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{sourcePond, destPond}, constants.FarmStatusActive)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — success; activity and both active pond ids returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.Equal(s.T(), int64(20), resp.ToActivePondId)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_Success_DestInMaintenance() {
	// GIVEN — source active; dest in maintenance (no active cycle); tx mocks set up
	sourcePondId := 1
	req := validPondMoveRequest()
	sourcePond := &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	sourceActive := &model.ActivePond{
		Id:          10,
		PondId:      sourcePondId,
		IsActive:    true,
		TotalFish:   100,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
		FishTypes:   []string{constants.FishTypeNil},
	}
	destPond := &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusMaintenance}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: sourcePond, ClientId: 1, ActivePond: sourceActive,
	}, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: destPond, ClientId: 1, ActivePond: nil,
	}, nil)
	s.setupReposWithTxForTransaction()
	destAfter := &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusActive}
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{sourcePond, destAfter}, constants.FarmStatusActive)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — success; new dest active pond created
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.Greater(s.T(), resp.ToActivePondId, int64(0))
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestMovePond_Success_MarkToClose() {
	// GIVEN — source and dest active; MarkToClose true; capture pond Update
	sourcePondId := 1
	req := validPondMoveRequest()
	req.MarkToClose = true
	sourcePond := &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	sourceActive := &model.ActivePond{
		Id:          10,
		PondId:      sourcePondId,
		IsActive:    true,
		TotalFish:   100,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
		FishTypes:   []string{constants.FishTypeNil},
	}
	destPond := &model.Pond{Id: req.ToPondId, FarmId: 1, Name: "P2", Status: constants.FarmStatusActive}
	destActive := &model.ActivePond{
		Id:          20,
		PondId:      req.ToPondId,
		IsActive:    true,
		TotalFish:   0,
		TotalCost:   decimal.Zero,
		TotalProfit: decimal.Zero,
		NetResult:   decimal.Zero,
		FishTypes:   []string{},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, sourcePondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: sourcePond, ClientId: 1, ActivePond: sourceActive,
	}, nil)
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, req.ToPondId).Return(&repository.PondWithFarmAndActivePond{
		Pond: destPond, ClientId: 1, ActivePond: destActive,
	}, nil)

	var updatedPond *model.Pond
	s.pondRepo.On("Update", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
		if p, ok := args.Get(1).(*model.Pond); ok && p.Id == sourcePondId {
			updatedPond = p
		}
	}).Return(nil)
	s.setupReposWithTxForTransaction()
	sourceAfter := &model.Pond{Id: sourcePondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusMaintenance}
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{sourceAfter, destPond}, constants.FarmStatusActive)

	// WHEN — MovePond is called
	resp, err := s.pondService.MovePond(fillPondCtx(), sourcePondId, req, "user")

	// THEN — success; source pond updated to maintenance
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.Equal(s.T(), int64(20), resp.ToActivePondId)
	assert.NotNil(s.T(), updatedPond, "pondRepo.Update should be called for source pond when MarkToClose is true")
	assert.Equal(s.T(), constants.FarmStatusMaintenance, updatedPond.Status)
	assert.Equal(s.T(), 0, sourceActive.TotalFish, "closing via move empties the pond, same as a sell that closes")
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_Success_WithAdditionalCosts() {
	// GIVEN — pond with active cycle; sell request includes additionalCosts
	pondId := 1
	req := validPondSellRequest()
	req.AdditionalCosts = []dto.AdditionalCostItem{
		{Title: "Transport", Cost: decimal.RequireFromString("200")},
		{Title: "Packaging", Cost: decimal.RequireFromString("50")},
	}
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id:          10,
		PondId:      pondId,
		IsActive:    true,
		TotalCost:   decimal.RequireFromString("1000"),
		TotalProfit: decimal.Zero,
		NetResult:   decimal.RequireFromString("-1000"),
		FishTypes:   []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{
		Pond: pond, ClientId: 1, ActivePond: activePond,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(activePond, nil)
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{pond}, constants.FarmStatusActive)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — success; additional costs persisted via CreateBatch
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	s.additionalCostRepo.AssertCalled(s.T(), "CreateBatch", mock.Anything, mock.MatchedBy(func(items []*model.AdditionalCost) bool {
		return len(items) == 2 && items[0].Title == "Transport" && items[1].Title == "Packaging"
	}))
	s.farmRepo.AssertExpectations(s.T())
}

func validPondSellRequest() dto.PondSellRequest {
	return dto.PondSellRequest{
		ActivityDate: "2025-07-01",
		Details: []dto.PondSellDetailItem{
			{
				FishSizeGradeId: 1,
				Weight:          decimal.RequireFromString("100"),
				PricePerUnit:    decimal.RequireFromString("50"),
			},
		},
	}
}

func TestBuildSellDetailModels(t *testing.T) {
	// GIVEN — activity id and two detail items
	details := []dto.PondSellDetailItem{
		{FishSizeGradeId: 1, Weight: decimal.RequireFromString("10"), PricePerUnit: decimal.RequireFromString("5")},
		{FishSizeGradeId: 2, Weight: decimal.RequireFromString("20"), PricePerUnit: decimal.RequireFromString("10")},
	}

	// WHEN — buildSellDetailModels is called
	out := buildSellDetailModels(99, details)

	// THEN — two models with correct SellId and fields
	require.Len(t, out, 2)
	assert.Equal(t, 99, out[0].SellId)
	assert.Equal(t, 1, out[0].FishSizeGradeId)
	assert.True(t, out[0].Weight.Equal(decimal.RequireFromString("10")))
	assert.True(t, out[0].PricePerUnit.Equal(decimal.RequireFromString("5")))
	assert.Equal(t, 99, out[1].SellId)
	assert.Equal(t, 2, out[1].FishSizeGradeId)
	assert.True(t, out[1].Weight.Equal(decimal.RequireFromString("20")))
}

func (s *PondServiceTestSuite) TestSellPond_PondNotFound() {
	// GIVEN — pond id does not exist
	pondId := 1
	req := validPondSellRequest()
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, nil)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrPondNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_PondNotActive() {
	// GIVEN — pond has no active cycle
	pondId := 1
	req := validPondSellRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusMaintenance},
		ClientId:   1,
		ActivePond: nil,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrPondInMaintenance (maintenance ponds cannot sell)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondInMaintenance)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_FarmNotFound() {
	// GIVEN — pond data has ClientId 0
	pondId := 1
	req := validPondSellRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   0,
		ActivePond: &model.ActivePond{Id: 10, PondId: pondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrFarmNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrFarmNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_PermissionDenied() {
	// GIVEN — pond belongs to client 2; user has no access
	pondId := 1
	req := validPondSellRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   2,
		ActivePond: &model.ActivePond{Id: 10, PondId: pondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	// WHEN — SellPond is called with no-access context
	resp, err := s.pondService.SellPond(fillPondCtxNoAccess(), pondId, req, "user")

	// THEN — ErrAuthPermissionDenied; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_MerchantNotFound() {
	// GIVEN — pond exists with active cycle; request has unknown merchantId
	pondId := 1
	merchantId := 5
	req := validPondSellRequest()
	req.MerchantId = &merchantId
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: pondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()
	s.merchantRepo.On("GetByID", merchantId).Return(nil, nil)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — ErrMerchantNotFound; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrMerchantNotFound)
	s.pondRepo.AssertExpectations(s.T())
	s.merchantRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_InvalidActivityDate() {
	// GIVEN — valid pond and active cycle; request has invalid activity date
	pondId := 1
	req := validPondSellRequest()
	req.ActivityDate = "invalid"
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive},
		ClientId:   1,
		ActivePond: &model.ActivePond{Id: 10, PondId: pondId, IsActive: true},
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — validation error; no response
	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorContains(s.T(), err, "Validation")
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_Success() {
	// GIVEN — pond with active cycle; valid sell request; tx mocks set up
	pondId := 1
	req := validPondSellRequest()
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id:          10,
		PondId:      pondId,
		IsActive:    true,
		TotalCost:   decimal.RequireFromString("1000"),
		TotalProfit: decimal.Zero,
		NetResult:   decimal.RequireFromString("-1000"),
		FishTypes:   []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{
		Pond: pond, ClientId: 1, ActivePond: activePond,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(activePond, nil)
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{pond}, constants.FarmStatusActive)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — success; activity and active pond ids returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestSellPond_Success_MarkToClose() {
	// GIVEN — pond with active cycle; MarkToClose true; capture pond Update
	pondId := 1
	req := validPondSellRequest()
	req.MarkToClose = true
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id:          10,
		PondId:      pondId,
		IsActive:    true,
		TotalCost:   decimal.RequireFromString("1000"),
		TotalProfit: decimal.Zero,
		NetResult:   decimal.RequireFromString("-1000"),
		FishTypes:   []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{
		Pond: pond, ClientId: 1, ActivePond: activePond,
	}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	var updatedPond *model.Pond
	s.pondRepo.On("Update", mock.Anything, mock.Anything).Maybe().Run(func(args mock.Arguments) {
		if p, ok := args.Get(1).(*model.Pond); ok && p.Id == pondId {
			updatedPond = p
		}
	}).Return(nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(activePond, nil)
	pondAfter := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusMaintenance}
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{pondAfter}, constants.FarmStatusActive)

	// WHEN — SellPond is called
	resp, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	// THEN — success; pond updated to maintenance
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.NotNil(s.T(), updatedPond, "pondRepo.Update should be called when MarkToClose is true")
	assert.Equal(s.T(), constants.FarmStatusMaintenance, updatedPond.Status)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

// TestSellPond_ReducesTotalFish verifies a (non-closing) sale draws the pond's
// head count down by the sum of the sell lines' FishCount.
func (s *PondServiceTestSuite) TestSellPond_ReducesTotalFish() {
	pondId := 1
	req := validPondSellRequest()
	fishCount := 200
	req.Details[0].FishCount = &fishCount
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id: 10, PondId: pondId, IsActive: true,
		TotalCost: decimal.RequireFromString("1000"), TotalProfit: decimal.Zero,
		TotalFish: 500, FishTypes: []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{Pond: pond, ClientId: 1, ActivePond: activePond}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(activePond, nil)
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{pond}, constants.FarmStatusActive)

	_, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 300, activePond.TotalFish, "500 stocked − 200 sold")
	assert.False(s.T(), activePond.FeedCost.Valid, "feed cost is only snapshotted on close")
}

// TestSellPond_InsufficientFish rejects a sale of more fish than the pond holds.
func (s *PondServiceTestSuite) TestSellPond_InsufficientFish() {
	pondId := 1
	req := validPondSellRequest()
	fishCount := 999
	req.Details[0].FishCount = &fishCount
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id: 10, PondId: pondId, IsActive: true, TotalFish: 500,
		FishTypes: []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{Pond: pond, ClientId: 1, ActivePond: activePond}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()

	_, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	assert.ErrorIs(s.T(), err, errors.ErrPondInsufficientFish)
}

// TestSellPond_ConcurrentOversell_Rejected covers the race the row lock guards:
// the pre-transaction check passes against a stale head count (500), but by the
// time the transaction acquires the lock a concurrent sell has drawn stock down
// to 50, so the locked re-check rejects the oversell — and the specific
// ErrPondInsufficientFish (422) survives the transaction wrapper rather than
// collapsing to a generic 500.
func (s *PondServiceTestSuite) TestSellPond_ConcurrentOversell_Rejected() {
	pondId := 1
	req := validPondSellRequest()
	fishCount := 100
	req.Details[0].FishCount = &fishCount
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	// Stale snapshot seen before the transaction: 500 fish, so 100 passes the
	// pre-check.
	activePond := &model.ActivePond{
		Id: 10, PondId: pondId, IsActive: true, TotalFish: 500,
		FishTypes: []string{constants.FishTypeNil},
	}
	// Authoritative locked row: a concurrent sell already dropped stock to 50.
	locked := &model.ActivePond{
		Id: 10, PondId: pondId, IsActive: true, TotalFish: 50,
		FishTypes: []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{Pond: pond, ClientId: 1, ActivePond: activePond}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(locked, nil)

	_, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	assert.ErrorIs(s.T(), err, errors.ErrPondInsufficientFish)
}

// TestSellPond_MarkToClose_SnapshotsFeedCost verifies closing empties the pond,
// freezes the cycle feed cost into feed_cost, and folds it into net_result.
func (s *PondServiceTestSuite) TestSellPond_MarkToClose_SnapshotsFeedCost() {
	pondId := 1
	req := validPondSellRequest() // revenue = 100kg × 50 = 5000
	req.MarkToClose = true
	fishCount := 300
	req.Details[0].FishCount = &fishCount
	s.feedCostReturn = decimal.RequireFromString("1500")
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	activePond := &model.ActivePond{
		Id: 10, PondId: pondId, IsActive: true,
		TotalCost: decimal.RequireFromString("1000"), TotalProfit: decimal.Zero,
		TotalFish: 500, FishTypes: []string{constants.FishTypeNil},
	}
	data := &repository.PondWithFarmAndActivePond{Pond: pond, ClientId: 1, ActivePond: activePond}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)
	s.pondRepo.On("Update", mock.Anything, mock.Anything).Maybe().Return(nil)
	s.mockFishSizeGradesForValidRequest()
	s.setupReposWithTxForTransaction()
	s.activePondRepo.On("GetByIDForUpdate", mock.Anything, activePond.Id).Return(activePond, nil)
	pondAfter := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusMaintenance}
	s.expectFarmStatusSyncAfterMutation(1, []*model.Pond{pondAfter}, constants.FarmStatusActive)

	_, err := s.pondService.SellPond(fillPondCtx(), pondId, req, "user")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, activePond.TotalFish, "closing empties the pond")
	assert.False(s.T(), activePond.IsActive)
	require.True(s.T(), activePond.FeedCost.Valid)
	assert.True(s.T(), decimal.RequireFromString("1500").Equal(activePond.FeedCost.Decimal))
	// net = revenue 5000 − cost 1000 − feed 1500 = 2500
	assert.True(s.T(), decimal.RequireFromString("2500").Equal(activePond.NetResult),
		"net_result should include feed cost, got %s", activePond.NetResult)
}

// TestListCycles verifies cycle history P&L: the active cycle derives feed cost
// live (net = revenue − cost − feed), a post-feature closed cycle uses its
// snapshot, and a legacy closed cycle (feed_cost NULL) is returned as stored.
func (s *PondServiceTestSuite) TestListCycles() {
	pondId := 1
	s.feedCostReturn = decimal.RequireFromString("700")
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "P1", Status: constants.FarmStatusActive}
	data := &repository.PondWithFarmAndActivePond{Pond: pond, ClientId: 1}
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	end := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	feed1500 := decimal.RequireFromString("1500")
	cycles := []*model.ActivePond{
		// active: net = 2000 − 1000 − 700(derived) = 300
		{Id: 30, PondId: pondId, IsActive: true, TotalCost: decimal.RequireFromString("1000"),
			TotalProfit: decimal.RequireFromString("2000"), NetResult: decimal.RequireFromString("1000"), TotalFish: 500},
		// closed post-feature: snapshot feed 1500, stored net 1500
		{Id: 20, PondId: pondId, IsActive: false, EndDate: &end, TotalCost: decimal.RequireFromString("5000"),
			TotalProfit: decimal.RequireFromString("8000"), NetResult: feed1500, FeedCost: decimal.NullDecimal{Decimal: feed1500, Valid: true}},
		// legacy closed: feed_cost NULL, stored net 1000 (pre-feed)
		{Id: 10, PondId: pondId, IsActive: false, EndDate: &end, TotalCost: decimal.RequireFromString("3000"),
			TotalProfit: decimal.RequireFromString("4000"), NetResult: decimal.RequireFromString("1000")},
	}
	s.activePondRepo.On("ListByPondID", mock.Anything, pondId).Return(cycles, nil)

	out, err := s.pondService.ListCycles(fillPondCtx(), pondId)

	require.NoError(s.T(), err)
	require.Len(s.T(), out, 3)
	// active
	assert.True(s.T(), out[0].IsActive)
	require.NotNil(s.T(), out[0].FeedCost)
	assert.InDelta(s.T(), 700, *out[0].FeedCost, 0.001)
	assert.InDelta(s.T(), 300, out[0].NetResult, 0.001)
	// closed post-feature: snapshot
	assert.False(s.T(), out[1].IsActive)
	require.NotNil(s.T(), out[1].FeedCost)
	assert.InDelta(s.T(), 1500, *out[1].FeedCost, 0.001)
	assert.InDelta(s.T(), 1500, out[1].NetResult, 0.001)
	// legacy closed: no feed cost, stored net unchanged
	assert.Nil(s.T(), out[2].FeedCost, "legacy closed cycle has no feed snapshot")
	assert.InDelta(s.T(), 1000, out[2].NetResult, 0.001)
}

// TestListCycles_PondNotFound returns a not-found error when the pond is missing.
func (s *PondServiceTestSuite) TestListCycles_PondNotFound() {
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, 99).Return((*repository.PondWithFarmAndActivePond)(nil), nil)
	_, err := s.pondService.ListCycles(fillPondCtx(), 99)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
}

// --- BulkImportFarmPond validation ---
// validateBulkImportRequest is a pure function; no repo mocks needed.

func newDecimal(v string) *decimal.Decimal {
	d := decimal.RequireFromString(v)
	return &d
}

func (s *PondServiceTestSuite) TestBulkImportValidate_HappyPath() {
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{
				{Name: "P1", Area: newDecimal("2.5")},
				{Name: "P2"}, // no area is fine
			}},
		},
	})
	assert.NoError(s.T(), err)
}

func (s *PondServiceTestSuite) TestBulkImportValidate_DuplicateInsideRequest() {
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{
				{Name: "P1"},
				{Name: "P2"},
				{Name: "p1"}, // case-insensitive dup of #1
			}},
		},
	})
	require.Error(s.T(), err)
	// Message reports the duplicate-row entry (lowercased "p1") and points
	// at "item #1" so the user can find the original.
	assert.Contains(s.T(), err.Error(), "duplicate pond")
	assert.Contains(s.T(), err.Error(), "p1")
	assert.Contains(s.T(), err.Error(), "item #1")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_NegativeAreaRejected() {
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{
				{Name: "P1", Area: newDecimal("-1")},
			}},
		},
	})
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "area must be >= 0")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_EmptyFarmNameAfterNormalize() {
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			// "ฟาร์ม " is the Thai display prefix; normalize strips it to "".
			{Name: "ฟาร์ม ", Ponds: []dto.BulkImportPondItem{{Name: "P1"}}},
		},
	})
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "empty name")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_EmptyPondNameAfterNormalize() {
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			// "บ่อ " normalizes to "".
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{{Name: "บ่อ "}}},
		},
	})
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "empty pond name")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_NameTooLong() {
	svc := s.pondService.(*pondService)
	long := make([]byte, 101)
	for i := range long {
		long[i] = 'a'
	}
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: string(long), Ponds: []dto.BulkImportPondItem{{Name: "P1"}}},
		},
	})
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "exceeds 100 chars")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_CollectsAllIssues() {
	// Multiple distinct issues across the payload should be reported together
	// so the user can fix the whole file in one pass.
	svc := s.pondService.(*pondService)
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{
				{Name: "P1", Area: newDecimal("-2")}, // negative area
				{Name: "P1"},                         // duplicate of #1
			}},
		},
	})
	require.Error(s.T(), err)
	msg := err.Error()
	assert.Contains(s.T(), msg, "area must be >= 0")
	assert.Contains(s.T(), msg, "duplicate pond")
}

func (s *PondServiceTestSuite) TestBulkImportValidate_TooManyPonds() {
	svc := s.pondService.(*pondService)
	ponds := make([]dto.BulkImportPondItem, bulkImportMaxPonds+1)
	for i := range ponds {
		ponds[i] = dto.BulkImportPondItem{Name: fmt.Sprintf("P%d", i)}
	}
	err := svc.validateBulkImportRequest(dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{{Name: "Farm A", Ponds: ponds}},
	})
	require.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "too many ponds")
}

// --- BulkImportFarmPond end-to-end (upsert transaction body) ---
// These cover the actual create/update/leave-alone behavior, complementing
// the pure-function validation tests above.

func (s *PondServiceTestSuite) TestBulkImportFarmPond_NewFarmAndPonds() {
	// GIVEN — farm doesn't exist; both ponds are new
	clientId := 1
	farmName := "Farm A"
	req := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: farmName, Ponds: []dto.BulkImportPondItem{
				{Name: "P1", Area: newDecimal("2.5")},
				{Name: "P2"},
			}},
		},
	}
	// Lookup says farm is new.
	s.farmRepo.On("GetByNameAndClientId", farmName, clientId).Return(nil, nil)
	// Create farm sets Id=10.
	s.farmRepo.On("Create", mock.Anything, mock.MatchedBy(func(f *model.Farm) bool {
		return f.Name == farmName && f.ClientId == clientId && f.Status == constants.FarmStatusMaintenance
	})).Return(nil).Run(func(args mock.Arguments) {
		f := args.Get(1).(*model.Farm)
		f.Id = 10
	})
	// Both pond lookups return nil → create path.
	s.pondRepo.On("GetByFarmIdAndName", 10, "P1").Return(nil, nil)
	s.pondRepo.On("GetByFarmIdAndName", 10, "P2").Return(nil, nil)
	s.pondRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.FarmId == 10 && p.Name == "P1" && p.Status == constants.FarmStatusMaintenance && p.Area.Valid
	})).Return(nil)
	s.pondRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.FarmId == 10 && p.Name == "P2" && p.Status == constants.FarmStatusMaintenance && !p.Area.Valid
	})).Return(nil)
	// Sync farm status: both new ponds are maintenance → status stays maintenance, no Update.
	s.expectFarmStatusSyncAfterMutation(10, []*model.Pond{
		{Id: 1, FarmId: 10, Status: constants.FarmStatusMaintenance},
		{Id: 2, FarmId: 10, Status: constants.FarmStatusMaintenance},
	}, constants.FarmStatusMaintenance)

	// WHEN
	resp, err := s.pondService.BulkImportFarmPond(context.Background(), clientId, req)

	// THEN — response reports 1 new farm + 2 new ponds.
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	assert.Equal(s.T(), 1, resp.FarmsCreated)
	assert.Equal(s.T(), 0, resp.FarmsExisting)
	assert.Equal(s.T(), 2, resp.PondsCreated)
	assert.Equal(s.T(), 0, resp.PondsUpdated)
	assert.Equal(s.T(), 0, resp.PondsUnchanged)
	require.Len(s.T(), resp.Farms, 1)
	assert.Equal(s.T(), farmName, resp.Farms[0].Name)
	assert.True(s.T(), resp.Farms[0].IsNew)
	assert.Equal(s.T(), 2, resp.Farms[0].PondsCreated)
	s.farmRepo.AssertExpectations(s.T())
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestBulkImportFarmPond_ExistingFarmMixedPonds() {
	// GIVEN — farm exists; one new pond, one existing-with-area-update, one
	// existing-without-area (unchanged).
	clientId := 1
	farmId := 7
	farmName := "Farm A"
	req := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: farmName, Ponds: []dto.BulkImportPondItem{
				{Name: "PNew"}, // create
				{Name: "PUpdate", Area: newDecimal("3.0")}, // update area
				{Name: "PUnchanged"},                       // matched, no area → no write
			}},
		},
	}
	existingFarm := &model.Farm{Id: farmId, ClientId: clientId, Name: farmName, Status: constants.FarmStatusMaintenance}
	s.farmRepo.On("GetByNameAndClientId", farmName, clientId).Return(existingFarm, nil)
	// PNew → not found, create.
	s.pondRepo.On("GetByFarmIdAndName", farmId, "PNew").Return(nil, nil)
	s.pondRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.FarmId == farmId && p.Name == "PNew"
	})).Return(nil)
	// PUpdate → found, Update called with new area.
	existingUpdate := &model.Pond{Id: 100, FarmId: farmId, Name: "PUpdate", Status: constants.FarmStatusActive}
	s.pondRepo.On("GetByFarmIdAndName", farmId, "PUpdate").Return(existingUpdate, nil)
	s.pondRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.Id == 100 && p.Area.Valid && p.Area.Decimal.String() == "3"
	})).Return(nil)
	// PUnchanged → found, but no area provided → no Update call.
	existingUnchanged := &model.Pond{Id: 101, FarmId: farmId, Name: "PUnchanged", Status: constants.FarmStatusMaintenance}
	s.pondRepo.On("GetByFarmIdAndName", farmId, "PUnchanged").Return(existingUnchanged, nil)
	// Status sync: one pond is active → farm becomes active. Helper registers Update only if status differs.
	s.expectFarmStatusSyncAfterMutation(farmId, []*model.Pond{
		{Id: 100, FarmId: farmId, Status: constants.FarmStatusActive},
		{Id: 101, FarmId: farmId, Status: constants.FarmStatusMaintenance},
	}, constants.FarmStatusMaintenance)

	// WHEN
	resp, err := s.pondService.BulkImportFarmPond(context.Background(), clientId, req)

	// THEN — counts split correctly across created/updated/unchanged.
	require.NoError(s.T(), err)
	require.NotNil(s.T(), resp)
	assert.Equal(s.T(), 0, resp.FarmsCreated)
	assert.Equal(s.T(), 1, resp.FarmsExisting)
	assert.Equal(s.T(), 1, resp.PondsCreated)
	assert.Equal(s.T(), 1, resp.PondsUpdated)
	assert.Equal(s.T(), 1, resp.PondsUnchanged)
	require.Len(s.T(), resp.Farms, 1)
	assert.False(s.T(), resp.Farms[0].IsNew)
	assert.Equal(s.T(), 1, resp.Farms[0].PondsCreated)
	assert.Equal(s.T(), 1, resp.Farms[0].PondsUpdated)
	assert.Equal(s.T(), 1, resp.Farms[0].PondsUnchanged)
	s.farmRepo.AssertExpectations(s.T())
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestBulkImportFarmPond_NoDeletesForOmittedPond() {
	// GIVEN — farm exists with two ponds in DB, file references only one of
	// them. The omitted pond must NOT be deleted (no-delete contract).
	clientId := 1
	farmId := 7
	farmName := "Farm A"
	req := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: farmName, Ponds: []dto.BulkImportPondItem{{Name: "P1"}}},
		},
	}
	s.farmRepo.On("GetByNameAndClientId", farmName, clientId).Return(
		&model.Farm{Id: farmId, ClientId: clientId, Name: farmName, Status: constants.FarmStatusMaintenance}, nil)
	s.pondRepo.On("GetByFarmIdAndName", farmId, "P1").Return(
		&model.Pond{Id: 200, FarmId: farmId, Name: "P1", Status: constants.FarmStatusMaintenance}, nil)
	// Sync sees both ponds — the one we mentioned AND the omitted one.
	s.expectFarmStatusSyncAfterMutation(farmId, []*model.Pond{
		{Id: 200, FarmId: farmId, Status: constants.FarmStatusMaintenance},
		{Id: 201, FarmId: farmId, Name: "POmitted", Status: constants.FarmStatusMaintenance},
	}, constants.FarmStatusMaintenance)

	// WHEN
	_, err := s.pondService.BulkImportFarmPond(context.Background(), clientId, req)

	// THEN — no error and Delete was never called on any pond.
	require.NoError(s.T(), err)
	s.pondRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
	s.farmRepo.AssertNotCalled(s.T(), "Delete", mock.Anything, mock.Anything)
}

func (s *PondServiceTestSuite) TestBulkImportFarmPond_ValidationErrorShortCircuits() {
	// GIVEN — request fails validation (negative area). No repo work should
	// happen — even the farm lookup must not be called.
	req := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: "Farm A", Ponds: []dto.BulkImportPondItem{
				{Name: "P1", Area: newDecimal("-1")},
			}},
		},
	}

	// WHEN
	resp, err := s.pondService.BulkImportFarmPond(context.Background(), 1, req)

	// THEN — validation error returned; nothing touched the repos.
	require.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), "area must be >= 0")
	s.farmRepo.AssertNotCalled(s.T(), "GetByNameAndClientId", mock.Anything, mock.Anything)
	s.farmRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
	s.pondRepo.AssertNotCalled(s.T(), "Create", mock.Anything, mock.Anything)
}

func (s *PondServiceTestSuite) TestBulkImportFarmPond_RollbackOnPondCreateFailure() {
	// GIVEN — farm gets created, first pond create succeeds, second pond
	// create fails. The transaction must roll back; the response must be nil.
	// We don't assert specific DB state (sqlite is in-memory), but we do
	// assert the service surfaces the error.
	clientId := 1
	farmName := "Farm A"
	req := dto.BulkImportFarmPondRequest{
		Farms: []dto.BulkImportFarmItem{
			{Name: farmName, Ponds: []dto.BulkImportPondItem{
				{Name: "P1"},
				{Name: "P2"},
			}},
		},
	}
	s.farmRepo.On("GetByNameAndClientId", farmName, clientId).Return(nil, nil)
	s.farmRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		args.Get(1).(*model.Farm).Id = 10
	})
	s.pondRepo.On("GetByFarmIdAndName", 10, "P1").Return(nil, nil)
	s.pondRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.Name == "P1"
	})).Return(nil)
	s.pondRepo.On("GetByFarmIdAndName", 10, "P2").Return(nil, nil)
	s.pondRepo.On("Create", mock.Anything, mock.MatchedBy(func(p *model.Pond) bool {
		return p.Name == "P2"
	})).Return(fmt.Errorf("simulated db error"))

	// WHEN
	resp, err := s.pondService.BulkImportFarmPond(context.Background(), clientId, req)

	// THEN — error surfaced; nil response; sync should NOT have run since
	// the transaction was aborted before the sync loop.
	require.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	s.pondRepo.AssertNotCalled(s.T(), "ListByFarmId", mock.Anything)
}
