package handlers

import "net/http"

func AuthMiddleware(repo ServiceRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("X-Service-Token")
			if token == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			svc, err := repo.GetServiceByToken(r.Context(), token)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "service_id", svc.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
