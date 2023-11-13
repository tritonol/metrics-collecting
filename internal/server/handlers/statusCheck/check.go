package statuscheck

import (
	"context"
	"net/http"
)

type statusChecker interface {
	Ping(context.Context) error
}

func Ping(ctx context.Context, pgstorage statusChecker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := pgstorage.Ping(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}
