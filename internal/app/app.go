package app

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/realPointer/banknoteExchange/config"
	v1 "github.com/realPointer/banknoteExchange/internal/controller/http/v1"
	"github.com/realPointer/banknoteExchange/internal/service"
	"github.com/realPointer/banknoteExchange/pkg/httpserver"
	"github.com/realPointer/banknoteExchange/pkg/logger"
)

func Run() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Logger
	l, err := logger.New(cfg.Log.Level)
	if err != nil {
		log.Fatalf("Logger error: %s", err)
	}

	// Services
	services := service.NewServices(l)
	l.Debug().Msg("Services created")

	// HTTP Server
	handler := v1.NewRouter(l, services)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	l.Info().Msgf("HTTP server started on port %v", cfg.HTTP.Port)

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	l.Debug().Msg("Waiting for interrupt signal")

	select {
	case s := <-interrupt:
		l.Info().Msg("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Err(err).Msg("app - Run - httpServer.Notify")
	}

	l.Info().Msg("Shutting down the application")

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Err(err).Msg("app - Run - httpServer.Shutdown")
	} else {
		l.Info().Msg("HTTP server shutdown successfully")
	}

	l.Info().Msg("Application shutdown completed")
}
