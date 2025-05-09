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

	"github.com/aik27/otus_go_image_previewer/internal/config"
	"github.com/aik27/otus_go_image_previewer/internal/routes"
)

func RunHTTP(ctx context.Context, cnf *config.Config, wg *sync.WaitGroup) {
	listenAddrPort := fmt.Sprintf("%s:%d", cnf.HTTPServer.ListenAddr, cnf.HTTPServer.ListenPort)

	server := &http.Server{
		Addr:         listenAddrPort,
		Handler:      routes.ChiRouter(),
		ReadTimeout:  time.Duration(cnf.HTTPServer.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cnf.HTTPServer.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cnf.HTTPServer.IdleTimeout) * time.Second,
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
