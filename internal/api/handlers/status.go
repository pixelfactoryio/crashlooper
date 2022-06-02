package handlers

import (
	"encoding/json"
	"net/http"
)

type statusHandler struct{}

type status struct {
	Status string `json:"status"`
}

// NewStatusHandler returns a new statusHandler instance.
func NewStatusHandler() http.Handler {
	return &statusHandler{}
}

// ServeHTTP respond with the service status
func (h *statusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(status{
		Status: "OK",
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
