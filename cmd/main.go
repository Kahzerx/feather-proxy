package main

import (
	"context"
	"errors"
	"feather-proxy/internal/logger"
	"feather-proxy/internal/proxy"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		apiPort = "8000"
	}
	log := logger.NewColorLogger("API")
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%s", apiPort),
		Handler: proxy.NewProxy(),
	}
	go func() {
		log.Info("listening on port: %s", apiPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Critical("Failed to start the server %s", err.Error())
		}
	}()

	waitForShutdown(server)
}

func waitForShutdown(server *http.Server) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
