package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type FarmRepositoryTestSuite struct {
	suite.Suite
	db           *gorm.DB
	farmRepo     FarmRepository
}

func (s *FarmRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.Farm{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.farmRepo = NewFarmRepository(s.db)
}

func (s *FarmRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *FarmRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM farms")
}

func TestFarmRepositorySuite(t *testing.T) {
	suite.Run(t, new(FarmRepositoryTestSuite))
}

func (s *FarmRepositoryTestSuite) TestCreate_Success() {
	farm := &model.Farm{
		ClientId: 1,
		Code:     "FARM001",
		Name:     "Test Farm",
	}

	err := s.farmRepo.Create(farm)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), farm.Id)
	assert.Equal(s.T(), "FARM001", farm.Code)
	assert.Equal(s.T(), "Test Farm", farm.Name)
}

func (s *FarmRepositoryTestSuite) TestGetByID_Success() {
	farm := &model.Farm{
		ClientId: 1,
		Code:     "FARM001",
		Name:     "Test Farm",
	}
	s.farmRepo.Create(farm)

	result, err := s.farmRepo.GetByID(farm.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), farm.Id, result.Id)
	assert.Equal(s.T(), "FARM001", result.Code)
}

func (s *FarmRepositoryTestSuite) TestGetByID_NotFound() {
	result, err := s.farmRepo.GetByID(999)

	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *FarmRepositoryTestSuite) TestGetByCodeAndClientId_Success() {
	farm := &model.Farm{
		ClientId: 1,
		Code:     "FARM001",
		Name:     "Test Farm",
	}
	s.farmRepo.Create(farm)

	result, err := s.farmRepo.GetByCodeAndClientId("FARM001", 1)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "FARM001", result.Code)
}

func (s *FarmRepositoryTestSuite) TestListByClientId_Success() {
	farm1 := &model.Farm{ClientId: 1, Code: "FARM001", Name: "Farm 1"}
	farm2 := &model.Farm{ClientId: 1, Code: "FARM002", Name: "Farm 2"}
	farm3 := &model.Farm{ClientId: 2, Code: "FARM003", Name: "Farm 3"}
	s.farmRepo.Create(farm1)
	s.farmRepo.Create(farm2)
	s.farmRepo.Create(farm3)

	results, err := s.farmRepo.ListByClientId(1)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2)
}

func (s *FarmRepositoryTestSuite) TestUpdate_Success() {
	farm := &model.Farm{
		ClientId: 1,
		Code:     "FARM001",
		Name:     "Test Farm",
	}
	s.farmRepo.Create(farm)

	farm.Name = "Updated Farm"
	err := s.farmRepo.Update(farm)

	assert.NoError(s.T(), err)
	
	updated, _ := s.farmRepo.GetByID(farm.Id)
	assert.Equal(s.T(), "Updated Farm", updated.Name)
}

