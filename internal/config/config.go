package config

import (
	"time"
)

type AgentCfg struct {
	Host           string
	Port           string
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewConfig(pollInterval int, reportInterval int) *AgentCfg {
	var agentCfg AgentCfg
	agentCfg.Host = "127.0.0.1"
	agentCfg.Port = "8080"
	agentCfg.PollInterval = 2 * time.Second
	agentCfg.ReportInterval = 10 * time.Second
	return &agentCfg
}
