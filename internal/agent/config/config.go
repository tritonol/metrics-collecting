package config

import "flag"

type Config struct {
	Address        string
	ReportInterval int64
	PollInterval   int64
}

func MustLoad() *Config {
	var cfg Config
	addr := flag.String("a", "http://localhost:8080", "endpoint address")
	flag.Int64Var(&cfg.PollInterval, "p", 2, "set poll interval")
	flag.Int64Var(&cfg.ReportInterval, "r", 10, "set report interval")
	flag.Parse()

	cfg.Address = "http://" + *addr
	return &cfg
}
