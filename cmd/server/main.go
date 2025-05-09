package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/logger"
	"github.com/aik27/otus_go_image_previewer/internal/server"
)

func main() {
	cnf := config.GetConfig()
	logger.SetupLogger(cnf)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	servWg := &sync.WaitGroup{}
	servWg.Add(1)

	go server.RunHTTP(ctx, cnf, servWg)

	servWg.Wait()
}
