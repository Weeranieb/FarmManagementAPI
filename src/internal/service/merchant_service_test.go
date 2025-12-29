package service

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/dto"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type MerchantServiceTestSuite struct {
	suite.Suite
	merchantRepo   *mocks.MockMerchantRepository
	merchantService MerchantService
}

func (s *MerchantServiceTestSuite) SetupTest() {
	s.merchantRepo = mocks.NewMockMerchantRepository(s.T())
	s.merchantService = NewMerchantService(s.merchantRepo)
}

func (s *MerchantServiceTestSuite) TearDownTest() {
	s.merchantRepo.ExpectedCalls = nil
}

func TestMerchantServiceSuite(t *testing.T) {
	suite.Run(t, new(MerchantServiceTestSuite))
}

func (s *MerchantServiceTestSuite) TestCreate_Success() {
	req := dto.CreateMerchantRequest{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}
	username := "admin"

	s.merchantRepo.On("GetByContactNumberAndName", req.ContactNumber, req.Name).Return(nil, nil)

	expectedTime := time.Now()
	expectedMerchant := &model.Merchant{
		Id:            1,
		Name:          req.Name,
		ContactNumber: req.ContactNumber,
		Location:      req.Location,
		BaseModel: model.BaseModel{
			CreatedAt: expectedTime,
			UpdatedAt: expectedTime,
			CreatedBy: username,
			UpdatedBy: username,
		},
	}

	s.merchantRepo.On("Create", mock.AnythingOfType("*model.Merchant")).Return(nil).Run(func(args mock.Arguments) {
		merchant := args.Get(0).(*model.Merchant)
		merchant.Id = expectedMerchant.Id
		merchant.CreatedAt = expectedMerchant.CreatedAt
		merchant.UpdatedAt = expectedMerchant.UpdatedAt
	})

	result, err := s.merchantService.Create(req, username)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Name, result.Name)
	s.merchantRepo.AssertExpectations(s.T())
}

func (s *MerchantServiceTestSuite) TestGet_Success() {
	merchantId := 1
	expectedMerchant := &model.Merchant{
		Id:            merchantId,
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}

	s.merchantRepo.On("GetByID", merchantId).Return(expectedMerchant, nil)

	result, err := s.merchantService.Get(merchantId)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), merchantId, result.Id)
	s.merchantRepo.AssertExpectations(s.T())
}

func (s *MerchantServiceTestSuite) TestGetList_Success() {
	merchants := []*model.Merchant{
		{Id: 1, Name: "Merchant 1", ContactNumber: "111", Location: "Loc 1"},
		{Id: 2, Name: "Merchant 2", ContactNumber: "222", Location: "Loc 2"},
	}

	s.merchantRepo.On("List").Return(merchants, nil)

	result, err := s.merchantService.GetList()

	assert.NoError(s.T(), err)
	assert.Len(s.T(), result, 2)
	s.merchantRepo.AssertExpectations(s.T())
}

