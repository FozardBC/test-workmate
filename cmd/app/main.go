package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-app/internal/api"
	"test-app/internal/config"
	"test-app/internal/logger"
	taskManager "test-app/internal/services/task-manager"
	"test-app/internal/storage/ram"
	"time"
)

func main() {

	ctx := context.Background()

	cfg := config.MustRead()

	log := logger.New(cfg.Log)

	log.Info("App is starting")

	//Starting services(storage, API, taskHandler)

	storage := ram.New(log)

	service := taskManager.New(storage, log)

	API := api.New(log, service)

	API.Setup()

	srv := http.Server{
		Addr:    cfg.ServerHost + ":" + cfg.ServerPort,
		Handler: API.Router,
	}

	// Graceful shutdown
	chanErrors := make(chan error, 1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("HTTP server started", "Addres", srv.Addr)
		chanErrors <- srv.ListenAndServe()
	}()

	go func() {
		log.Info("Started to ping databse")
		for {
			time.Sleep(5 * time.Second)
			err := storage.Ping(ctx)
			if err != nil {
				chanErrors <- err
				break
			}
		}

	}()

	// gracefull shutdown
	select {
	case err := <-chanErrors:
		log.Error("Shutting down. Critical error:", "err", err)

		shutdown <- syscall.SIGTERM
	case sig := <-shutdown:
		log.Error("received signal, starting graceful shutdown", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Error("server graceful shutdown failed", "err", err)
			err = srv.Close()
			if err != nil {
				log.Error("forced shutdown failed", "err", err)
			}
		}

		storage.Close()

		log.Info("shutdown completed")

	}
}
