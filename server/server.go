package server

import (
	"encoding/json"
	"net/http"

	"storex/handlers"
	"storex/middlewares"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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

	//public routes
	app := r.PathPrefix("/store-x").Subrouter()
	app.HandleFunc("/refresh", handlers.RefreshToken).Methods("POST")
	app.HandleFunc("/sign-in", handlers.SignIn).Methods("POST")

	//protected routes only for admin and employee manager
	requireRoleAEM := app.PathPrefix("/private").Subrouter()
	requireRoleAEM.Use(middlewares.Auth, middlewares.RequireRole("admin", "employee_manager"))
	requireRoleAEM.HandleFunc("/register-employee", handlers.SignUp).Methods("POST")

	//private routes only for admin
	adminOnly := app.PathPrefix("/admin").Subrouter()
	adminOnly.Use(middlewares.Auth, middlewares.RequireRole("admin"))
	adminOnly.HandleFunc("/employees/{employee_id}/role", handlers.UpdateRole).Methods("PATCH")
	return r
}
