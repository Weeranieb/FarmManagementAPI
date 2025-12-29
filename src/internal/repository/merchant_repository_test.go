package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MerchantRepositoryTestSuite struct {
	suite.Suite
	db            *gorm.DB
	merchantRepo  MerchantRepository
}

func (s *MerchantRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.Merchant{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.merchantRepo = NewMerchantRepository(s.db)
}

func (s *MerchantRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *MerchantRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM merchants")
}

func TestMerchantRepositorySuite(t *testing.T) {
	suite.Run(t, new(MerchantRepositoryTestSuite))
}

func (s *MerchantRepositoryTestSuite) TestCreate_Success() {
	merchant := &model.Merchant{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}

	err := s.merchantRepo.Create(merchant)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), merchant.Id)
	assert.Equal(s.T(), "Test Merchant", merchant.Name)
}

func (s *MerchantRepositoryTestSuite) TestGetByID_Success() {
	merchant := &model.Merchant{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}
	s.merchantRepo.Create(merchant)

	result, err := s.merchantRepo.GetByID(merchant.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), merchant.Id, result.Id)
}

func (s *MerchantRepositoryTestSuite) TestGetByID_NotFound() {
	result, err := s.merchantRepo.GetByID(999)

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *MerchantRepositoryTestSuite) TestGetByContactNumberAndName_Success() {
	merchant := &model.Merchant{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}
	s.merchantRepo.Create(merchant)

	result, err := s.merchantRepo.GetByContactNumberAndName("1234567890", "Test Merchant")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "1234567890", result.ContactNumber)
}

func (s *MerchantRepositoryTestSuite) TestList_Success() {
	merchant1 := &model.Merchant{Name: "Merchant 1", ContactNumber: "111", Location: "Loc 1"}
	merchant2 := &model.Merchant{Name: "Merchant 2", ContactNumber: "222", Location: "Loc 2"}
	s.merchantRepo.Create(merchant1)
	s.merchantRepo.Create(merchant2)

	results, err := s.merchantRepo.List()

	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2)
}

func (s *MerchantRepositoryTestSuite) TestUpdate_Success() {
	merchant := &model.Merchant{
		Name:          "Test Merchant",
		ContactNumber: "1234567890",
		Location:      "Test Location",
	}
	s.merchantRepo.Create(merchant)

	merchant.Name = "Updated Merchant"
	err := s.merchantRepo.Update(merchant)

	assert.NoError(s.T(), err)
	
	updated, _ := s.merchantRepo.GetByID(merchant.Id)
	assert.Equal(s.T(), "Updated Merchant", updated.Name)
}

