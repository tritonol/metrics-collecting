package main

import (
	"net/http"

	"github.com/tritonol/metrics-collecting.git/internal/server/handlers/metrics/save"
	"github.com/tritonol/metrics-collecting.git/internal/storage/memstorage"
)

func main() {
	storage := memstorage.NewMemStorage()

	mux := http.NewServeMux()

	mux.HandleFunc("/update/", save.New(storage))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
