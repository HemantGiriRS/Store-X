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
	requireRoleAEM.HandleFunc("/employees", handlers.GetEmployees).Methods("GET")
	requireRoleAEM.HandleFunc("/employees/{employee_id}", handlers.DeleteEmployee).Methods("DELETE")
	requireRoleAEM.HandleFunc("/employees/{employee_id}/timeline", handlers.GetEmployeeTimeline).Methods("GET")

	//protected routes only for admin and asset manager
	requireRoleAAM := app.PathPrefix("/private").Subrouter()
	requireRoleAAM.Use(middlewares.Auth, middlewares.RequireRole("admin", "asset_manager"))
	requireRoleAAM.HandleFunc("/assets", handlers.CreateAsset).Methods("POST")
	requireRoleAAM.HandleFunc("/assets/{asset_id}/assign", handlers.AssignAsset).Methods("POST")
	requireRoleAAM.HandleFunc("/assets", handlers.GetAssets).Methods("GET")
	requireRoleAAM.HandleFunc("/assets/{asset_id}/unassign", handlers.UnassignAsset).Methods("POST")
	requireRoleAAM.HandleFunc("/assets/{asset_id}", handlers.DeleteAsset).Methods("DELETE")
	requireRoleAAM.HandleFunc("/assets/{asset_id}/service", handlers.SendAssetForService).Methods("POST")
	requireRoleAAM.HandleFunc("/assets/{asset_id}/receive", handlers.ReceiveAssetFromService).Methods("POST")
	requireRoleAAM.HandleFunc("/assets/{asset_id}/timeline", handlers.GetAssetTimeline).Methods("GET")

	//private routes only for admin
	adminOnly := app.PathPrefix("/admin").Subrouter()
	adminOnly.Use(middlewares.Auth, middlewares.RequireRole("admin"))
	adminOnly.HandleFunc("/employees/{employee_id}/role", handlers.UpdateRole).Methods("PATCH")

	return r
}
