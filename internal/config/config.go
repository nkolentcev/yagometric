package config

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v6"
)

type ServerCfg struct {
	Address       string        `env:"ADDRESS"`
	FilePath      string        `env:"STORE_FILE"`
	StoreInterval time.Duration `env:"STORE_INTERVAL"`
	Restore       bool          `env:"RESTORE"`
}

type AgentCfg struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
}

func NewConfig() *AgentCfg {
	var agentCfg AgentCfg

	flag.StringVar(&agentCfg.Address, "a", "127.0.0.1:8080", "srv host and port")
	flag.DurationVar(&agentCfg.PollInterval, "p", 2*time.Second, "update interval")
	flag.DurationVar(&agentCfg.ReportInterval, "r", 10*time.Second, "report interval")
	flag.Parse()
	_ = env.Parse(&agentCfg)

	return &agentCfg
}

func NewServerCfg() *ServerCfg {
	var scfg ServerCfg

	flag.StringVar(&scfg.Address, "a", "127.0.0.1:8080", "srv host and port")
	flag.DurationVar(&scfg.StoreInterval, "i", 30*time.Second, "update cache interval")
	flag.BoolVar(&scfg.Restore, "r", true, "init recover")
	flag.StringVar(&scfg.FilePath, "f", "/tmp/devops-metrics-db.json", "temporary cache filepath")
	flag.Parse()

	_ = env.Parse(&scfg)
	return &scfg
}
