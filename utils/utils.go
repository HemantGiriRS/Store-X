package utils

import (
	"encoding/json"
	"net/http"
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
	var valid bool = false
	if email == "" {
		return valid
	}
	email = strings.ToLower(email)
	substr := "@remotestate.com"
	valid = strings.Contains(email, substr)
	return valid
}

func GetName(email string) string {
	name := strings.Split(email, "@")[0]
	name = strings.ReplaceAll(name, ".", "  ")
	return name
}
