package compress

import (
	"net/http"
	"strings"

	"github.com/tritonol/metrics-collecting.git/internal/compress"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		contentType := r.Header.Get("Content-Type")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		isValidContentType := strings.HasPrefix(contentType, "application/json") || strings.HasPrefix(contentType, "text/html")

		if supportsGzip && isValidContentType {
			cw := compress.NewCompressWriter(w)
			ow = cw
			
			// cw.Header().Set("Accept-Encoding", "gzip")
			defer cw.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := compress.NewCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, r)
	})
}
