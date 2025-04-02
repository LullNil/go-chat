package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"go-chat/echo"
	"go-chat/general"
	"go-chat/internal/app/auth"
	"go-chat/internal/app/auth/middleware"
	ssogrpc "go-chat/internal/clients/sso/grpc"
	"go-chat/internal/config"
	"go-chat/internal/lib/jwtutil"
	"go-chat/internal/lib/logger"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/storage/redis"
)

func main() {
	cfg := config.MustLoad()
	redisClient := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	jwtutil.SecretKey = []byte(cfg.AppSecret)

	log := logger.SetupLogger(cfg.Env, os.Stdout)

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
	authHandler := auth.NewAuthHandler(ssoClient, redisClient)

	// Register the endpoints for echo and general chat.
	http.HandleFunc("/ws/echo", echoServer.Handle)
	http.HandleFunc("/ws/general", generalChat.Handle)

	// Register the login and registration endpoints.
	http.HandleFunc("/register", middleware.CORSMiddleware(authHandler.RegisterHandler))
	http.HandleFunc("/login", middleware.CORSMiddleware(authHandler.LoginHandler))
	http.HandleFunc("/get-user", middleware.CORSMiddleware(authHandler.GetUserHandler))
	http.HandleFunc("/logout", middleware.CORSMiddleware(authHandler.LogoutHandler))

	// Start the broadcast routine for the general chat.
	go generalChat.Broadcast()

	if err := http.ListenAndServe(cfg.Address, nil); err != nil {
		log.Error("ListenAndServe error:", slog.String("error", err.Error()))
	}
}
