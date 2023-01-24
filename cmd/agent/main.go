package main

import (
	"context"
	"fmt"
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
	cfg := config.NewConfig(2, 10)
	agent := agent.NewAgent(cfg)
	agent.Start(ctx)
	<-sig
	cancel()
}
