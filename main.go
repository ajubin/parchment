package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"syscall"
	"time"

	"github.com/microcosm-cc/bluemonday"

	p "github.com/ajubin/parchment/printer"
)

// Structure du payload attendu
type PrintRequest struct {
	Content string `json:"content"`
}

func sanitizeMarkdown(input string) string {
	p := bluemonday.StrictPolicy() // 🚀 **Only allows safe text, no HTML**
	return p.Sanitize(input)
}

func handlePrint(p p.Printer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var req PrintRequest
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "Format JSON invalide", http.StatusBadRequest)
			return
		}

		// ✅ **Sanitize Markdown before parsing**
		cleanContent := sanitizeMarkdown(req.Content)

		// Parser le Markdown et générer le buffer
		buffer := parseMarkup(cleanContent)

		// Envoyer au printer
		err = p.Print(buffer)
		if err != nil {
			http.Error(w, "Erreur lors de l'impression", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	}
}

func handleTestPrint(p p.Printer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
			return
		}
		// 🚀 **Limit the request size to 1MB (or adjust as needed)**
		r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

		p.TestPrint()

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	}
}

// 🖨️ **Parse le markup et construit le buffer à imprimer**
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

func main() {
	const PORT = 8080

	var printer p.Printer

	if strings.ToLower(strings.TrimSpace(getEnv("MODE", "dev"))) == "prod" {

		// En mode production, on utilise l'imprimante série
		fmt.Println("Utilisation de l'imprimante en serie")
		printer = &p.SerialPrinter{PortName: "/dev/ttyS0", BaudRate: 9600}
	} else {
		// En mode dev, on utilise le mock
		fmt.Println("Utilisation de l'imprimante en mock")
		printer = &p.MockPrinter{}

	}

	// Créer un multiplexer personnalisé
	mux := http.NewServeMux()
	mux.HandleFunc("/print", handlePrint(printer))
	mux.HandleFunc("/test-print", handleTestPrint(printer))

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", PORT),
		ReadTimeout:  10 * time.Second, // ⏳ Prevents slow request attacks
		WriteTimeout: 10 * time.Second, // ⏳ Prevents slow response attacks
		IdleTimeout:  120 * time.Second,
		Handler:      secureMiddleware(mux), // 🔒 Security Middleware
	}

	fmt.Printf("Serveur démarré sur http://localhost:%d", PORT)
	log.Fatal(server.ListenAndServe())

}

func getEnv(key, defaultValue string) string {
	if value, exists := syscall.Getenv(key); exists {
		return value
	}
	return defaultValue
}
