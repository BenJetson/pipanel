package pipanel

// ContextKey is the type used for context keys by the server package.
type ContextKey string

// RequestIDKey is the key for request IDs set on the incoming context.
const RequestIDKey ContextKey = "requestID"
