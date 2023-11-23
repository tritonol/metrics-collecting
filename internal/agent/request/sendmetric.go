package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	mj "github.com/tritonol/metrics-collecting.git/internal/structs/JSON"
)

const RetryCount = 3

type MetricRequest interface {
	CollectCounter() map[string]int64
	CollectGauge() map[string]float64
}

func sendJSONMetrics(serverAddress, mtype, mname string, mvalue interface{}) error {
	url := fmt.Sprintf("%s/update/", serverAddress)

	var body mj.Metrics
	var delta int64
	var value float64

	switch mtype {
	case "gauge":
		value = mvalue.(float64)
	case "counter":
		delta = mvalue.(int64)
	default:
		return fmt.Errorf("invalid metric type: %s", mtype)
	}

	body = mj.Metrics{
		ID:    mname,
		MType: mtype,
		Delta: &delta,
		Value: &value,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error requesting metrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned non-200 status code: %s", resp.Status)
	}

	return nil
}

func retryableHTTPPost(ctx context.Context, url string, data *bytes.Buffer) (*http.Response, error) {
	var resp *http.Response
	var lastErr error

	for retry := 0; retry < RetryCount; retry++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err() // Context canceled or deadline exceeded
		default:
			resp, lastErr = http.Post(url, "application/json", data)
			if lastErr == nil && resp.StatusCode == http.StatusOK {
				return resp, nil
			}
			time.Sleep(time.Second * time.Duration((retry+1)*2))
		}
	}

	return resp, lastErr
}

func Send(metricRequest MetricRequest, serverAddress string) {
	gaugeMetrics := metricRequest.CollectGauge()
	counterMetrics := metricRequest.CollectCounter()

	for metric, value := range gaugeMetrics {
		if err := sendJSONMetrics(serverAddress, "gauge", metric, value); err != nil {
			log.Printf("Error sending gauge metric %s: %v", metric, err)
		}
	}

	for metric, value := range counterMetrics {
		if err := sendJSONMetrics(serverAddress, "counter", metric, value); err != nil {
			log.Printf("Error sending counter metric %s: %v", metric, err)
		}
	}
}

func SendBatch(metricRequest MetricRequest, serverAddress string) {
	url := fmt.Sprintf("%s/updates/", serverAddress)

	metrics := make([]mj.Metrics, 0)

	gaugeMetrics := metricRequest.CollectGauge()
	counterMetrics := metricRequest.CollectCounter()

	for name, gauge := range gaugeMetrics {
		v := gauge
		var standart int64 = 0
		metrics = append(metrics, mj.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &v,
			Delta: &standart,
		})
	}
	for name, counter := range counterMetrics {
		v := counter
		var standart float64 = 0
		metrics = append(metrics, mj.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &v,
			Value: &standart,
		})
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(metrics)
	if err != nil {
		log.Printf("Error encoding metrics: %v", err)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	resp, err := retryableHTTPPost(ctx, url, &buf)
	if err != nil {
		log.Printf("Error sending batch metrics after %d retries: %v", RetryCount, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Server returned non-200 status code: %s", resp.Status)
		return
	}
}
