package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"storex/handlers"
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

	r.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST")
	r.HandleFunc("/sign-in", handlers.SignIn).Methods("POST")

	return r
}
