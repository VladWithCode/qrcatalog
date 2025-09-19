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
	"github.com/vladwithcode/qrcatalog/internal/templates/pages"
)

func NewRouter() http.Handler {
	router := NewCustomServeMux()

	RegisterDashboardRoutes(router)
	RegisterImagesRoutes(router)
	RegisterImageSelectorRoutes(router)
	RegisterCategoriesRoutes(router)
	RegisterProductsRoutes(router)
	RegisterEventKindsRoutes(router)
	RegisterContactRequestsRoutes(router)
	RegisterCatalogRoutes(router)
	RegisterCartRoutes(router)
	RegisterUserRoutes(router)

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
	resData := map[string]any{
		"isAuthenticated": false,
		"message":         "No se encontró token de sesión",
	}
	a, err := auth.ExtractAuthFromReq(r)
	if err != nil {
		respondWithJSON(w, r, http.StatusUnauthorized, resData)
		return
	}

	if a.ID != "" && a.ID != auth.InvalidTokenID && a.ID != auth.ExpiredTokenID {
		resData["isAuthenticated"] = true
		resData["user"] = a.ID
		respondWithJSON(w, r, http.StatusOK, resData)
		return
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	a, _ := auth.ExtractAuthFromReq(r)
	if a.ID != "" && a.ID != auth.InvalidTokenID && a.ID != auth.ExpiredTokenID {
		respondWithJSON(w, r, http.StatusFound, map[string]any{"redirect": "/panel"})
		return
	}

	signinPage := pages.SignIn
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = signinPage(&pages.FormState{
			ServerError: "Error inesperado",
		}).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error inesperado"))
		}
		return
	}

	username := r.FormValue("user")
	password := r.FormValue("password")

	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = signinPage(&pages.FormState{
			UserError:     "El nombre de usuario es requerido",
			UserValue:     username,
			PasswordError: "La contraseña es requerida",
		}).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error inesperado"))
		}
		return
	}

	user, err := db.GetUserByUsername(username)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		signinPage(&pages.FormState{
			UserError:     "Revisa que el nombre de usuario sea correcto",
			UserValue:     username,
			PasswordError: "Revisa que la contraseña sea correcta",
		}).Render(r.Context(), w)
		log.Printf("signin failed: %v\n", err)
		return
	}

	err = user.ValidatePass(password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		err = signinPage(&pages.FormState{
			UserError:     "Revisa que el nombre de usuario sea correcto",
			UserValue:     username,
			PasswordError: "Revisa que la contraseña sea correcta",
		}).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error inesperado"))
		}
		return
	}

	token, err := auth.CreateToken(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = signinPage(&pages.FormState{
			ServerError: "Error inesperado",
		}).Render(r.Context(), w)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("Error inesperado"))
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24 * 7),
		Path:     "/",
		HttpOnly: true,
		// Secure:   true,
	})

	respondWithJSON(w, r, http.StatusFound, map[string]any{"redirect": "/panel"})
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
	resData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Ocurrió un error inesperado en el servidor"))
		log.Printf("request to [%s] failed to marshal data: %v\n", r.URL.Path, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resData)
}

func respondWithError(w http.ResponseWriter, r *http.Request, code int, reason string, err error) {
	resData := map[string]string{
		"error": reason,
	}
	respondWithJSON(w, r, code, resData)
	log.Printf("request to [%s] failed: %v\n", r.URL.Path, err)
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
