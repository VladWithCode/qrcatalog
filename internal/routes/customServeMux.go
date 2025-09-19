package routes

import (
	"net/http"
	"os"
)

// customServeMux builds on top of http.ServeMux to provide the ability to customize
// the handler for not found routes
type customServeMux struct {
	*http.ServeMux

	// Must not try to WriteHeader as it is set to 404 by default
	//
	// Defaults to http.NotFoundHandler().ServeHTTP
	notFoundHandle http.HandlerFunc
}

func NewCustomServeMux() *customServeMux {
	return &customServeMux{
		http.NewServeMux(),
		http.NotFoundHandler().ServeHTTP,
	}
}

// Will search for the handler appropiate for the received request, if found
// processes the request normally, otherwise responds with a 404
func (csm *customServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	csm.setCORSHeaders(w)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	handler, pattern := csm.Handler(r)

	if pattern == "" {
		csm.notFoundHandle(w, r)
		return
	}

	handler.ServeHTTP(w, r)
}

// Set the custom NotFoundHandler
func (csm *customServeMux) NotFoundHandleFunc(handler http.HandlerFunc) {
	csm.notFoundHandle = handler
}

func (csm *customServeMux) setCORSHeaders(w http.ResponseWriter) {
	corsAllowOrigin := os.Getenv("CORS_ALLOW_ORIGIN")
	if corsAllowOrigin == "" {
		corsAllowOrigin = "*"
	}

	w.Header().Set("Access-Control-Allow-Origin", corsAllowOrigin)
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}
