package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/livingpool/config"
	"github.com/livingpool/middleware"
	"github.com/livingpool/router"
)

const IsProduction = false

func main() {
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := godotenv.Load(".env"); err != nil {
		slog.Error("error loading .env file")
		return
	}

	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	logger := slog.New(h)

	loggingConfig := middleware.LoggingConfig{
		DefaultLevel:     slog.LevelInfo,
		ServerErrorLevel: slog.LevelError,
		ClientErrorLevel: slog.LevelWarn,
	}

	stack := middleware.CreateStack(
		middleware.Logging(logger, loggingConfig),
	)

	conf, err := config.Setup(IsProduction)
	if err != nil {
		slog.Error("error setting up config", "error", err)
	}

	router := router.Init(conf)

	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}
	port = ":" + port

	server := http.Server{
		Addr:    port,
		Handler: stack(router),
	}

	go func() {
		slog.Info(fmt.Sprintf("server listening on port %s", port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-rootCtx.Done()
	stop()
}
