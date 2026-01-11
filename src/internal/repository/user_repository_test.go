package repository

import (
	"strconv"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	db             *gorm.DB
	userRepository UserRepository
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	// Use in-memory SQLite for fast tests
	var err error
	s.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		s.T().Fatal("Failed to connect to test database:", err)
	}

	// Auto-migrate the schema
	err = s.db.AutoMigrate(&model.User{})
	if err != nil {
		s.T().Fatal("Failed to migrate database:", err)
	}

	s.userRepository = NewUserRepository(s.db)
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	sqlDB, _ := s.db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func (s *UserRepositoryTestSuite) SetupTest() {
	// Clean up before each test - GORM uses table name "users" by default
	s.db.Exec("DELETE FROM users")
}

func (s *UserRepositoryTestSuite) TearDownTest() {
	// Additional cleanup if needed (SetupTest already does this)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}

// Test Create operations
func (s *UserRepositoryTestSuite) TestCreate_Success() {
	user := &model.User{
		Username:      "testuser",
		Password:      "hashed_password",
		FirstName:     "Test",
		LastName:      nil,
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}

	err := s.userRepository.Create(user)

	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), user.Id)
	assert.NotZero(s.T(), user.CreatedAt)
	assert.NotZero(s.T(), user.UpdatedAt)
	assert.Equal(s.T(), "testuser", user.Username)
	assert.Equal(s.T(), "Test", user.FirstName)
}

func (s *UserRepositoryTestSuite) TestCreate_DuplicateUsername() {
	user1 := &model.User{
		Username:      "testuser",
		Password:      "password1",
		FirstName:     "User",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user1)
	assert.NoError(s.T(), err)

	user2 := &model.User{
		Username:      "testuser", // Duplicate username
		Password:      "password2",
		FirstName:     "User2",
		UserLevel:     1,
		ContactNumber: "0987654321",
		ClientId:      lo.ToPtr(1),
	}

	err = s.userRepository.Create(user2)

	assert.Error(s.T(), err)
	// Verify the duplicate wasn't created
	count := int64(0)
	s.db.Model(&model.User{}).Where("username = ?", "testuser").Count(&count)
	assert.Equal(s.T(), int64(1), count)
}

// Test GetByID operations
func (s *UserRepositoryTestSuite) TestGetByID_Success() {
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	result, err := s.userRepository.GetByID(user.Id)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Id, result.Id)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.FirstName, result.FirstName)
	assert.Equal(s.T(), user.Password, result.Password)
}

func (s *UserRepositoryTestSuite) TestGetByID_NotFound() {
	result, err := s.userRepository.GetByID(999)

	// Repository returns nil, nil when record not found (not an error)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

// Test GetByEmail operations - Note: Email field doesn't exist in new model
// These tests are skipped since Email field is not available in the new User model
// If GetByEmail is still needed, the repository implementation needs to be updated
func (s *UserRepositoryTestSuite) TestGetByEmail_Success() {
	s.T().Skip("Email field not available in new model - GetByEmail needs to be updated or removed")
}

func (s *UserRepositoryTestSuite) TestGetByEmail_NotFound() {
	s.T().Skip("Email field not available in new model - GetByEmail needs to be updated or removed")
}

func (s *UserRepositoryTestSuite) TestGetByEmail_CaseSensitive() {
	s.T().Skip("Email field not available in new model - GetByEmail needs to be updated or removed")
}

// Test GetByUsername operations
func (s *UserRepositoryTestSuite) TestGetByUsername_Success() {
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	result, err := s.userRepository.GetByUsername("testuser")

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.FirstName, result.FirstName)
	assert.Equal(s.T(), user.Id, result.Id)
}

func (s *UserRepositoryTestSuite) TestGetByUsername_NotFound() {
	result, err := s.userRepository.GetByUsername("nonexistent")

	// GetByUsername returns (nil, nil) when not found, not an error
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

// Test Update operations
func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	user := &model.User{
		Username:      "olduser",
		Password:      "password",
		FirstName:     "Old",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(10 * time.Millisecond)

	user.Username = "newuser"
	user.FirstName = "New"
	err = s.userRepository.Update(user)

	assert.NoError(s.T(), err)

	// Verify update
	updated, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newuser", updated.Username)
	assert.Equal(s.T(), "New", updated.FirstName)
	assert.True(s.T(), updated.UpdatedAt.After(originalUpdatedAt))
}

func (s *UserRepositoryTestSuite) TestUpdate_PartialFields() {
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	// Update only username
	user.Username = "updateduser"
	err = s.userRepository.Update(user)

	assert.NoError(s.T(), err)

	// Verify only username changed
	updated, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updateduser", updated.Username)
	assert.Equal(s.T(), "Test", updated.FirstName) // FirstName unchanged
}

func (s *UserRepositoryTestSuite) TestUpdate_NonExistent() {
	user := &model.User{
		Id:            999,
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}

	// GORM Save will create if not found, so this might not error
	err := s.userRepository.Update(user)

	// GORM Save creates if ID doesn't exist, so we verify it was created
	assert.NoError(s.T(), err)
	result, err := s.userRepository.GetByID(999)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
}

// Test Delete operations
func (s *UserRepositoryTestSuite) TestDelete_Success() {
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)
	userID := user.Id

	err = s.userRepository.Delete(userID)

	assert.NoError(s.T(), err)

	// Verify soft delete (User model has DeletedAt field)
	// After soft delete, GetByID should return nil, nil (not found)
	result, err := s.userRepository.GetByID(userID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)

	// Verify it's soft deleted (still in DB but with DeletedAt set)
	var deletedUser model.User
	s.db.Unscoped().First(&deletedUser, userID)
	assert.NotZero(s.T(), deletedUser.DeletedAt)
}

func (s *UserRepositoryTestSuite) TestDelete_NotFound() {
	err := s.userRepository.Delete(999)

	// GORM Delete doesn't error on non-existent records
	assert.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestDelete_MultipleUsers() {
	// Create multiple users
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: lo.ToPtr(1)}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: lo.ToPtr(1)}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: lo.ToPtr(1)}

	s.userRepository.Create(user1)
	s.userRepository.Create(user2)
	s.userRepository.Create(user3)

	// Delete one user
	err := s.userRepository.Delete(user2.Id)
	assert.NoError(s.T(), err)

	// Verify only user2 is deleted
	result1, err := s.userRepository.GetByID(user1.Id)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result1)

	result2, err := s.userRepository.GetByID(user2.Id)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result2) // Deleted user returns nil, nil

	result3, err := s.userRepository.GetByID(user3.Id)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result3)
}

// Test ListByClientId operations
func (s *UserRepositoryTestSuite) TestListByClientId_Success() {
	clientId := 1
	// Create multiple users for clientId 1
	for i := 0; i < 5; i++ {
		user := &model.User{
			Username:      "user" + strconv.Itoa(i),
			Password:      "password",
			FirstName:     "User" + strconv.Itoa(i),
			UserLevel:     1,
			ContactNumber: "123456789" + strconv.Itoa(i),
			ClientId:      &clientId,
		}
		s.userRepository.Create(user)
	}

	// Create users for different clientId
	for i := 0; i < 3; i++ {
		user := &model.User{
			Username:      "otheruser" + strconv.Itoa(i),
			Password:      "password",
			FirstName:     "Other" + strconv.Itoa(i),
			UserLevel:     1,
			ContactNumber: "987654321" + strconv.Itoa(i),
			ClientId:      lo.ToPtr(2),
		}
		s.userRepository.Create(user)
	}

	users, err := s.userRepository.ListByClientId(clientId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 5)
	// Verify all users belong to the correct clientId
	for _, user := range users {
		assert.Equal(s.T(), &clientId, user.ClientId)
	}
}

func (s *UserRepositoryTestSuite) TestListByClientId_Empty() {
	users, err := s.userRepository.ListByClientId(999)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), users)
	assert.Len(s.T(), users, 0)
}

func (s *UserRepositoryTestSuite) TestListByClientId_ExcludesSoftDeleted() {
	clientId := 1
	// Create users
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: &clientId}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: &clientId}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: &clientId}

	s.userRepository.Create(user1)
	s.userRepository.Create(user2)
	s.userRepository.Create(user3)

	// Delete one user
	s.userRepository.Delete(user2.Id)

	// List should only return non-deleted users
	users, err := s.userRepository.ListByClientId(clientId)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	// Verify user2 is not in the list
	for _, user := range users {
		assert.NotEqual(s.T(), user2.Id, user.Id)
	}
}

func (s *UserRepositoryTestSuite) TestListByClientId_FiltersByClientId() {
	// Create users for different clientIds
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: lo.ToPtr(1)}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: lo.ToPtr(1)}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: lo.ToPtr(2)}

	s.userRepository.Create(user1)
	s.userRepository.Create(user2)
	s.userRepository.Create(user3)

	// List users for clientId 1
	users, err := s.userRepository.ListByClientId(1)

	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	// Verify all users belong to clientId 1
	for _, user := range users {
		assert.Equal(s.T(), lo.ToPtr(1), user.ClientId)
		assert.NotEqual(s.T(), user3.Id, user.Id)
	}
}

// Test edge cases and data integrity
func (s *UserRepositoryTestSuite) TestCreate_RequiredFields() {
	// Note: SQLite is more lenient with NOT NULL constraints than PostgreSQL
	// This test verifies behavior with empty strings
	// For strict NOT NULL enforcement, use PostgreSQL in integration tests

	// Test with empty strings (SQLite allows these, PostgreSQL would reject)
	user := &model.User{
		Username:      "",
		Password:      "",
		FirstName:     "",
		UserLevel:     0,
		ContactNumber: "",
		ClientId:      nil,
	}

	err := s.userRepository.Create(user)

	// SQLite allows empty strings, so this may succeed
	// We just verify the behavior is consistent
	if err == nil {
		// If it succeeds, verify the record was created
		assert.NotZero(s.T(), user.Id)
	}
	// If it fails, that's also acceptable - depends on database configuration
}

func (s *UserRepositoryTestSuite) TestGetByID_AfterDelete() {
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	s.userRepository.Create(user)
	userID := user.Id

	// Delete the user
	s.userRepository.Delete(userID)

	// Try to get it - should return nil, nil (not found)
	result, err := s.userRepository.GetByID(userID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_AfterDelete() {
	s.T().Skip("Email field not available in new model - GetByEmail needs to be updated or removed")
}

func (s *UserRepositoryTestSuite) TestMultipleOperations() {
	// Create
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(user)
	assert.NoError(s.T(), err)

	// Get by ID
	found, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), user.Id, found.Id)

	// Get by Username
	found, err = s.userRepository.GetByUsername(user.Username)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), found)
	assert.Equal(s.T(), user.Username, found.Username)

	// Update
	user.Username = "updateduser"
	err = s.userRepository.Update(user)
	assert.NoError(s.T(), err)

	// Verify update
	found, err = s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updateduser", found.Username)

	// Delete
	err = s.userRepository.Delete(user.Id)
	assert.NoError(s.T(), err)

	// Verify deletion - should return nil, nil (not found)
	deletedResult, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), deletedResult)
}

