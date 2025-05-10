package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	lrucache "github.com/aik27/otus_go_image_previewer/internal/cache/lru"
	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/routes"
)

func RunHTTP(ctx context.Context, cfg *config.Config, wg *sync.WaitGroup, cache lrucache.Cache) {
	listenAddrPort := fmt.Sprintf("%s:%d", cfg.HTTPServer.ListenAddr, cfg.HTTPServer.ListenPort)

	server := &http.Server{
		Addr:         listenAddrPort,
		Handler:      routes.ChiRouter(ctx, cfg, cache),
		ReadTimeout:  time.Duration(cfg.HTTPServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.HTTPServer.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.HTTPServer.IdleTimeout) * time.Second,
	}

	go func(ctx context.Context) {
		<-ctx.Done()

		slog.Info(fmt.Sprintf("HTTP server shutdown on: '%s'", listenAddrPort))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
			slog.Error(fmt.Sprintf("HTTP server shutdown error: %v", shutdownErr))
		}

		slog.Info(fmt.Sprintf("HTTP server stop on: '%s'", listenAddrPort))

		wg.Done()
	}(ctx)

	listener, listenErr := net.Listen("tcp", listenAddrPort)
	if listenErr != nil {
		slog.Error(fmt.Sprintf("HTTP server start error on '%s': %v", listenAddrPort, listenErr))
		panic("HTTP server start error: " + listenErr.Error())
	}

	slog.Info(fmt.Sprintf("HTTP server start on: '%s'", listenAddrPort))

	if serveErr := server.Serve(listener); !errors.Is(serveErr, http.ErrServerClosed) {
		slog.Error(fmt.Sprintf("HTTP server error: %v", serveErr))
		panic("HTTP server close error: " + serveErr.Error())
	}
}
