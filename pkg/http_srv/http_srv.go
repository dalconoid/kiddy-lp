package http_srv

import (
	"encoding/json"
	"github.com/dalconoid/kiddy-lp/pkg/storage"
	"github.com/gorilla/mux"
	"net/http"
)

// HTTPServer represents HTTP server
type HTTPServer struct {
	router *mux.Router
}

// New creates server
func New() *HTTPServer {
	s := HTTPServer{router: mux.NewRouter()}
	return &s
}

// Start starts sever
func (s *HTTPServer) Start(addr string) error {
	return http.ListenAndServe(addr, s.router)
}

// ConfigureRouter binds handles to routes
func (s *HTTPServer) ConfigureRouter(storage storage.Storage) {
	s.router.HandleFunc("/ready", handleReady(storage)).Methods("GET")
}

// handleReady handles "/ready"
func handleReady(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := storage.CheckConnection(); err != nil {
			resp := &struct {
				Message string
			}{Message: err.Error()}
			data, _ := json.Marshal(resp)
			http.Error(w, string(data), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}