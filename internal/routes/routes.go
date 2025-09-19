// Package routes contains the routing logic for the application
package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vladwithcode/qrcatalog/internal/auth"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

func NewRouter() http.Handler {
	router := NewCustomServeMux()

	RegisterSectionsRoutes(router)

	// Api
	router.HandleFunc("GET /api/auth", auth.PopulateAuth(CheckAuth))
	router.HandleFunc("POST /api/sign-in", auth.PopulateAuth(SignIn))

	// Serve static files
	fs := http.FileServer(http.Dir("web/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fs))

	router.NotFoundHandleFunc(respondWithNotFound)

	return router
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	a, err := auth.ExtractAuthFromReq(r)
	if err != nil {
		resData := map[string]any{
			"isAuthenticated": false,
			"message":         "No se encontró token de sesión",
		}
		respondWithJSON(w, r, http.StatusUnauthorized, resData)
		return
	}

	if a.ID != "" && a.ID != auth.InvalidTokenID && a.ID != auth.ExpiredTokenID {
		resData := map[string]any{
			"isAuthenticated": true,
			"user":            a.ID,
			"message":         "Usuario autenticado",
		}
		respondWithJSON(w, r, http.StatusOK, resData)
		return
	}

	resData := map[string]any{
		"isAuthenticated": false,
		"message":         "Usuario no autenticado",
	}
	respondWithJSON(w, r, http.StatusUnauthorized, resData)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	a, err := auth.ExtractAuthFromReq(r)
	if err == nil {
		if a.ID != "" && a.ID != auth.InvalidTokenID && a.ID != auth.ExpiredTokenID {
			respondWithJSON(w, r, http.StatusFound, map[string]any{
				"redirect":        "/panel",
				"message":         "Usuario ya ha sido autenticado",
				"isAuthenticated": true,
				"user":            a.ID,
			})
			return
		}
	}

	data := db.User{}
	err = json.NewDecoder(r.Body).Decode(&data)
	if data.Username == "" || data.Password == "" {
		respondWithError(w, r, http.StatusBadRequest, "El nombre de usuario y la contraseña son requeridos", nil)
		return
	}

	user, err := db.GetUserByUsername(data.Username)
	if err != nil {
		respondWithError(w, r, http.StatusUnauthorized, "El nombre de usuario o contraseña son incorrectos", err)
		return
	}

	err = user.ValidatePass(data.Password)
	if err != nil {
		respondWithError(w, r, http.StatusUnauthorized, "El nombre de usuario o contraseña son incorrectos", err)
		return
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		respondWithError(w, r, http.StatusInternalServerError, "Ocurrió un error inesperado", err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Path:     "/",
		HttpOnly: auth.UseHTTPOnlyCookies,
		Secure:   auth.UseSecureCookies,
	})

	respondWithJSON(w, r, http.StatusCreated, map[string]any{
		"redirect":        "/panel",
		"message":         "Usuario autenticado con éxito",
		"isAuthenticated": true,
		"user":            user.ID,
	})
}

func SignOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -24),
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
	})

	respondWithJSON(w, r, http.StatusFound, map[string]any{"redirect": "/"})
}

func respondWithNotFound(w http.ResponseWriter, r *http.Request) {
	resData := map[string]any{
		"error":       "No se encontró la página solicitada",
		"routeExists": false,
	}
	respondWithJSON(w, r, http.StatusNotFound, resData)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, data any) {
	w.Header().Set("Connection", "close")
	w.Header().Set("Content-Type", "application/json")

	// Encode to buffer first to catch errors before writing headers
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[%s] %s JSON marshal failed: %v\n", r.Method, r.URL.Path, err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
		return
	}

	w.WriteHeader(code)
	w.Write(jsonData)
}

func respondWithError(w http.ResponseWriter, r *http.Request, code int, reason string, err error) {
	resData := map[string]string{
		"error": reason,
	}
	respondWithJSON(w, r, code, resData)
	log.Printf("[%s] %s failed: %v\n", r.Method, r.URL.Path, err)
}

func publicMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Gives templ's ctx access to the request path and query params
		req := r.WithContext(context.WithValue(r.Context(), "urlPath", r.URL.Path))
		req = r.WithContext(context.WithValue(r.Context(), "urlQueryParams", r.URL.Query()))

		// Ensure there's a cart id for the catalog to use
		_, err := r.Cookie("cart_id")
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:    "cart_id",
				Value:   uuid.Must(uuid.NewV7()).String(),
				Expires: time.Now().Add(time.Hour * 24 * 30),
				Path:    "/",
			})
		}
		next(w, req)
	})
}
