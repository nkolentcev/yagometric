package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nkolentcev/yagometric/internal/agent"
	"github.com/nkolentcev/yagometric/internal/config"
)

type gauge float64
type counter int64

func main() {

	fmt.Println("start")
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	cfg := config.NewConfig()
	log.Printf("%v %v %v", cfg.Address, cfg.PollInterval, cfg.ReportInterval)
	agent := agent.NewAgent(cfg)
	agent.Start(ctx)
	<-sig
	cancel()
}
