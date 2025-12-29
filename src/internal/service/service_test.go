package service

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/weeranieb/boonmafarm-backend/src/internal/repository/mocks"
)

type ServiceTestSuite struct {
	suite.Suite
	userRepo    *mocks.MockUserRepository
	userService UserService
}

func (s *ServiceTestSuite) SetupTest() {
	s.userRepo = mocks.NewMockUserRepository(s.T())
	s.userService = NewUserService(s.userRepo)
}

func (s *ServiceTestSuite) TearDownTest() {
	s.userRepo.ExpectedCalls = nil
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
