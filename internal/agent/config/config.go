package config

import (
	"flag"
	"os"
	"strconv"
)

type Config struct {
	Address        string
	ReportInterval int64
	PollInterval   int64
	Key            string
}

func MustLoad() *Config {
	cfg := &Config{
		Address:        "localhost:8080",
		ReportInterval: 2,
		PollInterval:   10,
		Key: "",
	}

	addr := flag.String("a", "localhost:8080", "endpoint address")

	flag.StringVar(&cfg.Key, "k", "", "Secret key")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "set poll interval")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "set report interval")
	flag.Parse()

	cfg.PollInterval = getEnvAsInt64("POLL_INTERVAL", cfg.PollInterval)
	cfg.ReportInterval = getEnvAsInt64("REPORT_INTERVAL", cfg.ReportInterval)
	cfg.Key = getEnv("KEY", cfg.Key)

	cfg.Address = "http://" + getEnv("ADDRESS", *addr)

	return cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt64(name string, defaultVal int64) int64 {
	strVal := getEnv(name, "")
	if value, err := strconv.ParseInt(strVal, 10, 64); err == nil {
		return value
	}

	return defaultVal
}
