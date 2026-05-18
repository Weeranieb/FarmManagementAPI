package errors

import "net/http"

// codeToHTTP captures every AppError code that doesn't follow its category's
// default mapping. Anything not in this map falls through to the range-based
// switch in HTTPStatusFor.
var codeToHTTP = map[int]int{
	// Generic
	ErrGeneric.Code: http.StatusInternalServerError,

	// Validation
	ErrValidationFailed.Code:   http.StatusUnprocessableEntity, // 500010 → 422
	ErrInvalidRequestBody.Code: http.StatusBadRequest,          // 500011 → 400

	// Auth
	ErrAuthUserAlreadyExists.Code:  http.StatusConflict,     // 500020 → 409
	ErrAuthInvalidCredentials.Code: http.StatusUnauthorized, // 500021 → 401
	ErrAuthTokenInvalid.Code:       http.StatusUnauthorized, // 500022 → 401
	ErrAuthTokenExpired.Code:       http.StatusUnauthorized, // 500023 → 401
	ErrAuthPermissionDenied.Code:   http.StatusForbidden,    // 500024 → 403

	// User
	ErrUserNotFound.Code:               http.StatusNotFound,            // 500030 → 404
	ErrUserInvalidInput.Code:           http.StatusUnprocessableEntity, // 500031 → 422
	ErrUserAlreadyExists.Code:          http.StatusConflict,            // 500032 → 409
	ErrUserEmailAlreadyExists.Code:     http.StatusConflict,            // 500033 → 409
	ErrUserCannotModifySuperAdmin.Code: http.StatusForbidden,           // 500034 → 403
	ErrUserCannotDeleteSelf.Code:       http.StatusForbidden,           // 500035 → 403
	ErrUserCannotAssignSuperAdmin.Code: http.StatusForbidden,           // 500036 → 403

	// Farm
	ErrFarmNotFound.Code:      http.StatusNotFound,            // 500040
	ErrFarmAlreadyExists.Code: http.StatusConflict,            // 500041
	ErrFarmInvalidInput.Code:  http.StatusUnprocessableEntity, // 500042

	// Database
	ErrDatabaseError.Code: http.StatusInternalServerError, // 500050

	// Merchant
	ErrMerchantNotFound.Code:      http.StatusNotFound,
	ErrMerchantAlreadyExists.Code: http.StatusConflict,
	ErrMerchantInvalidInput.Code:  http.StatusUnprocessableEntity,

	// Pond
	ErrPondNotFound.Code:      http.StatusNotFound,
	ErrPondAlreadyExists.Code: http.StatusConflict,
	ErrPondInvalidInput.Code:  http.StatusUnprocessableEntity,
	ErrInvalidFishType.Code:   http.StatusUnprocessableEntity,
	ErrPondSourceNotActive.Code: http.StatusUnprocessableEntity,
	ErrPondNotActive.Code:       http.StatusUnprocessableEntity,
	ErrPondInMaintenance.Code:   http.StatusUnprocessableEntity,

	// Worker
	ErrWorkerNotFound.Code:      http.StatusNotFound,
	ErrWorkerAlreadyExists.Code: http.StatusConflict,
	ErrWorkerInvalidInput.Code:  http.StatusUnprocessableEntity,

	// FeedCollection
	ErrFeedCollectionNotFound.Code:      http.StatusNotFound,
	ErrFeedCollectionAlreadyExists.Code: http.StatusConflict,
	ErrFeedCollectionInvalidInput.Code:  http.StatusUnprocessableEntity,

	// FeedPriceHistory
	ErrFeedPriceHistoryNotFound.Code:      http.StatusNotFound,
	ErrFeedPriceHistoryAlreadyExists.Code: http.StatusConflict,
	ErrFeedPriceHistoryInvalidInput.Code:  http.StatusUnprocessableEntity,

	// FishSizeGrade
	ErrFishSizeGradeNotFound.Code: http.StatusNotFound,

	// Client
	ErrClientNotFound.Code:      http.StatusNotFound,
	ErrClientAlreadyExists.Code: http.StatusConflict,
	ErrClientInvalidInput.Code:  http.StatusUnprocessableEntity,

	// FarmGroup
	ErrFarmGroupNotFound.Code:      http.StatusNotFound,
	ErrFarmGroupAlreadyExists.Code: http.StatusConflict,
	ErrFarmGroupInvalidInput.Code:  http.StatusUnprocessableEntity,
}

// HTTPStatusFor maps an AppError numeric code to the HTTP status the server
// should return. Specific codes take precedence over range-based fallbacks
// so any new domain code lands on a sensible default without needing a map
// entry (e.g. anything ≥ 500010 that's unmapped becomes 422 — a safer
// "client probably did something wrong" default than 500).
func HTTPStatusFor(code int) int {
	if status, ok := codeToHTTP[code]; ok {
		return status
	}

	switch {
	case code >= 500001 && code <= 500009:
		return http.StatusInternalServerError
	case code >= 500050 && code <= 500059:
		return http.StatusInternalServerError
	case code >= 500010 && code <= 500019:
		return http.StatusUnprocessableEntity
	case code >= 500020 && code <= 500029:
		return http.StatusUnauthorized
	case code >= 500030:
		// Unknown domain code — assume the client's input wasn't acceptable.
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}
