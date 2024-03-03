package middleware

import (
	"go.uber.org/zap"
	"net/http"
)

const (
	headerKey            = "Content-Type"
	contentType          = "application/json"
	invalidMethodMessage = "method not allowed"
)

func MethodValidationMiddleware(logger *zap.SugaredLogger, httpMethod string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != httpMethod {
			w.Header().Set(headerKey, contentType)
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err := w.Write([]byte(invalidMethodMessage))
			if err != nil {
				logger.Infof("ошибка при ответе - %s", err)
				return
			}
		} else {
			logger.Infof("метод: %s, запрос: %v", r.Method, r)
			next.ServeHTTP(w, r)
		}
	}
}
