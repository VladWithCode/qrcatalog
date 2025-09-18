// Package routes contains the routing logic for the application
package routes
func NewRouter() http.Handler {
	router := NewCustomServeMux()

	return router
}
