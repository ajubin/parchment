package api

import (
	"net/http"

	"github.com/microcosm-cc/bluemonday"
)

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
