package jwtutil_test

import (
	"boonmafarm/api/utils/jwtutil"
	"errors"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetUserId(t *testing.T) {
	// Mock a Gin context with a user ID set
	ctxWithUserID := &gin.Context{}
	userID := 123
	ctxWithUserID.Set("userId", float64(userID))

	// Mock a Gin context without a user ID set
	ctxWithoutUserID := &gin.Context{}

	tests := []struct {
		name          string
		ctx           *gin.Context
		expectedID    int
		expectedError error
	}{
		{
			name:          "ShouldReturnUserId",
			ctx:           ctxWithUserID,
			expectedID:    userID,
			expectedError: nil,
		},
		{
			name:          "ShouldReturnError",
			ctx:           ctxWithoutUserID,
			expectedID:    0, // We expect the user ID to be 0 when not found
			expectedError: errors.New("user id not found"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userID, err := jwtutil.GetUserId(test.ctx)

			// Check if the returned user ID matches the expected value
			if userID != test.expectedID {
				t.Errorf("Test %q: expected user ID %d, got %d", test.name, test.expectedID, userID)
			}

			// Check if the returned error matches the expected error
			if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && err.Error() != test.expectedError.Error()) {
				t.Errorf("Test %q: expected error %v, got %v", test.name, test.expectedError, err)
			}
		})
	}
}

func TestGetClientId(t *testing.T) {
	// Mock a Gin context with a user ID set
	ctxWithClientID := &gin.Context{}
	clientId := 123
	ctxWithClientID.Set("clientId", float64(clientId))

	// Mock a Gin context without a client ID set
	ctxWithoutClientID := &gin.Context{}

	tests := []struct {
		name          string
		ctx           *gin.Context
		expectedID    int
		expectedError error
	}{
		{
			name:          "ShouldReturnClientId",
			ctx:           ctxWithClientID,
			expectedID:    clientId,
			expectedError: nil,
		},
		{
			name:          "ShouldReturnError",
			ctx:           ctxWithoutClientID,
			expectedID:    0, // We expect the client ID to be 0 when not found
			expectedError: errors.New("client id not found"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			clientId, err := jwtutil.GetClientId(test.ctx)

			// Check if the returned client ID matches the expected value
			if clientId != test.expectedID {
				t.Errorf("Test %q: expected client ID %d, got %d", test.name, test.expectedID, clientId)
			}

			// Check if the returned error matches the expected error
			if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && err.Error() != test.expectedError.Error()) {
				t.Errorf("Test %q: expected error %v, got %v", test.name, test.expectedError, err)
			}
		})
	}
}

func TestGetUserLevel(t *testing.T) {
	// Mock a Gin context with a UserLevel set
	ctxWithUserLevel := &gin.Context{}
	userLevel := 123
	ctxWithUserLevel.Set("userLevel", float64(userLevel))

	// Mock a Gin context without a client ID set
	ctxWithoutUserLevel := &gin.Context{}

	tests := []struct {
		name          string
		ctx           *gin.Context
		expectedID    int
		expectedError error
	}{
		{
			name:          "ShouldReturnUserLevel",
			ctx:           ctxWithUserLevel,
			expectedID:    userLevel,
			expectedError: nil,
		},
		{
			name:          "ShouldReturnError",
			ctx:           ctxWithoutUserLevel,
			expectedID:    0, // We expect the user level to be 0 when not found
			expectedError: errors.New("user level not found"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			userLevel, err := jwtutil.GetUserLevel(test.ctx)

			// Check if the returned user level matches the expected value
			if userLevel != test.expectedID {
				t.Errorf("Test %q: expected client ID %d, got %d", test.name, test.expectedID, userLevel)
			}

			// Check if the returned error matches the expected error
			if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && err.Error() != test.expectedError.Error()) {
				t.Errorf("Test %q: expected error %v, got %v", test.name, test.expectedError, err)
			}
		})
	}
}

func TestGetUsername(t *testing.T) {
	// Mock a Gin context with a user ID set
	ctxWithUsername := &gin.Context{}
	username := "test"
	ctxWithUsername.Set("username", username)

	// Mock a Gin context without a username set
	ctxWithoutUsername := &gin.Context{}

	tests := []struct {
		name          string
		ctx           *gin.Context
		expectedID    string
		expectedError error
	}{
		{
			name:          "ShouldReturnUsername",
			ctx:           ctxWithUsername,
			expectedID:    username,
			expectedError: nil,
		},
		{
			name:          "ShouldReturnError",
			ctx:           ctxWithoutUsername,
			expectedID:    "", // We expect the username to be "" when not found
			expectedError: errors.New("username not found"),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			username, err := jwtutil.GetUsername(test.ctx)

			// Check if the returned client ID matches the expected value
			if username != test.expectedID {
				t.Errorf("Test %q: expected username %s, got %s", test.name, test.expectedID, username)
			}

			// Check if the returned error matches the expected error
			if (err == nil && test.expectedError != nil) || (err != nil && test.expectedError == nil) || (err != nil && err.Error() != test.expectedError.Error()) {
				t.Errorf("Test %q: expected error %v, got %v", test.name, test.expectedError, err)
			}
		})
	}
}
