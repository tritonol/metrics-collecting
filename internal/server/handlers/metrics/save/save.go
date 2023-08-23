package save

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type MetricSaver interface {
	StoreGauge(name string, value float64)
	IncrCounter(name string, value int64)
}

func New(metricSaver MetricSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		metricValue := chi.URLParam(r, "value")

		if metricName == "" || metricType == "" || metricValue == "" {
			http.Error(w, "missing param", http.StatusNotFound)
			return
		}

		switch metricType {
		case gauge:
			metricValue, err := strconv.ParseFloat(metricValue, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			metricSaver.StoreGauge(metricName, metricValue)
		case counter:
			metricValue, err := strconv.ParseInt(metricValue, 10, 64)
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
	}
}

func NewJSON(metricSaver MetricSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric jsonstructs.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "Invalid JSON string", http.StatusBadRequest)
			return
		}

		if metric.ID == "" {
			http.Error(w, "missing param", http.StatusNotFound)
			return
		}

		switch metric.MType {
		case gauge:
			metricSaver.StoreGauge(metric.ID, *metric.Value)
		case counter:
			metricSaver.IncrCounter(metric.ID, *metric.Delta)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
	}
}
