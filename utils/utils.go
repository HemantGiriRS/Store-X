package utils

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"
)

func DecodeRequest(r *http.Request, req interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}
	return nil
}

func EncodeResponse(w http.ResponseWriter, code int, res interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		return err
	}
	return nil
}

func IsValidEmail(email string) bool {
	if !strings.HasSuffix(strings.ToLower(email), "@remotestate.com") {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func GetName(email string) string {
	name := strings.Split(email, "@")[0]
	name = strings.ReplaceAll(name, ".", "  ")
	return name
}

func IsValidEmployeeType(empType string) bool {
	return empType == "full-time" || empType == "intern" || empType == "freelancer"
}

func IsValidEmployeeRole(role string) bool {
	return role == "admin" || role == "asset_manager" || role == "employee_manager" || role == "employee"
}
