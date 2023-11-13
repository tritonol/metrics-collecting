package config

import (
	"flag"
	"os"
)

type Config struct {
	Server HTTPServer
	Backup Backup
}

type HTTPServer struct {
	Address string
}

type Backup struct{
	StoreInterval int64
	FilePath string
	Restore bool
}

func MustLoad() *Config {
	var cfg Config
	flag.StringVar(&cfg.Server.Address, "a", ":8080", "address to run server")
	flag.StringVar(&cfg.Backup.FilePath, "f", "/tmp/metrics-db.json", "Path to backup file")
	flag.BoolVar(&cfg.Backup.Restore, "r", true, "Load previously saved values from file")
	flag.Int64Var(&cfg.Backup.StoreInterval, "i", 300, "Backup interval")
	flag.Parse()

	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Server.Address = envAddr
	}
	return &cfg
}