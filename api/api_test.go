// In file api/handlers_test.go
package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePrintMethodNotAllowed(t *testing.T) {
	server := NewServer(":8080", nil, BasicAuth{"user", "pass"})
	req, err := http.NewRequest(http.MethodGet, "/print", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.handlePrint)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusMethodNotAllowed)
	}
}
