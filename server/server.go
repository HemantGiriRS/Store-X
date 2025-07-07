package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func SetupRoutes() http.Handler {
	r := mux.NewRouter()

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(struct{ Message string }{Message: "server is running"})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logrus.Errorf("Error encoding response: %v", err)
			return
		}
	}).Methods("GET")

	return r
}
