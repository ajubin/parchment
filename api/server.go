package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	p "github.com/ajubin/parchment/printer"
	t "github.com/ajubin/parchment/types"
	"github.com/microcosm-cc/bluemonday"
)

type Server struct {
	listenAddr string
	printer    p.Printer
}

func NewServer(listenAddr string, printer p.Printer) *Server {
	return &Server{listenAddr: listenAddr, printer: printer}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:         s.listenAddr,
		ReadTimeout:  10 * time.Second, // ⏳ Prevents slow request attacks
		WriteTimeout: 10 * time.Second, // ⏳ Prevents slow response attacks
		IdleTimeout:  120 * time.Second,
		Handler:      secureMiddleware(mux), // 🔒 Security Middleware
	}

	// mux.HandleFunc("/print", handlePrint(printer))
	mux.HandleFunc("/test-print", s.handleTestPrint)

	return server.ListenAndServe()
}

func secureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleTestPrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	// 🚀 **Limit the request size to 1MB (or adjust as needed)**
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

	s.printer.TestPrint()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}
func (s *Server) handlePrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// 🚀 **Limit the request size to 1MB (or adjust as needed)**
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1) // 1MB limit

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Impossible de lire la requête: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req t.PrintRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	// ✅ **Sanitize Markdown before parsing**
	cleanContent := sanitizeMarkdown(req.Content)

	// Parser le Markdown et générer le buffer
	buffer := parseMarkup(cleanContent)

	// Envoyer au printer
	err = s.printer.Print(buffer)
	if err != nil {
		http.Error(w, "Erreur lors de l'impression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Ok")
}

func parseMarkup(content string) bytes.Buffer {
	var buffer bytes.Buffer
	// buffer.Write(resetPrinter) // Réinitialisation

	lines := strings.Split(content, "\n")

	// TODO: add markup parsing (maybe with markdown)
	for _, line := range lines {
		buffer.WriteString(line + "\n") // Texte normal

	}
	// Ajouter un saut de ligne final
	return buffer
}

func sanitizeMarkdown(input string) string {
	p := bluemonday.StrictPolicy() // 🚀 **Only allows safe text, no HTML**
	return p.Sanitize(input)
}
