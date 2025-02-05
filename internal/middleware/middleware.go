package middleware

import (
	"context"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"goth/internal/store"
	"log"
	"net/http"
	"strings"
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Generate Nonces for this request
		nonceSet := Nonces{
			Htmx:            generateRandomString(16),
			ResponseTargets: generateRandomString(16),
			Tw:              generateRandomString(16),
			HtmxCSSHash:     "sha256-pgn1TCGZX6O77zDvy0oTODMOxemn0oj0LeCnQTRj7Kg=",
		}

		// Set nonces in context
		ctx := context.WithValue(r.Context(), NonceKey, nonceSet)

		// Construct the CSP header
		cspHeader := fmt.Sprintf(
			"default-src 'self'; "+
				"script-src 'unsafe-inline' 'nonce-%s' 'nonce-%s' unpkg.com; "+
				"style-src 'unsafe-inline' 'nonce-%s' '%s' fonts.googleapis.com; "+
				"script-src-attr	'unsafe-inline' 'self'; "+
				"script-src-elem	'unsafe-inline' 'self'; "+
				"img-src 'self' assets.mycartpanda.com; "+
				"font-src fonts.gstatic.com fonts.googleapis.com; "+
				// enable inline scripts
				"report-uri /csp-violation-report-endpoint",
			nonceSet.Htmx,
			nonceSet.ResponseTargets,
			nonceSet.Tw,
			nonceSet.HtmxCSSHash,
		)

		// Set the CSP header
		w.Header().Set("Content-Security-Policy", cspHeader)

		// // Pass the updated request down the chain
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
	sessionStore      store.SessionStore
	sessionCookieName string
}

func NewAuthMiddleware(sessionStore store.SessionStore, sessionCookieName string) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore:      sessionStore,
		sessionCookieName: sessionCookieName,
	}
}

type UserContextKey string

var UserKey UserContextKey = "user"

func (m *AuthMiddleware) AddUserToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sessionCookie, err := r.Cookie(m.sessionCookieName)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		decodedValue, err := b64.StdEncoding.DecodeString(sessionCookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(string(decodedValue), ":")

		if len(splitValue) != 2 {
			next.ServeHTTP(w, r)
			return
		}

		sessionID := splitValue[0]
		userID := splitValue[1]

		fmt.Println("sessionID", sessionID)
		fmt.Println("userID", userID)

		user, err := m.sessionStore.GetUserFromSession(sessionID, userID)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), UserKey, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUser(ctx context.Context) *store.User {
	user := ctx.Value(UserKey)
	if user == nil {
		return nil
	}

	return user.(*store.User)
}
