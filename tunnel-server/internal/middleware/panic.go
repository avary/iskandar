package middleware

import (
	"net/http"

	"github.com/igneel64/iskandar/server/internal/logger"
)

func PanicRecoveryMiddleware(next http.Handler, logger logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.PanicRecovered(r.URL.Path, err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
