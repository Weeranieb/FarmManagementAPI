package service

import (
	"context"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/constants"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/errors"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
	"github.com/weeranieb/boonmafarm-backend/src/internal/transaction"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PondServiceTestSuite struct {
	suite.Suite
	pondRepo       *mocks.MockPondRepository
	farmRepo       *mocks.MockFarmRepository
	activePondRepo *mocks.MockActivePondRepository
	activityRepo   *mocks.MockActivityRepository
	db             *gorm.DB
	pondService    PondService
}

func (s *PondServiceTestSuite) SetupTest() {
	s.pondRepo = mocks.NewMockPondRepository(s.T())
	s.farmRepo = mocks.NewMockFarmRepository(s.T())
	s.activePondRepo = mocks.NewMockActivePondRepository(s.T())
	s.activityRepo = mocks.NewMockActivityRepository(s.T())
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	s.Require().NoError(err)
	err = s.db.AutoMigrate(&model.Pond{}, &model.ActivePond{}, &model.Activity{}, &model.AdditionalCost{})
	s.Require().NoError(err)
	s.pondService = NewPondService(PondServiceParams{
		PondRepo:       s.pondRepo,
		FarmRepo:       s.farmRepo,
		ActivePondRepo: s.activePondRepo,
		ActivityRepo:   s.activityRepo,
		TxManager:      transaction.NewManager(s.db),
	})
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
}

func TestPondServiceSuite(t *testing.T) {
	suite.Run(t, new(PondServiceTestSuite))
}

func (s *PondServiceTestSuite) TestCreatePonds_Success() {
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}
	username := "admin"

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

	err := s.pondService.CreatePonds(context.Background(), req, username)

	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestCreatePonds_PondAlreadyExists() {
	req := dto.CreatePondsRequest{
		FarmId: 1,
		Names:  []string{"Pond 1", "Pond 2"},
	}
	username := "admin"

	// First name is free, second name already exists for this farm
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 1").Return(nil, nil)
	existingPond := &model.Pond{Id: 99, FarmId: 1, Name: "Pond 2", Status: "active"}
	s.pondRepo.On("GetByFarmIdAndName", 1, "Pond 2").Return(existingPond, nil)

	err := s.pondService.CreatePonds(context.Background(), req, username)

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
	// CreateBatch must not be called
	s.pondRepo.AssertNotCalled(s.T(), "CreateBatch")
}

func (s *PondServiceTestSuite) TestGet_Success() {
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

	result, err := s.pondService.Get(pondId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), pondId, result.Id)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGet_NotFound() {
	pondId := 999
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, nil)

	result, err := s.pondService.Get(pondId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGet_RepoError() {
	pondId := 1
	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, assert.AnError)

	result, err := s.pondService.Get(pondId)

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestGetList_Success() {
	farmId := 1
	list := []*repository.PondWithFarmAndActivePond{
		{Pond: &model.Pond{Id: 1, FarmId: farmId, Name: "Pond 1", Status: "active"}, ClientId: 1, ActivePond: nil},
		{Pond: &model.Pond{Id: 2, FarmId: farmId, Name: "Pond 2", Status: "active"}, ClientId: 1, ActivePond: nil},
	}
	s.pondRepo.On("ListByFarmIdWithActivePond", mock.Anything, farmId).Return(list, nil)

	result, err := s.pondService.GetList(farmId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestUpdate_Success() {
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Old Name", Status: "maintenance"}
	req := dto.UpdatePondRequest{Id: 1, Name: "New Name", Status: "active"}

	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "New Name").Return(nil, nil) // normalized name, no duplicate
	s.pondRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pond")).Return(nil)

	err := s.pondService.Update(context.Background(), req, "user")

	assert.NoError(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestUpdate_PondNotFound() {
	req := dto.UpdatePondRequest{Id: 999, Name: "Pond"}

	s.pondRepo.On("GetByID", 999).Return(nil, nil)

	err := s.pondService.Update(context.Background(), req, "user")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
	s.pondRepo.AssertNotCalled(s.T(), "Update")
}

func (s *PondServiceTestSuite) TestUpdate_DuplicateName() {
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Old", Status: "active"}
	otherPond := &model.Pond{Id: 2, FarmId: 1, Name: "New Name", Status: "active"}
	req := dto.UpdatePondRequest{Id: 1, Name: "New Name"}

	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("GetByFarmIdAndName", 1, "New Name").Return(otherPond, nil)

	err := s.pondService.Update(context.Background(), req, "user")

	assert.Error(s.T(), err)
	assert.ErrorIs(s.T(), err, errors.ErrPondAlreadyExists)
	s.pondRepo.AssertExpectations(s.T())
	s.pondRepo.AssertNotCalled(s.T(), "Update")
}

func (s *PondServiceTestSuite) TestUpdate_RepoError() {
	existing := &model.Pond{Id: 1, FarmId: 1, Name: "Pond", Status: "active"}
	req := dto.UpdatePondRequest{Id: 1, Status: "maintenance"}

	s.pondRepo.On("GetByID", 1).Return(existing, nil)
	s.pondRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.Pond")).Return(assert.AnError)

	err := s.pondService.Update(context.Background(), req, "user")

	assert.Error(s.T(), err)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_PondNotFound() {
	pondId := 999
	req := validPondFillRequest()

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, nil)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrPondNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_RepoError() {
	pondId := 1
	req := validPondFillRequest()

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(nil, assert.AnError)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_FarmNotFound() {
	pondId := 1
	req := validPondFillRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   0, // no farm
		ActivePond: nil,
	}

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrFarmNotFound)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_PermissionDenied() {
	pondId := 1
	req := validPondFillRequest()
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   2, // user is client 1, cannot access client 2
		ActivePond: nil,
	}

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	resp, err := s.pondService.FillPond(fillPondCtxNoAccess(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrAuthPermissionDenied)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_InvalidFishType() {
	pondId := 1
	req := validPondFillRequest()
	req.FishType = "invalid"
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   1,
		ActivePond: nil,
	}

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.ErrorIs(s.T(), err, errors.ErrInvalidFishType)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_InvalidActivityDate() {
	pondId := 1
	req := validPondFillRequest()
	req.ActivityDate = "not-a-date"
	data := &repository.PondWithFarmAndActivePond{
		Pond:       &model.Pond{Id: pondId, FarmId: 1, Name: "P", Status: "active"},
		ClientId:   1,
		ActivePond: nil,
	}

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.Error(s.T(), err)
	assert.Nil(s.T(), resp)
	assert.Contains(s.T(), err.Error(), errors.ErrValidationFailed.Message)
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_Success_NewActivePond() {
	pondId := 1
	req := validPondFillRequest()
	pond := &model.Pond{Id: pondId, FarmId: 1, Name: "Pond", Status: constants.FarmStatusMaintenance}
	data := &repository.PondWithFarmAndActivePond{
		Pond:       pond,
		ClientId:   1,
		ActivePond: nil,
	}

	s.pondRepo.On("GetByIDWithFarmAndActivePond", mock.Anything, pondId).Return(data, nil)

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Greater(s.T(), resp.ActivePondId, int64(0))
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	s.pondRepo.AssertExpectations(s.T())
}

func (s *PondServiceTestSuite) TestFillPond_Success_ExistingActivePond() {
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

	resp, err := s.pondService.FillPond(fillPondCtx(), pondId, req, "user")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), resp)
	assert.Equal(s.T(), int64(10), resp.ActivePondId)
	assert.Greater(s.T(), resp.ActivityId, int64(0))
	s.pondRepo.AssertExpectations(s.T())
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
