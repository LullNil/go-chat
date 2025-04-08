package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-chat/echo"
	"go-chat/general"
	"go-chat/internal/app/auth"
	mw "go-chat/internal/app/middleware"
	"go-chat/internal/app/profile"
	ssogrpc "go-chat/internal/clients/sso/grpc"
	"go-chat/internal/config"
	"go-chat/internal/lib/jwtutil"
	"go-chat/internal/lib/logger"
	"go-chat/internal/lib/logger/sl"
	userservice "go-chat/internal/service/user"
	"go-chat/internal/storage/postgres"
	"go-chat/internal/storage/redis"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()
	// Redis init
	redisClient := redis.New(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	// Setting up JWT
	jwtutil.SecretKey = []byte(cfg.AppSecret)
	// Logger init
	log := logger.SetupLogger(cfg.Env, os.Stdout)

	log.Info("starting go-chat",
		slog.String("address", cfg.Address),
		slog.String("env", cfg.Env),
		slog.String("version", "0.0.1"),
	)
	log.Debug("debug messages enabled")

	// PostgreSQL init
	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	pgPool, err := postgres.NewStorage(connectCtx, cfg.Postgres.DSN, 5, 1*time.Second)
	if err != nil {
		log.Error("failed to connect to postgres", sl.Err(err))
		os.Exit(1)
	}
	log.Info("successfully connected to postgresql")

	// Init gRPC SSO client
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

	// Repos init
	userProfileRepo := postgres.NewUserProfileRepository(pgPool)
	// chatRepo := postgres.NewChatRepository(pgPool)
	// messageRepo := postgres.NewMessageRepository(pgPool)
	log.Info("repositories initialized")

	// Init services
	userService := userservice.NewUserService(
		log,
		userProfileRepo,
		ssoClient,
		redisClient,
		cfg.Redis.SSOCacheTTL,
		cfg.Redis.ProfileCacheTTL,
	)
	log.Info("services initialized")

	// Init handlers
	authHandler := auth.NewHandler(ssoClient, redisClient, log)
	profileHandler := profile.NewHandler(userService, log)
	echoServer := echo.NewEchoServer()
	generalChat := general.NewGeneralChat()
	// TODO: init Chat Handler when done

	// Chi router init
	router := chi.NewRouter()

	// Standard Chi middleware
	router.Use(chimiddleware.RequestID)
	router.Use(chimiddleware.RealIP)
	// router.Use(logger.New(log))
	router.Use(chimiddleware.Recoverer)
	router.Use(chimiddleware.Timeout(cfg.HTTPServer.Timeout))

	// Registering routes
	router.Use(mw.CORS)

	router.Route("/api", func(r chi.Router) {
		// Public routes group
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
			// TODO: add route for GET /api/profiles/{username} (public profile view)
		})

		// Protected routes group
		r.Group(func(r chi.Router) {
			r.Use(mw.Authenticate)
			r.Post("/auth/logout", authHandler.Logout)

			r.Route("/users/me", func(r chi.Router) {
				r.Get("/", authHandler.GetUser)
				r.Get("/profile", profileHandler.GetProfile)
				r.Put("/profile", profileHandler.UpdateProfile)
			})
			// TODO: Add route groups for chats and messages
			// r.Get("/api/chats", chatHandler.HandleGetChats)
			// r.Post("/api/chats", chatHandler.HandleCreateChat)
			// r.Get("/api/chats/{chat_id}/messages", chatHandler.HandleGetMessages)
			// r.Post("/api/chats/{chat_id}/messages", chatHandler.HandlePostMessage)
			// r.Put("/api/messages/{message_id}", chatHandler.HandleUpdateMessage)
			// r.Delete("/api/messages/{message_id}", chatHandler.HandleDeleteMessage)
		})
	})

	// WebSocket routes (use Get for connetion request)
	// Auth for WS is handled separately during handshake
	router.Get("/ws/echo", echoServer.Handle)
	router.Get("/ws/general", generalChat.Handle)

	// Launching background processes
	// TODO: Review the logic for stopping this process during a graceful shutdown
	go generalChat.Broadcast()

	// Starting HTTP server
	log.Info("starting http server", slog.String("address", cfg.HTTPServer.Address))

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	// Starting server
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(
			err, http.ErrServerClosed) {
			log.Error("http server error", sl.Err(err))
		}
	}()

	// Waiting for a signal from the OS in the main thread
	sig := <-osSignals
	log.Info("received signal", slog.String("signal", sig.String()))

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(),
		cfg.HTTPServer.ShutDownTimeout)
	defer shutdownCancel()

	log.Info("stopping http server", slog.String("op", "main.shutdown"))

	// Run graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("http server shutdown error", sl.Err(err))
	}

	// Closing resources after stopping the server
	log.Info("closing postgres connection pool")
	pgPool.Close()

	// TODO: add graceful shutdown for WebSocket

	log.Info("application stopped")
}
