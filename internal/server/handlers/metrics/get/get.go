package get

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	gauge string = "gauge"
	counter string = "counter"
)

type metricGetter interface{
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

func MainPage(storage metricGetter) http.HandlerFunc{
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