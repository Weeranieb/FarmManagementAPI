package handler

import (
	"testing"

	"github.com/stretchr/testify/suite"
	mocks "github.com/weeranieb/boonmafarm-backend/src/internal/service/mocks"
)

type HandlerTestSuite struct {
	suite.Suite
	userService *mocks.MockUserService
	userHandler UserHandler
}

func (s *HandlerTestSuite) SetupTest() {
	s.userService = mocks.NewMockUserService(s.T())
	s.userHandler = NewUserHandler(s.userService)
}

func (s *HandlerTestSuite) TearDownTest() {
	s.userService.ExpectedCalls = nil
}

func TestHandlerSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
