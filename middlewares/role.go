package middlewares

import (
	"net/http"
	"strings"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Extract role from context
			roleVal := r.Context().Value(RoleKey)
			if roleVal == nil {
				http.Error(w, "Role not found in context", http.StatusForbidden)
				return
			}
			role, ok := roleVal.(string)
			if !ok {
				http.Error(w, "Invalid role type in context", http.StatusInternalServerError)
				return
			}

			// check role is match with allowed roles
			for _, allowed := range allowedRoles {
				if strings.EqualFold(role, allowed) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// If not allowed
			http.Error(w, "Forbidden: insufficient role permissions", http.StatusForbidden)
		})
	}
}
