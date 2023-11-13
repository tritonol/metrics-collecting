package get

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	m "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type metricGetter interface {
	GetMetrics() map[string]m.Metric
	GetMetric(name string, mType string) (m.Metric, error)
}

func Get(storage metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		
		var err error
		var metric m.Metric
		var response string

		switch metricType {
		case counter:
			metric, err = storage.GetMetric(metricName, counter)
			response = fmt.Sprintf("%v", metric.Delta)
		case gauge:
			metric, err = storage.GetMetric(metricName, gauge)
			response = fmt.Sprintf("%v", metric.Value)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Cant find metric", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	}
}

func GetJSON(storage metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var rawMetric m.Metric
		var metric m.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "Invalid JSON string", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case gauge:
			rawMetric, err = storage.GetMetric(metric.ID, "gauge")
			metric.Value = &rawMetric.Value
		case counter:
			rawMetric, err = storage.GetMetric(metric.ID, "counter")
			metric.Delta = &rawMetric.Delta
		}

		if err != nil {
			http.Error(w, "Cant find metric", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metric)
	}
}

func MainPage(storage metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := make([]string, 0, 50)
		data := storage.GetMetrics()

		for k, v := range data {
			switch v.Type {
			case "gauge":
				resp = append(resp, fmt.Sprintf("%s: %f", k, v.Value))
			case "counter":
				resp = append(resp, fmt.Sprintf("%s: %d", k, v.Delta))
			}
		}

		render.HTML(w, r, strings.Join(resp, "\n"))
	}
}
