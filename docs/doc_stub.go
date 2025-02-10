//go:build !dev

package docs

// Stub to prevent import errors in production
var SwaggerInfo struct {
	Host string
}
