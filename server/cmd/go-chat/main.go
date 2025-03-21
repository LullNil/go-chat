package main

import (
	"context"
	"net/http"
	"os"

	"log/slog"

	ssogrpc "go-chat/internal/clients/sso/grpc"

	"go-chat/echo"
	"go-chat/general"
	"go-chat/internal/app/auth"
	"go-chat/internal/config"
	"go-chat/internal/lib/logger/handlers/slogpretty"
	"go-chat/internal/lib/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting go-chat",
		slog.String("address", cfg.Address),
		slog.String("env", cfg.Env),
		slog.String("version", "0.0.1"),
	)
	log.Debug("debug messages enabled")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to create sso client", sl.Err(err))
		os.Exit(1)
	}

	// Initialize the echo server.
	echoServer := echo.NewEchoServer()
	// Initialize the general chat server.
	generalChat := general.NewGeneralChat()
	// Initialize the auth handler.
	authHandler := auth.NewAuthHandler(ssoClient)

	// Register the endpoints for echo and general chat.
	http.HandleFunc("/ws/echo", echoServer.Handle)
	http.HandleFunc("/ws/general", generalChat.Handle)

	// Register the login and registration endpoints.
	http.HandleFunc("/login", authHandler.LoginHandler)
	http.HandleFunc("/register", authHandler.RegisterHandler)

	// Start the broadcast routine for the general chat.
	go generalChat.Broadcast()

	if err := http.ListenAndServe(cfg.Address, nil); err != nil {
		log.Error("ListenAndServe error:", slog.String("error", err.Error()))
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
