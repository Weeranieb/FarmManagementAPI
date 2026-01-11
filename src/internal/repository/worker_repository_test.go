package repository

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type WorkerRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	workerRepo WorkerRepository
}

func (s *WorkerRepositoryTestSuite) SetupSuite() {
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	err = s.db.AutoMigrate(&model.Worker{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.workerRepo = NewWorkerRepository(s.db)
}

func (s *WorkerRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *WorkerRepositoryTestSuite) SetupTest() {
	s.db.Exec("DELETE FROM workers")
}

func TestWorkerRepositorySuite(t *testing.T) {
	suite.Run(t, new(WorkerRepositoryTestSuite))
}

func (s *WorkerRepositoryTestSuite) TestCreate_Success() {
	worker := &model.Worker{
		ClientId:    1,
		FarmGroupId: 1,
		FirstName:   "John",
		LastName:    lo.ToPtr("Doe"),
		Nationality: "Thai",
		Salary:      50000,
		IsActive:    true,
	}

	err := s.workerRepo.Create(worker)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), worker.Id)
	assert.Equal(s.T(), "John", worker.FirstName)
}

func (s *WorkerRepositoryTestSuite) TestGetByID_Success() {
	worker := &model.Worker{
		ClientId:    1,
		FarmGroupId: 1,
		FirstName:   "John",
		Nationality: "Thai",
		Salary:      50000,
		IsActive:    true,
	}
	s.workerRepo.Create(worker)

	result, err := s.workerRepo.GetByID(worker.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), worker.Id, result.Id)
}

func (s *WorkerRepositoryTestSuite) TestGetPage_Success() {
	worker1 := &model.Worker{ClientId: 1, FarmGroupId: 1, FirstName: "John", Nationality: "Thai", Salary: 50000, IsActive: true}
	worker2 := &model.Worker{ClientId: 1, FarmGroupId: 1, FirstName: "Jane", Nationality: "Thai", Salary: 60000, IsActive: true}
	worker3 := &model.Worker{ClientId: 2, FarmGroupId: 1, FirstName: "Bob", Nationality: "Thai", Salary: 55000, IsActive: true}
	s.workerRepo.Create(worker1)
	s.workerRepo.Create(worker2)
	s.workerRepo.Create(worker3)

	results, total, err := s.workerRepo.GetPage(1, 0, 10, "", "")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(2), total)
	assert.Len(s.T(), results, 2)
}

func (s *WorkerRepositoryTestSuite) TestGetPage_WithKeyword() {
	worker1 := &model.Worker{ClientId: 1, FarmGroupId: 1, FirstName: "John", Nationality: "Thai", Salary: 50000, IsActive: true}
	worker2 := &model.Worker{ClientId: 1, FarmGroupId: 1, FirstName: "Jane", Nationality: "Thai", Salary: 60000, IsActive: true}
	s.workerRepo.Create(worker1)
	s.workerRepo.Create(worker2)

	results, total, err := s.workerRepo.GetPage(1, 0, 10, "", "John")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), total)
	assert.Len(s.T(), results, 1)
	assert.Equal(s.T(), "John", results[0].FirstName)
}

func (s *WorkerRepositoryTestSuite) TestUpdate_Success() {
	worker := &model.Worker{
		ClientId:    1,
		FarmGroupId: 1,
		FirstName:   "John",
		Nationality: "Thai",
		Salary:      50000,
		IsActive:    true,
	}
	s.workerRepo.Create(worker)

	worker.FirstName = "Johnny"
	err := s.workerRepo.Update(worker)

	assert.NoError(s.T(), err)
	
	updated, _ := s.workerRepo.GetByID(worker.Id)
	assert.Equal(s.T(), "Johnny", updated.FirstName)
}


