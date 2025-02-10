//go:build !dev

package routes

import "github.com/Project-IPCA/ipca-backend/server"

// ConfigureSwaggerRoutes is a no-op in production.
func ConfigureSwaggerRoutes(server *server.Server) {
	// Do nothing in production
}