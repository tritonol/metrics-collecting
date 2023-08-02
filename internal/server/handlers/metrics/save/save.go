package save

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	gauge string = "gauge"
	counter string = "counter"
)

type MetricSaver interface {
	StoreGauge(name string, value float64)
	IncrCounter(name string, value int64)
}

func New(metricSaver MetricSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 5 {
			http.Error(w, "Invalid URL format", http.StatusNotFound)
			return
		}

		metricType := parts[2]
		metricName := parts[3]

		switch metricType{
		case gauge:
			metricValue, err := strconv.ParseFloat(parts[4], 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			metricSaver.StoreGauge(metricName, metricValue)
		case counter:
			metricValue, err := strconv.ParseInt(parts[4], 10, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			metricSaver.IncrCounter(metricName, metricValue)
		default: 
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
	}
}