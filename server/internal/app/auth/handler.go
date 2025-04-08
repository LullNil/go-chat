package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	mw "go-chat/internal/app/middleware"
	ssogrpc "go-chat/internal/clients/sso/grpc"
	"go-chat/internal/lib/logger/sl"
	"go-chat/internal/lib/response"
	"go-chat/internal/storage/redis"
)

type Handler struct {
	SSOClient   *ssogrpc.Client
	RedisClient *redis.Client
	Log         *slog.Logger
}

func NewHandler(ssoClient *ssogrpc.Client, redisClient *redis.Client, log *slog.Logger) *Handler {
	return &Handler{
		SSOClient:   ssoClient,
		RedisClient: redisClient,
		Log:         log,
	}
}

// Register handles requests to register a user.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	_, err := h.SSOClient.Register(r.Context(), req.Email, req.Password,
		req.Username)
	if err != nil {
		h.Log.Error("gRPC register error", sl.Err(err))
		response.RespondError(w, http.StatusInternalServerError,
			"User already exists")
		return
	}

	response.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "register_success",
	})
}

// Login handles requests to login a user.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Username string `json:"username"`
		AppID    int32  `json:"app_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.RespondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	token, err := h.SSOClient.Login(r.Context(), req.Email, req.Password,
		req.Username, req.AppID)
	if err != nil {
		h.Log.Error("gRPC login error", sl.Err(err))
		response.RespondError(w, http.StatusInternalServerError, "Invalid email or password")
		return
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // TODO: Change to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, &cookie)

	response.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "login_success",
	})
}

// GetUser handles requests to get user data.
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	const op = "handler.auth.GetUser"

	userID, ok := mw.GetUserIDFromContext(r.Context())
	if !ok {
		h.Log.Error("userID not found in context",
			slog.String("op", op))
		response.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Try to get user data from Redis.
	userData, err := h.RedisClient.GetSSOUserCache(r.Context(), userID)
	switch {
	case err != nil:
		h.Log.Error("failed to get user data from redis", sl.Err(err),
			slog.String("op", op))
	case userData == nil:
		h.Log.Info("cache miss", slog.Int64("userID", userID),
			slog.String("op", op))
	default:
		h.Log.Info("cache hit", slog.Int64("userID", userID),
			slog.String("op", op))
		response.RespondJSON(w, http.StatusOK, userData)
		return
	}

	// If there is no data in Redis or an error occurred while retrieving,
	// turn to gRPC.
	userData, err = h.SSOClient.GetUser(r.Context(), userID)
	if err != nil {
		h.Log.Error("gRPC get user error", sl.Err(err),
			slog.String("op", op))
		response.RespondError(w, http.StatusInternalServerError, "Failed to get user data")
		return
	}

	// Cache the received data in Redis with a TTL of 5 minutes.
	if err := h.RedisClient.SetSSOUserCache(r.Context(), userID,
		userData, 5*time.Minute); err != nil {
		h.Log.Error("failed to cache user data", sl.Err(err))
	} else {
		h.Log.Info("user data cached successfully", slog.Int(
			"userID", int(userID)))
	}

	response.RespondJSON(w, http.StatusOK, userData)
}

// Logout handles requests to log out a user by expiring the
// authentication cookie.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	expiredCookie := http.Cookie{
		Name:     "token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // TODO: Change to true in production
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, &expiredCookie)
	response.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status": "logged_out",
	})
}
