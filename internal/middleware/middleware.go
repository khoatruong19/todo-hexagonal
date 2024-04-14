package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo-hexagonal/internal/adapters/primary/web/httperror"
	"todo-hexagonal/internal/constants"
	"todo-hexagonal/internal/core/services"

	"github.com/gorilla/sessions"
)

type key string

var NonceKey key = "nonces"

type Nonces struct {
	Htmx            string
	ResponseTargets string
	Tw              string
	HtmxCSSHash     string
}

func generateRandomString(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}

func CSPMiddleware(next http.Handler) http.Handler {
	// To use the same nonces in all responses, move the Nonces
	// struct creation to here, outside the handler.

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a new Nonces struct for every request when here.
		// move to outside the handler to use the same nonces in all responses
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// set nonces in context
		ctx := context.WithValue(r.Context(), NonceKey, nonceSet)
		// insert the nonces into the content security policy header
		cspHeader := fmt.Sprintf("default-src 'self'; script-src 'nonce-%s' 'nonce-%s' ; style-src 'nonce-%s' '%s';",
			nonceSet.Htmx,
			nonceSet.ResponseTargets,
			nonceSet.Tw,
			nonceSet.HtmxCSSHash)
		w.Header().Set("Content-Security-Policy", cspHeader)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TextHTMLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// get the Nonce from the context, it is a struct called Nonces,
// so we can get the nonce we need by the key, i.e. HtmxNonce
func GetNonces(ctx context.Context) Nonces {
	nonceSet := ctx.Value(NonceKey)
	if nonceSet == nil {
		log.Fatal("error getting nonce set - is nil")
	}

	nonces, ok := nonceSet.(Nonces)

	if !ok {
		log.Fatal("error getting nonce set - not ok")
	}

	return nonces
}

func GetHtmxNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)

	return nonceSet.Htmx
}

func GetResponseTargetsNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.ResponseTargets
}

func GetTwNonce(ctx context.Context) string {
	nonceSet := GetNonces(ctx)
	return nonceSet.Tw
}

type AuthMiddleware struct {
	sessionStore      *sessions.CookieStore
	sessionCookieName string
	userService       *services.UserService
}

type NewAuthMiddlewareParams struct {
	SessionStore      *sessions.CookieStore
	SessionCookieName string
	UserService       *services.UserService
}

func NewAuthMiddleware(params NewAuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      params.SessionStore,
		sessionCookieName: params.SessionCookieName,
		userService:       params.UserService,
	}
}

type UserContextKey string

var UserKey UserContextKey = UserContextKey(constants.UserKey)

func (m *AuthMiddleware) ValidateSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := m.sessionStore.Get(r, m.sessionCookieName)
		if auth, ok := sess.Values[constants.AuthKey].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		userId, ok := sess.Values[constants.UserIdKey].(string)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := m.userService.GetUser(userId)
		if err != nil {
			httperror.ServerErrorResponse(w)
			return
		}

		timezone, ok := sess.Values[constants.TimezoneKey].(time.Time)
		if !ok || time.Until(timezone) < 5*time.Minute {
			sess.Options.MaxAge = 60 * 60
			sess.Values[constants.TimezoneKey] = time.Now()
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) RedirectToHomeIfLoggedIn(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, _ := m.sessionStore.Get(r, m.sessionCookieName)
		authorized := true

		if auth, ok := sess.Values[constants.AuthKey].(bool); !ok || !auth {
			authorized = false
		}

		userId, ok := sess.Values[constants.UserIdKey].(string)
		if !ok {
			authorized = false
		}

		_, err := m.userService.GetUser(userId)
		if err != nil {
			authorized = false
		}

		urlPath := r.URL.Path
		if authorized && (urlPath == "/login" || urlPath == "/register") {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}
