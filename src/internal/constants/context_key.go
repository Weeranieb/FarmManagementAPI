package constants

// ContextKey is the type for request context keys (userId, clientId, userLevel, username).
type ContextKey string

const (
	UserIDKey    ContextKey = "userId"
	ClientIDKey  ContextKey = "clientId"
	UserLevelKey ContextKey = "userLevel"
	UsernameKey  ContextKey = "username"
)
