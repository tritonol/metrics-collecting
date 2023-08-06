package get

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	gauge string = "gauge"
	counter string = "counter"
)

type MetricGetter interface{
	GetCounter(name string) (int64, bool)
	GetGauge(name string) (float64, bool)
	GetAllGauge() (map[string]float64)
	GetAllCounter() (map[string]int64)
}

func Get(storage MetricGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")

		switch metricType {
		case counter:
			metric, ok := storage.GetCounter(metricName)
			if !ok {
				http.Error(w, "Cant find metric", http.StatusNotFound)
				return
			}
			response := strconv.FormatInt(metric, 10)
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(response))
			return
		case gauge:
			metric, ok := storage.GetGauge(metricName)
			if !ok {
				http.Error(w, "Cant find metric", http.StatusNotFound)
				return
			}
			response := strconv.FormatFloat(metric, 'f', -1, 64)

			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(response))
			return
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}
	}
}

func MainPage(storage MetricGetter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {	
		resp := make([]string, 0)
		for k, v := range storage.GetAllGauge() {
			resp = append(resp, fmt.Sprintf("%s: %f", k, v))
		}
		for k, v := range storage.GetAllCounter() {
			resp = append(resp, fmt.Sprintf("%s: %d", k, v))
		}

		render.HTML(w, r, strings.Join(resp, "\n"))
	}
}