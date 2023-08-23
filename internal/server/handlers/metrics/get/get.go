package get

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	jsonstructs "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type metricGetter interface {
	GetCounter(name string) (int64, bool)
	GetGauge(name string) (float64, bool)
	GetAllGauge() map[string]float64
	GetAllCounter() map[string]int64
}

func Get(storage metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		var metric interface{}
		var ok bool

		switch metricType {
		case counter:
			metric, ok = storage.GetCounter(metricName)
		case gauge:
			metric, ok = storage.GetGauge(metricName)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		if !ok {
			http.Error(w, "Cant find metric", http.StatusNotFound)
			return
		}

		response := fmt.Sprintf("%v", metric)
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(response))
	}
}

func GetJSON(storage metricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ok bool
		var value float64
		var delta int64
		var metric jsonstructs.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "Invalid JSON string", http.StatusBadRequest)
			return
		}

		switch metric.MType {
		case gauge:
			value, ok = storage.GetGauge(metric.ID)
			metric.Value = &value
		case counter:
			delta, ok = storage.GetCounter(metric.ID)
			metric.Delta = &delta
		}

		if !ok {
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
		for k, v := range storage.GetAllGauge() {
			resp = append(resp, fmt.Sprintf("%s: %f", k, v))
		}
		for k, v := range storage.GetAllCounter() {
			resp = append(resp, fmt.Sprintf("%s: %d", k, v))
		}

		render.HTML(w, r, strings.Join(resp, "\n"))
	}
}
