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

var ValidRoles = map[string]bool{
	"admin":            true,
	"asset_manager":    true,
	"employee_manager": true,
	"employee":         true,
}

var ValidEmpTypes = map[string]bool{"full-time": true,
	"intern":     true,
	"freelancer": true,
}

var ValidAssignmentStatus = map[string]bool{
	"assigned":     true,
	"not_assigned": true,
}

var ValidAssetTypes = map[string]bool{
	"laptop":      true,
	"mouse":       true,
	"monitor":     true,
	"hard-disk":   true,
	"pen-drive":   true,
	"mobile":      true,
	"sim":         true,
	"accessories": true,
}

var ValidAssetStatuses = map[string]bool{
	"available":     true,
	"assigned":      true,
	"waitForRepair": true,
	"service":       true,
	"damage":        true,
	"deleted":       true,
}

var ValidAssetOwnedBy = map[string]bool{
	"RemoteState": true,
	"Client":      true,
}
