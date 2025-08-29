// Package handler have all server routes handler functions
package handler

import "net/http"

// SetRoutes registers the application's HTTP routes.
func SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /zerodha-auth", zAuth)
	mux.HandleFunc("/zerodha-authenticated", zAuthenticated)
}
