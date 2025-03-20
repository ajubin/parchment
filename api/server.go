package api

import (
	"net/http"
	"time"

	p "github.com/ajubin/parchment/printer"
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
		Handler:      SecureMiddleware(mux), // 🔒 Security Middleware
	}

	mux.HandleFunc("/print", s.handlePrint)
	mux.HandleFunc("/test-print", s.handleTestPrint)

	return server.ListenAndServe()
}
