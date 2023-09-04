package compressor

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type GzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (grw GzipResponseWriter) Write(data []byte) (int, error) {
	return grw.Writer.Write(data)
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		acceptEncoding := r.Header.Get("Accept-Encoding")
		if !strings.Contains(acceptEncoding, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		contentTypes := r.Header.Get("Content-Type")
		if !strings.Contains(contentTypes, "application/json") || !strings.Contains(contentTypes, "text/html") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		gzipWriter := gzip.NewWriter(w)
		defer gzipWriter.Close()

		gzipResponseWriter := GzipResponseWriter{
			Writer:         gzipWriter,
			ResponseWriter: w,
		}

		next.ServeHTTP(gzipResponseWriter, r)
	})
}
