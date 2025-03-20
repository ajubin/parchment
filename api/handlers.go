package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ajubin/parchment/types"
)

func (s *Server) handleTestPrint(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}
	// 🚀 **Limit the request size to 1MB (or adjust as needed)**
	r.Body = http.MaxBytesReader(w, r.Body, 1024*1024) // 1MB limit

	s.printer.TestPrint()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Response{Message: "ok"})
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

	var req types.PrintRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "Format JSON invalide", http.StatusBadRequest)
		return
	}

	cleanContent := SanitizeMarkdown(req.Content)

	buffer := parseMarkup(cleanContent)

	err = s.printer.Print(buffer)
	if err != nil {
		http.Error(w, "Erreur lors de l'impression", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types.Response{Message: "ok"})

}

func parseMarkup(content string) bytes.Buffer {
	var buffer bytes.Buffer

	lines := strings.Split(content, "\n")

	// TODO: add markup parsing (maybe with markdown)
	for _, line := range lines {
		buffer.WriteString(line + "\n") // Texte normal

	}

	return buffer
}
