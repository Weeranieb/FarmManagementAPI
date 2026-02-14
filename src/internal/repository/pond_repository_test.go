package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type PondRepositoryTestSuite struct {
	suite.Suite
	db       *gorm.DB
	pondRepo PondRepository
}

func (s *PondRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.Pond{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.pondRepo = NewPondRepository(s.db)
}

func (s *PondRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *PondRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM ponds")
}

func TestPondRepositorySuite(t *testing.T) {
	suite.Run(t, new(PondRepositoryTestSuite))
}

func (s *PondRepositoryTestSuite) TestCreate_Success() {
	pond := &model.Pond{
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}

	err := s.pondRepo.Create(context.Background(), pond)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), pond.Id)
	assert.Equal(s.T(), "Test Pond", pond.Name)
}

func (s *PondRepositoryTestSuite) TestCreateBatch_Success() {
	ponds := []*model.Pond{
		{FarmId: 1, Name: "Pond 1", Status: "active"},
		{FarmId: 1, Name: "Pond 2", Status: "active"},
	}

	err := s.pondRepo.CreateBatch(context.Background(), ponds)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), ponds[0].Id)
	assert.NotZero(s.T(), ponds[1].Id)
}

func (s *PondRepositoryTestSuite) TestGetByID_Success() {
	pond := &model.Pond{
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}
	s.pondRepo.Create(context.Background(), pond)

	result, err := s.pondRepo.GetByID(pond.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), pond.Id, result.Id)
}

func (s *PondRepositoryTestSuite) TestGetByFarmIdAndName_Success() {
	pond := &model.Pond{
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}
	s.pondRepo.Create(context.Background(), pond)

	result, err := s.pondRepo.GetByFarmIdAndName(1, "Test Pond")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), "Test Pond", result.Name)
}

func (s *PondRepositoryTestSuite) TestListByFarmId_Success() {
	pond1 := &model.Pond{FarmId: 1, Name: "Pond 1", Status: "active"}
	pond2 := &model.Pond{FarmId: 1, Name: "Pond 2", Status: "active"}
	pond3 := &model.Pond{FarmId: 2, Name: "Pond 3", Status: "active"}
	s.pondRepo.Create(context.Background(), pond1)
	s.pondRepo.Create(context.Background(), pond2)
	s.pondRepo.Create(context.Background(), pond3)

	results, err := s.pondRepo.ListByFarmId(1)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), results, 2)
}

func (s *PondRepositoryTestSuite) TestUpdate_Success() {
	pond := &model.Pond{
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}
	s.pondRepo.Create(context.Background(), pond)

	pond.Name = "Updated Pond"
	err := s.pondRepo.Update(context.Background(), pond)

	assert.NoError(s.T(), err)

	updated, _ := s.pondRepo.GetByID(pond.Id)
	assert.Equal(s.T(), "Updated Pond", updated.Name)
}

func (s *PondRepositoryTestSuite) TestDelete_Success() {
	pond := &model.Pond{
		FarmId: 1,
		Name:   "Test Pond",
		Status: "active",
	}
	s.pondRepo.Create(context.Background(), pond)

	err := s.pondRepo.Delete(pond.Id)

	assert.NoError(s.T(), err)

	deleted, _ := s.pondRepo.GetByID(pond.Id)
	assert.Nil(s.T(), deleted)
}
