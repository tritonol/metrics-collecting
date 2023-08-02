package save

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
)

func TestSaveHanndler(t *testing.T) {
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
			storage := memstorage.NewMemStorage()
			request := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(New(storage))
			h(w, request)

			result := w.Result()

			assert.Equal(t, test.want.code, result.StatusCode)
            assert.Equal(t, test.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
