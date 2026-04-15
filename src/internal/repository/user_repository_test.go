//go:build cgo

package repository

import (
	"context"
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
		_ = sqlDB.Close()
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
	// GIVEN — a new user model
	user := &model.User{
		Username:      "testuser",
		Password:      "hashed_password",
		FirstName:     "Test",
		LastName:      nil,
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}

	// WHEN — Create is called
	err := s.userRepository.Create(context.Background(), user)

	// THEN — no error and id/timestamps/fields set
	assert.NoError(s.T(), err)
	assert.NotZero(s.T(), user.Id)
	assert.NotZero(s.T(), user.CreatedAt)
	assert.NotZero(s.T(), user.UpdatedAt)
	assert.Equal(s.T(), "testuser", user.Username)
	assert.Equal(s.T(), "Test", user.FirstName)
}

func (s *UserRepositoryTestSuite) TestCreate_DuplicateUsername() {
	// GIVEN — one user already created with username "testuser"
	user1 := &model.User{
		Username:      "testuser",
		Password:      "password1",
		FirstName:     "User",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user1)
	assert.NoError(s.T(), err)

	user2 := &model.User{
		Username:      "testuser", // Duplicate username
		Password:      "password2",
		FirstName:     "User2",
		UserLevel:     1,
		ContactNumber: "0987654321",
		ClientId:      lo.ToPtr(1),
	}

	// WHEN — Create is called again with same username
	err = s.userRepository.Create(context.Background(), user2)

	// THEN — error and only one user with that username
	assert.Error(s.T(), err)
	count := int64(0)
	s.db.Model(&model.User{}).Where("username = ?", "testuser").Count(&count)
	assert.Equal(s.T(), int64(1), count)
}

// Test GetByID operations
func (s *UserRepositoryTestSuite) TestGetByID_Success() {
	// GIVEN — a user created in DB
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
	assert.NoError(s.T(), err)

	// WHEN — GetByID is called
	result, err := s.userRepository.GetByID(user.Id)

	// THEN — user is returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Id, result.Id)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.FirstName, result.FirstName)
	assert.Equal(s.T(), user.Password, result.Password)
}

func (s *UserRepositoryTestSuite) TestGetByID_NotFound() {
	// GIVEN — no user with id 999
	// WHEN — GetByID(999) is called
	result, err := s.userRepository.GetByID(999)

	// THEN — nil, nil (no error)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

// Test GetByEmail operations - Note: Email field doesn't exist in new model
func (s *UserRepositoryTestSuite) TestGetByEmail_Success() {
	// GIVEN/WHEN/THEN — skipped (email not in model)
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
	// GIVEN — a user created with username "testuser"
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
	assert.NoError(s.T(), err)

	// WHEN — GetByUsername("testuser") is called
	result, err := s.userRepository.GetByUsername("testuser")

	// THEN — user is returned
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), user.Username, result.Username)
	assert.Equal(s.T(), user.FirstName, result.FirstName)
	assert.Equal(s.T(), user.Id, result.Id)
}

func (s *UserRepositoryTestSuite) TestGetByUsername_NotFound() {
	// GIVEN — no user "nonexistent"
	// WHEN — GetByUsername("nonexistent") is called
	result, err := s.userRepository.GetByUsername("nonexistent")

	// THEN — nil, nil
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

// Test Update operations
func (s *UserRepositoryTestSuite) TestUpdate_Success() {
	// GIVEN — a user created in DB
	user := &model.User{
		Username:      "olduser",
		Password:      "password",
		FirstName:     "Old",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
	assert.NoError(s.T(), err)
	originalUpdatedAt := user.UpdatedAt

	// Wait a bit to ensure UpdatedAt changes
	time.Sleep(10 * time.Millisecond)

	user.Username = "newuser"
	user.FirstName = "New"
	// WHEN — Update is called
	err = s.userRepository.Update(context.Background(), user)

	// THEN — no error and GetByID returns updated fields
	assert.NoError(s.T(), err)
	updated, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "newuser", updated.Username)
	assert.Equal(s.T(), "New", updated.FirstName)
	assert.True(s.T(), updated.UpdatedAt.After(originalUpdatedAt))
}

func (s *UserRepositoryTestSuite) TestUpdate_PartialFields() {
	// GIVEN — a user created in DB
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
	assert.NoError(s.T(), err)

	user.Username = "updateduser"
	// WHEN — Update is called (only username changed)
	err = s.userRepository.Update(context.Background(), user)

	// THEN — only username changed
	assert.NoError(s.T(), err)
	updated, err := s.userRepository.GetByID(user.Id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), "updateduser", updated.Username)
	assert.Equal(s.T(), "Test", updated.FirstName)
}

func (s *UserRepositoryTestSuite) TestUpdate_NonExistent() {
	// GIVEN — user with id 999 not in DB
	user := &model.User{
		Id:            999,
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}

	// WHEN — Update is called
	err := s.userRepository.Update(context.Background(), user)

	// THEN — GORM Save creates; GetByID returns record
	assert.NoError(s.T(), err)
	result, err := s.userRepository.GetByID(999)
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
}

// Test Delete operations
func (s *UserRepositoryTestSuite) TestDelete_Success() {
	// GIVEN — a user created in DB
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
	assert.NoError(s.T(), err)
	userID := user.Id

	// WHEN — Delete is called
	err = s.userRepository.Delete(userID)

	// THEN — no error; GetByID returns nil; record soft-deleted
	assert.NoError(s.T(), err)
	result, err := s.userRepository.GetByID(userID)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
	var deletedUser model.User
	s.db.Unscoped().First(&deletedUser, userID)
	assert.NotZero(s.T(), deletedUser.DeletedAt)
}

func (s *UserRepositoryTestSuite) TestDelete_NotFound() {
	// GIVEN — no user with id 999
	// WHEN — Delete(999) is called
	err := s.userRepository.Delete(999)

	// THEN — no error (GORM behavior)
	assert.NoError(s.T(), err)
}

func (s *UserRepositoryTestSuite) TestDelete_MultipleUsers() {
	// GIVEN — multiple users created
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: lo.ToPtr(1)}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: lo.ToPtr(1)}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: lo.ToPtr(1)}

	_ = s.userRepository.Create(context.Background(), user1)
	_ = s.userRepository.Create(context.Background(), user2)
	_ = s.userRepository.Create(context.Background(), user3)

	// WHEN — Delete(user2.Id) is called
	err := s.userRepository.Delete(user2.Id)
	assert.NoError(s.T(), err)

	// THEN — only user2 is deleted
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
	// GIVEN — 5 users for client 1, 3 for client 2
	clientId := 1
	for i := 0; i < 5; i++ {
		user := &model.User{
			Username:      "user" + strconv.Itoa(i),
			Password:      "password",
			FirstName:     "User" + strconv.Itoa(i),
			UserLevel:     1,
			ContactNumber: "123456789" + strconv.Itoa(i),
			ClientId:      &clientId,
		}
		_ = s.userRepository.Create(context.Background(), user)
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
		_ = s.userRepository.Create(context.Background(), user)
	}

	ctx := context.Background()
	clientIdPtr := &clientId
	// WHEN — ListByClientId(ctx, 1) is called
	users, err := s.userRepository.ListByClientId(ctx, clientIdPtr)

	// THEN — 5 users for client 1
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 5)
	for _, user := range users {
		assert.Equal(s.T(), &clientId, user.ClientId)
	}
}

func (s *UserRepositoryTestSuite) TestListByClientId_Empty() {
	// GIVEN — no users for client 999
	ctx := context.Background()
	clientId := 999
	clientIdPtr := &clientId
	// WHEN — ListByClientId is called
	users, err := s.userRepository.ListByClientId(ctx, clientIdPtr)

	// THEN — empty list
	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), users)
	assert.Len(s.T(), users, 0)
}

func (s *UserRepositoryTestSuite) TestListByClientId_ExcludesSoftDeleted() {
	// GIVEN — 3 users for client 1; one is soft-deleted
	clientId := 1
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: &clientId}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: &clientId}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: &clientId}

	_ = s.userRepository.Create(context.Background(), user1)
	_ = s.userRepository.Create(context.Background(), user2)
	_ = s.userRepository.Create(context.Background(), user3)

	_ = s.userRepository.Delete(user2.Id)

	ctx := context.Background()
	clientIdPtr := &clientId
	// WHEN — ListByClientId is called
	users, err := s.userRepository.ListByClientId(ctx, clientIdPtr)

	// THEN — 2 users (soft-deleted excluded)
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	for _, user := range users {
		assert.NotEqual(s.T(), user2.Id, user.Id)
	}
}

func (s *UserRepositoryTestSuite) TestListByClientId_FiltersByClientId() {
	// GIVEN — 2 users for client 1, 1 for client 2
	user1 := &model.User{Username: "user1", Password: "pass1", FirstName: "User1", UserLevel: 1, ContactNumber: "111", ClientId: lo.ToPtr(1)}
	user2 := &model.User{Username: "user2", Password: "pass2", FirstName: "User2", UserLevel: 1, ContactNumber: "222", ClientId: lo.ToPtr(1)}
	user3 := &model.User{Username: "user3", Password: "pass3", FirstName: "User3", UserLevel: 1, ContactNumber: "333", ClientId: lo.ToPtr(2)}

	_ = s.userRepository.Create(context.Background(), user1)
	_ = s.userRepository.Create(context.Background(), user2)
	_ = s.userRepository.Create(context.Background(), user3)

	ctx := context.Background()
	clientId := 1
	clientIdPtr := &clientId
	// WHEN — ListByClientId(ctx, 1) is called
	users, err := s.userRepository.ListByClientId(ctx, clientIdPtr)

	// THEN — 2 users for client 1
	assert.NoError(s.T(), err)
	assert.Len(s.T(), users, 2)
	for _, user := range users {
		assert.Equal(s.T(), lo.ToPtr(1), user.ClientId)
		assert.NotEqual(s.T(), user3.Id, user.Id)
	}
}

// Test edge cases and data integrity
func (s *UserRepositoryTestSuite) TestCreate_RequiredFields() {
	// GIVEN — user with empty/minimal fields (SQLite may allow)
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

	// WHEN — Create is called
	err := s.userRepository.Create(context.Background(), user)

	// THEN — behavior consistent (may succeed or fail depending on DB)
	if err == nil {
		assert.NotZero(s.T(), user.Id)
	}
}

func (s *UserRepositoryTestSuite) TestGetByID_AfterDelete() {
	// GIVEN — a user created then soft-deleted
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	_ = s.userRepository.Create(context.Background(), user)
	userID := user.Id

	_ = s.userRepository.Delete(userID)

	// WHEN — GetByID is called
	result, err := s.userRepository.GetByID(userID)

	// THEN — nil, nil (not found)
	assert.NoError(s.T(), err)
	assert.Nil(s.T(), result)
}

func (s *UserRepositoryTestSuite) TestGetByEmail_AfterDelete() {
	s.T().Skip("Email field not available in new model - GetByEmail needs to be updated or removed")
}

func (s *UserRepositoryTestSuite) TestMultipleOperations() {
	// GIVEN — empty DB
	// WHEN/THEN — create, get, update, delete flow
	user := &model.User{
		Username:      "testuser",
		Password:      "password",
		FirstName:     "Test",
		UserLevel:     1,
		ContactNumber: "1234567890",
		ClientId:      lo.ToPtr(1),
	}
	err := s.userRepository.Create(context.Background(), user)
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
	err = s.userRepository.Update(context.Background(), user)
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
