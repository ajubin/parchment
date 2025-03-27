package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

type BasicAuth struct {
	username string
	password string
}

func NewBasicAuth(username string, password string) *BasicAuth {
	return &BasicAuth{username: username, password: password}
}

// BasicAuthMiddleware validates the Basic Auth credentials
func BasicAuthMiddleware(auth BasicAuth, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Validate the provided credentials
		if !isValidBasicAuth(authHeader, auth) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isValidBasicAuth(authHeader string, auth BasicAuth) bool {
	const prefix = "Basic "
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	decoded, err := base64.StdEncoding.DecodeString(authHeader[len(prefix):])
	if err != nil {
		return false
	}

	credentials := strings.SplitN(string(decoded), ":", 2)
	if len(credentials) != 2 {
		return false
	}

	return credentials[0] == auth.username && credentials[1] == auth.password
}

func SecureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func SanitizeMarkdown(input string) string {
	p := bluemonday.StrictPolicy() // 🚀 **Only allows safe text, no HTML**
	return p.Sanitize(input)
}
