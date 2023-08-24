package config

import "flag"

type Config struct {
	Server HTTPServer
}

type HTTPServer struct {
	Address string
}

func MustLoad() *Config {
	var cfg Config
	flag.StringVar(&cfg.Server.Address, "a", ":8080", "address to run server")
	flag.Parse()
	return &cfg
}