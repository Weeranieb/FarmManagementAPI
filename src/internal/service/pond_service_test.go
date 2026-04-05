package service

import (
	"context"
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
		TxManager:          transaction.NewManager(s.db),
	})
	s.pondRepo.On("WithTx", mock.Anything).Maybe().Return(s.pondRepo)
	s.farmRepo.On("WithTx", mock.Anything).Maybe().Return(s.farmRepo)
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
	s.additionalCostRepo.On("Create", mock.Anything, mock.Anything).Maybe().Return(nil)
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
	// GIVEN — request with farm and names; repo returns no duplicate names
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
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

	// WHEN — CreatePonds is called
	err := s.pondService.CreatePonds(context.Background(), req)

	// THEN — no error; CreateBatch was used
	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
	s.farmRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestCreatePonds_PondAlreadyExists() {
	// GIVEN — request; second name already exists for this farm
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	existingPond := &model.Pond{Id: 99, FarmId: 1, Name: "Pond 2", Status: "active"}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(existingPond, nil)

	// WHEN — CreatePonds is called
	err := s.pondService.CreatePonds(context.Background(), req)

	// THEN — ErrPondAlreadyExists; CreateBatch not called
	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
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
