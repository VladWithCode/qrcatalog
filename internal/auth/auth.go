// Package auth contains functions pertaining to user authentication
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vladwithcode/qrcatalog/internal"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

var (
	ErrInvalidAuth = errors.New("invalid auth")
)

type Auth struct {
	ID       string
	Username string
	Fullname string
	Role     string
}

type AuthedHandler func(w http.ResponseWriter, r *http.Request, auth *Auth)

type AuthClaims struct {
	ID       string
	Username string
	Fullname string
	Role     string

	jwt.RegisteredClaims
}

const DefaultExpirationTime = time.Hour * 24
const InvalidTokenID = "invalid"
const ExpiredTokenID = "expired"

type AuthCtxKey string

const DefaultAuthCtxKey AuthCtxKey = "auth"

func CreateToken(user *db.User) (string, error) {
	var (
		t *jwt.Token
		k = os.Getenv("JWT_SECRET")
	)
	expTime := time.Now().Add(DefaultExpirationTime)

	t = jwt.NewWithClaims(jwt.SigningMethodHS256, AuthClaims{
		user.ID,
		user.Username,
		user.Fullname,
		user.Role,

		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expTime),
		},
	})

	return t.SignedString([]byte(k))
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
	var (
		t *jwt.Token
		k = os.Getenv("JWT_SECRET")
	)

	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(k), nil
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

func PopulateAuth(next AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieToken, err := r.Cookie("auth_token")
		var auth = &Auth{}
		authedReq := r.WithContext(context.WithValue(r.Context(), DefaultAuthCtxKey, auth))
		defer next(w, authedReq, auth)
		if err != nil {
			auth.ID = InvalidTokenID
			return
		}

		tokenStr := strings.Split(cookieToken.String(), "=")
		if len(tokenStr) < 2 {
			auth.ID = InvalidTokenID
			return
		}

		t, err := ParseToken(tokenStr[1])
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				auth.ID = ExpiredTokenID
				return
			}
			auth.ID = InvalidTokenID
			return
		}

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			var (
				id, ok1       = claims["ID"].(string)
				username, ok2 = claims["Username"].(string)
				role, ok3     = claims["Role"].(string)
				fullname, ok4 = claims["Fullname"].(string)
			)

			if !ok1 || !ok2 || !ok3 || !ok4 {
				auth.ID = InvalidTokenID
				return
			}

			auth.ID = id
			auth.Role = role
			auth.Username = username
			auth.Fullname = fullname
		}
	}
}

func ValidateAuth(next AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookieToken, err := r.Cookie("auth_token")
		if err != nil {
			RejectUnauthenticated(w, r, "No se encontro token")
			return
		}

		tokenStr := strings.Split(cookieToken.String(), "=")
		if len(tokenStr) < 2 {
			RejectUnauthenticated(w, r, "Token invalido")
			return
		}

		t, err := ParseToken(tokenStr[1])
		if err != nil {
			fmt.Printf("ParseToken err: %v\n", err)
			RejectUnauthenticated(w, r, "Sesion Token invalido")
			return
		}

		if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
			var (
				id, ok1       = claims["ID"].(string)
				username, ok2 = claims["Username"].(string)
				role, ok3     = claims["Role"].(string)
				fullname, ok4 = claims["Fullname"].(string)
			)

			if !ok1 || !ok2 || !ok3 || !ok4 {
				return
			}

			a := &Auth{
				ID:       id,
				Username: username,
				Role:     role,
				Fullname: fullname,
			}

			authedReq := r.WithContext(context.WithValue(r.Context(), DefaultAuthCtxKey, a))
			next(w, authedReq, &Auth{
				ID:       id,
				Username: username,
				Role:     role,
				Fullname: fullname,
			})
		} else {
			RejectUnauthenticated(w, r, "Sesion Token invalido")
		}
	}
}

func RejectUnauthenticated(w http.ResponseWriter, r *http.Request, reason string) {
	internal.HandleRedirect("/iniciar-sesion", http.StatusFound, w, r)
}

func ExtractAuthFromReq(r *http.Request) (*Auth, error) {
	auth, ok := r.Context().Value(DefaultAuthCtxKey).(*Auth)
	if !ok || auth == nil {
		return nil, ErrInvalidAuth
	}

	return auth, nil
}

func ExtractAuthFromCtx(ctx context.Context) (*Auth, error) {
	auth, ok := ctx.Value(DefaultAuthCtxKey).(*Auth)
	if !ok || auth == nil {
		return nil, ErrInvalidAuth
	}

	return auth, nil
}
