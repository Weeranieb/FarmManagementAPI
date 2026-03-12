package errors

// Generic/Internal errors (500001-500009)
var (
	ErrGeneric = &AppError{
		Code:    500001,
		Message: "An error occurred while processing your request",
	}
)

// Validation errors (500010-500019)
var (
	ErrValidationFailed = &AppError{
		Code:    500010,
		Message: "Validation failed",
	}

	ErrInvalidRequestBody = &AppError{
		Code:    500011,
		Message: "Invalid request body",
	}
)

// Auth errors (500020-500029)
var (
	ErrAuthUserAlreadyExists = &AppError{
		Code:    500020,
		Message: "User already exists",
	}

	ErrAuthInvalidCredentials = &AppError{
		Code:    500021,
		Message: "Invalid username or password",
	}

	ErrAuthTokenInvalid = &AppError{
		Code:    500022,
		Message: "Invalid token",
	}

	ErrAuthTokenExpired = &AppError{
		Code:    500023,
		Message: "Token expired",
	}

	ErrAuthPermissionDenied = &AppError{
		Code:    500024,
		Message: "Permission denied",
	}
)

// User errors (500030-500039)
var (
	ErrUserNotFound = &AppError{
		Code:    500030,
		Message: "User not found",
	}

	ErrUserInvalidInput = &AppError{
		Code:    500031,
		Message: "Invalid user input",
	}

	ErrUserAlreadyExists = &AppError{
		Code:    500032,
		Message: "User already exists",
	}
)

// Farm errors (500040-500049)
var (
	ErrFarmNotFound = &AppError{
		Code:    500040,
		Message: "Farm not found",
	}

	ErrFarmAlreadyExists = &AppError{
		Code:    500041,
		Message: "Farm already exists",
	}

	ErrFarmInvalidInput = &AppError{
		Code:    500042,
		Message: "Invalid farm input",
	}
)

// Database errors (500050-500059)
var (
	ErrDatabaseError = &AppError{
		Code:    500050,
		Message: "Database error occurred",
	}
)

// Merchant errors (500060-500069)
var (
	ErrMerchantNotFound = &AppError{
		Code:    500060,
		Message: "Merchant not found",
	}

	ErrMerchantAlreadyExists = &AppError{
		Code:    500061,
		Message: "Merchant already exists",
	}

	ErrMerchantInvalidInput = &AppError{
		Code:    500062,
		Message: "Invalid merchant input",
	}
)

// Pond errors (500070-500079)
var (
	ErrPondNotFound = &AppError{
		Code:    500070,
		Message: "Pond not found",
	}

	ErrPondAlreadyExists = &AppError{
		Code:    500071,
		Message: "Pond already exists",
	}

	ErrPondInvalidInput = &AppError{
		Code:    500072,
		Message: "Invalid pond input",
	}

	ErrInvalidFishType = &AppError{
		Code:    500073,
		Message: "Invalid fish type",
	}

	ErrPondSourceNotActive = &AppError{
		Code:    500074,
		Message: "Source pond has no active cycle",
	}

	ErrPondNotActive = &AppError{
		Code:    500075,
		Message: "Pond has no active cycle",
	}

	ErrPondInMaintenance = &AppError{
		Code:    500076,
		Message: "Pond is in maintenance; move and sell are not allowed",
	}
)

// Worker errors (500080-500089)
var (
	ErrWorkerNotFound = &AppError{
		Code:    500080,
		Message: "Worker not found",
	}

	ErrWorkerAlreadyExists = &AppError{
		Code:    500081,
		Message: "Worker already exists",
	}

	ErrWorkerInvalidInput = &AppError{
		Code:    500082,
		Message: "Invalid worker input",
	}
)

// FeedCollection errors (500090-500099)
var (
	ErrFeedCollectionNotFound = &AppError{
		Code:    500090,
		Message: "Feed collection not found",
	}

	ErrFeedCollectionAlreadyExists = &AppError{
		Code:    500091,
		Message: "Feed collection already exists",
	}

	ErrFeedCollectionInvalidInput = &AppError{
		Code:    500092,
		Message: "Invalid feed collection input",
	}
)

// FeedPriceHistory errors (500100-500109)
var (
	ErrFeedPriceHistoryNotFound = &AppError{
		Code:    500100,
		Message: "Feed price history not found",
	}

	ErrFeedPriceHistoryAlreadyExists = &AppError{
		Code:    500101,
		Message: "Feed price history already exists",
	}

	ErrFeedPriceHistoryInvalidInput = &AppError{
		Code:    500102,
		Message: "Invalid feed price history input",
	}
)

// Client errors (500110-500119)
var (
	ErrClientNotFound = &AppError{
		Code:    500110,
		Message: "Client not found",
	}

	ErrClientAlreadyExists = &AppError{
		Code:    500111,
		Message: "Client already exists",
	}

	ErrClientInvalidInput = &AppError{
		Code:    500112,
		Message: "Invalid client input",
	}
)
