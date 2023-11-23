package routes

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tritonol/metrics-collecting.git/internal/server/config"
	"github.com/tritonol/metrics-collecting.git/internal/storage"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
	"github.com/tritonol/metrics-collecting.git/internal/storage/pgstorage"
	"go.uber.org/zap"
)

func testRequest(t *testing.T, ts *httptest.Server, method,
                    path string) (*http.Response, string) {
    req, err := http.NewRequest(method, ts.URL+path, nil)
    require.NoError(t, err)

    resp, err := ts.Client().Do(req)
    require.NoError(t, err)
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    require.NoError(t, err)

    return resp, string(respBody)
}

func TestRouter(t *testing.T) {
	cfg := config.MustLoad()

	logger, _ := zap.NewProduction()
	ctx := context.Background()
	
	var storage storage.Storage
	var err error

	if cfg.DB.ConnString != "" {
		storage, err = pgstorage.NewPg(ctx, cfg.DB.ConnString)
		if err != nil {
			logger.Error("Can`t connect db", zap.Error(err))
		}
	} else {
		storage = memstorage.NewMemStorage()
	}
	
	ts := httptest.NewServer(MetricRouter(ctx, storage, logger))

	type want struct {
		code        int
		contentType string
	}

	tests := []struct {
		name   string
		method string
		url    string
		want   want
	}{
		{"Valid gauge metric", http.MethodPost, "/update/gauge/metric1/123.45", want{code: http.StatusOK, contentType: "text/plain; charset=utf-8"}},
		{"Valid counter metric", http.MethodPost, "/update/counter/metric2/50", want{code: http.StatusOK, contentType: "text/plain; charset=utf-8"}},
		{"Invalid metric type", http.MethodPost, "/update/invalidtype/metric3/100", want{code: http.StatusBadRequest, contentType: "text/plain; charset=utf-8"}},
		{"Invalid gauge metric value", http.MethodPost, "/update/gauge/metric4/not-a-number", want{code: http.StatusBadRequest, contentType: "text/plain; charset=utf-8"}},
		{"Invalid counter metric value", http.MethodPost, "/update/counter/metric4/12.3", want{code: http.StatusBadRequest, contentType: "text/plain; charset=utf-8"}},
		{"Missing metric name", http.MethodPost, "/update/gauge/789", want{code: http.StatusNotFound, contentType: "text/plain; charset=utf-8"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			resp, _ := testRequest(t, ts, "POST", test.url)
			assert.Equal(t, test.want.code, resp.StatusCode)
			assert.Equal(t, test.want.contentType, resp.Header.Get("Content-Type"))
			resp.Body.Close()
		})
	}
}