package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/btserial"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/mcpcomm"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/ollamaish"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/internal/state"
	"github.com/ghofuobudgeghikhdrhiehdihebd/pc-usbc-int/pkg/httphelpers"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := btserial.Run(ctx, btserial.Config{Device: os.Getenv("AIRLOCK_BT_DEVICE")}); err != nil {
			log.Printf("bluetooth console stopped: %v", err)
		}
	}()

	state := &statepkg.State{
		Mux: http.NewServeMux(),
	}
	ollamaish.RegisterAll(state)
	mcpcomm.RegisterAll(state)
	server := &http.Server{
		Addr:              ":11434",
		Handler:           httphelpers.WithCORS(httphelpers.WithLogging(state.Mux)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serverErr := make(chan error, 1)
	go func() { serverErr <- server.ListenAndServe() }()
	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("HTTP shutdown: %v", err)
		}
	}
}
