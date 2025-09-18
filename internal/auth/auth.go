// Package auth contains functions pertaining to user authentication
package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vladwithcode/qrcatalog/internal"
	"github.com/vladwithcode/qrcatalog/internal/db"
)

var (
	// UseSecureCookies is a flag to enable secure cookies, by default they are not secure
	// This should be set to true in production, through the USE_SECURE_COOKIES environment variable
	UseSecureCookies = false
	// UseHTTPOnlyCookies is a flag to enable HTTP only cookies, by default they are HTTP only
	// This may be changed through the USE_HTTP_ONLY_COOKIES environment variable if needed
	UseHTTPOnlyCookies = true
	// DefaultCookieName is the name of the cookie used to store the auth token
	DefaultCookieName = "auth_token"
	// DefaultCookieMaxAge is the max age of the cookie in seconds
	DefaultCookieMaxAge = 60 * 60 * 24 * 7 // 1 week
)

func SetAuthParameters() {
	UseSecureCookies = os.Getenv("USE_SECURE_COOKIES") == "true"
	UseHTTPOnlyCookies = os.Getenv("USE_HTTP_ONLY_COOKIES") != "false"

	envCookieName := os.Getenv("DEFAULT_COOKIE_NAME")
	if envCookieName != "" {
		DefaultCookieName = envCookieName
	}
	envMaxAge, _ := strconv.Atoi(os.Getenv("DEFAULT_COOKIE_MAX_AGE"))
	if envMaxAge > 0 {
		DefaultCookieMaxAge = envMaxAge
	}
}

var (
	ErrInvalidAuth = errors.New("invalid auth")
)

type AccessLevel uint8

const (
	AccessLevelUser AccessLevel = iota
	AccessLevelAdmin
	AccessLevelSuperAdmin
)

type Auth struct {
	ID       string
	Username string
	Fullname string
	Role     string
}

func (a *Auth) HasAccess(reqLv AccessLevel) bool {
	var roleLv AccessLevel = 0
	switch a.Role {
	case db.RoleUser:
		roleLv = 0
	case db.RoleAdmin:
		roleLv = 1
	}

	return roleLv >= reqLv
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
