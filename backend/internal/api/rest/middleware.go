package rest

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthService struct {
	ID int64
}

type AuthServiceRepo interface {
	GetServiceByToken(ctx context.Context, token string) (AuthService, error)
}

func AuthMiddleware(repo AuthServiceRepo) func(http.Handler) http.Handler {
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

// ApiKeyAuthMiddleware проверяет наличие валидного API ключа в заголовках запроса.
// Ищет ключ в заголовках: "X-API-Key" или "api_key".
// Если ключ отсутствует или не найден в БД, возвращает 401 Unauthorized.
func ApiKeyAuthMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/swagger/") {
			c.Next()
			return
		}

		apiKey := c.GetHeader("X-API-Key")
		if apiKey == "" {
			apiKey = c.GetHeader("api_key")
		}
		if apiKey == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		var id int64
		err := db.QueryRowContext(c.Request.Context(), "SELECT id FROM services WHERE api_key = $1 LIMIT 1", apiKey).Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.Set("service_id", id)
		c.Next()
	}
}
