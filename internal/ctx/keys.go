package ctx

// 1. Define a private type for keys to prevent collisions with other packages.
// Using 'int' is more memory-efficient than 'string' for context lookups.
type contextKey int

const (
	// 2. Use iota to automatically assign unique values (0, 1, 2...)
	// The first key is 'UserIDKey' (0), the next would be 1, and so on.
	UserKey contextKey = iota

	// Add future keys here easily:
	// UserRoleKey
	// RequestIDKey
)

// ContextUser is the simplified user data we carry through the request
type ContextUser struct {
	ID   int
	Role string
}
