package main

import (
	"context"
	"os/signal"
	"sync"
	"syscall"

	lrucache "github.com/aik27/otus_go_image_previewer/internal/cache/lru"
	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/logger"
	"github.com/aik27/otus_go_image_previewer/internal/server"
)

func main() {
	cfg := config.GetConfig()
	logger.SetupLogger(cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	servWg := &sync.WaitGroup{}
	servWg.Add(1)

	cache := lrucache.NewCache(cfg.Cache.Capacity, lrucache.OnEvictedEvent)
	go server.RunHTTP(ctx, cfg, servWg, cache)

	servWg.Wait()
}
