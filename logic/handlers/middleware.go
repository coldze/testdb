package handlers

import (
	"net/http"

	"github.com/coldze/testdb/utils"
)

func NewCheckAndSetLoggerMiddleware(newLogger LoggerFactory, makeId IDFactory, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		requestID := makeId()
		ctx := utils.SetRequestID(r.Context(), requestID)
		ctx = utils.SetLogger(ctx, newLogger(requestID))
		next(w, r.WithContext(ctx))
	}
}
