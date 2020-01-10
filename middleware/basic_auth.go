package middleware

// Middleware: Basic Auth
type BasicAuthMiddleware struct {
	username string
	password string
	headerName string
}

const BASIC_AUTH_HEADER = ""