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
}

func MustLoad() *Config {
	var cfg Config
	addr := flag.String("a", "localhost:8080", "endpoint address")
	flag.Int64("p", 2, "set poll interval")
	flag.Int64("r", 10, "set report interval")
	flag.Parse()

	if poll := os.Getenv("POLL_INTERVAL"); poll != "" {
		poll, err := strconv.ParseInt(poll, 10, 64)
		if err != nil {
			_ = poll
		} else {
			cfg.PollInterval = poll
		}
	}
	if report := os.Getenv("REPORT_INTERVAL"); report != "" {
		report, err := strconv.ParseInt(report, 10, 64)
		if err != nil {
			_ = err
		} else {
			cfg.ReportInterval = report
		}
	}
	if adr := os.Getenv("ADDRESS"); adr != "" {
		addr = &adr
	}

	cfg.Address = "http://" + *addr

	return &cfg
}
