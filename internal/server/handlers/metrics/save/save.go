package save

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	m "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const (
	gauge   string = "gauge"
	counter string = "counter"
)

type MetricSaver interface {
	StoreMetric(ctx context.Context, name string, mType string, value float64, delta int64) error
	BatchUpdate(ctx context.Context, metrics []m.Metrics) error
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
			metricSaver.StoreMetric(r.Context(), metricName, metricType, metricValue, 0)
		case counter:
			metricValue, err := strconv.ParseInt(metricValue, 10, 64)
			if err != nil {
				http.Error(w, "Invalid metric value", http.StatusBadRequest)
				return
			}
			metricSaver.StoreMetric(r.Context(), metricName, metricType, 0, metricValue)
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
}

func NewJSON(metricSaver MetricSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric m.Metrics

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			var syntaxError *json.SyntaxError
			var unmarshalTypeError *json.UnmarshalTypeError

			switch {
			case errors.As(err, &syntaxError):
				msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
				http.Error(w, msg, http.StatusBadRequest)
			case errors.Is(err, io.ErrUnexpectedEOF):
				msg := "Request body contains badly-formed JSON"
				http.Error(w, msg, http.StatusBadRequest)
			case errors.As(err, &unmarshalTypeError):
				msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
				http.Error(w, msg, http.StatusBadRequest)
			case strings.HasPrefix(err.Error(), "json: unknown field "):
				fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
				msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
				http.Error(w, msg, http.StatusBadRequest)
			case errors.Is(err, io.EOF):
				msg := "Request body must not be empty"
				http.Error(w, msg, http.StatusBadRequest)
			default:
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			return
		}

		if metric.ID == "" {
			http.Error(w, "missing param", http.StatusNotFound)
			return
		}

		switch metric.MType {
		case gauge:
			if metric.Value == nil {
				http.Error(w, "Empty value", http.StatusBadRequest)
				return
			}
			err := metricSaver.StoreMetric(r.Context(), metric.ID, metric.MType, *metric.Value, 0)
			if err != nil {
				http.Error(w, "Cant write", http.StatusInternalServerError)
				fmt.Printf("%s", err)
				return
			}
		case counter:
			if metric.Delta == nil {
				http.Error(w, "Empty value", http.StatusBadRequest)
				return
			}
			err := metricSaver.StoreMetric(r.Context(), metric.ID, metric.MType, 0, *metric.Delta)
			if err != nil {
				http.Error(w, "Cant write", http.StatusInternalServerError)
				fmt.Printf("%s", err)
				return
			}
		default:
			http.Error(w, "Invalid metric type", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(metric)
	}
}

func Update(storage MetricSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics []m.Metrics

		err := json.NewDecoder(r.Body).Decode(&metrics)
		if err != nil {
			http.Error(w, "something wrong", http.StatusBadRequest)
		}

		err = storage.BatchUpdate(r.Context(), metrics)
		if err != nil {
			http.Error(w, "something wrong", http.StatusInternalServerError)
		}
	}
}