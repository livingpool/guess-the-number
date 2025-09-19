package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/livingpool/middleware"
	"github.com/livingpool/router"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("error loading .env file")
		return
	}

	h := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
	logger := slog.New(h)

	config := middleware.LoggingConfig{
		DefaultLevel:     slog.LevelInfo,
		ServerErrorLevel: slog.LevelError,
		ClientErrorLevel: slog.LevelWarn,
	}

	stack := middleware.CreateStack(
		middleware.Logging(logger, config),
	)

	router := router.Init()

	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}
	port = ":" + port

	server := http.Server{
		Addr:    port,
		Handler: stack(router),
	}

	slog.Info(fmt.Sprintf("server listening on port %s", port))

	log.Fatal(server.ListenAndServe())
}
