package api

import (
	"net/http"
	"strings"

	"gitserv/internal/auth"
)

// AuthMiddleware validates JWT tokens and adds claims to context.
// It does NOT block unauthenticated requests — use RequireAuth for that.
func AuthMiddleware(jwtMgr *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}

			parts := strings.SplitN(header, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				next.ServeHTTP(w, r)
				return
			}

			claims, err := jwtMgr.Validate(parts[1])
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := auth.SetClaimsInContext(r.Context(), claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth is a middleware that blocks unauthenticated requests.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := auth.GetClaimsFromContext(r.Context())
		if claims == nil {
			writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "Authentication required"})
			return
		}
		next.ServeHTTP(w, r)
	})
}
